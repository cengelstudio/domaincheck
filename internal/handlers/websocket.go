package handlers

import (
	"log"
	"net/http"
	"sync"
	"time"

	"domaincheck/internal/models"
	"domaincheck/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	domainService *services.DomainService
	clients       map[*websocket.Conn]bool
	mutex         sync.RWMutex
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(domainService *services.DomainService) *WebSocketHandler {
	return &WebSocketHandler{
		domainService: domainService,
		clients:       make(map[*websocket.Conn]bool),
	}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Register client
	h.mutex.Lock()
	h.clients[conn] = true
	h.mutex.Unlock()

	// Remove client when connection closes
	defer func() {
		h.mutex.Lock()
		delete(h.clients, conn)
		h.mutex.Unlock()
	}()

	// Send welcome message
	welcomeMsg := models.WebSocketMessage{
		Type:    "connected",
		Message: "WebSocket connection established",
	}
	conn.WriteJSON(welcomeMsg)

	// Handle incoming messages
	for {
		var msg models.WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		switch msg.Type {
		case "check_all_extensions":
			h.handleCheckAllExtensions(conn, msg)
		case "ping":
			conn.WriteJSON(models.WebSocketMessage{
				Type: "pong",
				Data: time.Now().Unix(),
			})
		default:
			conn.WriteJSON(models.WebSocketMessage{
				Type:    "error",
				Message: "Unknown message type",
			})
		}
	}
}

// handleCheckAllExtensions handles bulk domain extension checks via WebSocket
func (h *WebSocketHandler) handleCheckAllExtensions(conn *websocket.Conn, msg models.WebSocketMessage) {
	// Extract domain name from message
	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		conn.WriteJSON(models.WebSocketMessage{
			Type:    "error",
			Message: "Invalid message format",
		})
		return
	}

	domainName, ok := data["domain_name"].(string)
	if !ok || domainName == "" {
		conn.WriteJSON(models.WebSocketMessage{
			Type:    "error",
			Message: "Domain name is required",
		})
		return
	}

	// Send start message
	conn.WriteJSON(models.WebSocketMessage{
		Type: "bulk_check_started",
		Data: map[string]interface{}{
			"domain_name": domainName,
		},
	})

	// Start bulk check with progress updates
	go h.performBulkCheckWithProgress(conn, domainName)
}

// performBulkCheckWithProgress performs bulk check and sends progress updates
func (h *WebSocketHandler) performBulkCheckWithProgress(conn *websocket.Conn, domainName string) {
	startTime := time.Now()

	// Get all extensions
	extensions := h.domainService.GetValidExtensions()
	totalExtensions := len(extensions)

	progress := models.WebSocketBulkProgress{
		DomainName:         domainName,
		TotalExtensions:    totalExtensions,
		CheckedCount:       0,
		AvailableCount:     0,
		UnavailableCount:   0,
		ErrorCount:         0,
		AvailableDomains:   []models.WebSocketDomainCheck{},
		UnavailableDomains: []models.WebSocketDomainCheck{},
		IsComplete:         false,
	}

	// Create channels for concurrent processing
	results := make(chan models.WebSocketDomainCheck, totalExtensions)
	errors := make(chan error, totalExtensions)

	// Start concurrent domain checks
	semaphore := make(chan struct{}, 10) // Limit concurrent connections
	var wg sync.WaitGroup

	for _, ext := range extensions {
		wg.Add(1)
		go func(extension string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			fullDomain := domainName + "." + extension
			start := time.Now()

			result, err := h.domainService.CheckDomain(nil, fullDomain)
			responseTime := time.Since(start).Milliseconds()

			if err != nil {
				errors <- err
				return
			}

			wsResult := models.WebSocketDomainCheck{
				Domain:       fullDomain,
				Status:       result.Domain.Status,
				IP:           result.Domain.IP,
				ResponseTime: responseTime,
				CheckedAt:    time.Now().Format(time.RFC3339),
			}

			results <- wsResult
		}(ext)
	}

	// Close channels when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Process results
	for result := range results {
		progress.CheckedCount++
		progress.CurrentDomain = &result

		if result.Status == "Available" {
			progress.AvailableCount++
			progress.AvailableDomains = append(progress.AvailableDomains, result)
		} else {
			progress.UnavailableCount++
			progress.UnavailableDomains = append(progress.UnavailableDomains, result)
		}

		// Send progress update
		conn.WriteJSON(models.WebSocketMessage{
			Type: "bulk_check_progress",
			Data: progress,
		})
	}

	// Process errors
	for range errors {
		progress.ErrorCount++
		progress.CheckedCount++
	}

	// Mark as complete
	progress.IsComplete = true
	progress.TotalTime = time.Since(startTime).Milliseconds()
	progress.CurrentDomain = nil

	// Send final result
	conn.WriteJSON(models.WebSocketMessage{
		Type: "bulk_check_complete",
		Data: progress,
	})
}

// BroadcastToAll sends a message to all connected clients
func (h *WebSocketHandler) BroadcastToAll(msg models.WebSocketMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("Failed to send message to client: %v", err)
			client.Close()
			delete(h.clients, client)
		}
	}
}
