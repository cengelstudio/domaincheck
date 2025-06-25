package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"domaincheck/internal/config"
	"domaincheck/internal/handlers"
	"domaincheck/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Initialize services
	domainService, err := services.NewDomainService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize domain service: %v", err)
	}

	// Initialize handlers
	domainHandler := handlers.NewDomainHandler(domainService)
	wsHandler := handlers.NewWebSocketHandler(domainService)

	// Setup router
	router := setupRouter(cfg, domainHandler, wsHandler)

	// Setup server
	srv := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Domain Check API Server starting on %s", cfg.Server.Port)
		log.Printf("📁 Domain extensions loaded from: %s", cfg.Domain.ExtensionsFile)
		log.Printf("⚙️  Max concurrent domain checks: %d", cfg.Domain.MaxConcurrentChecks)
		log.Printf("⏱️  Domain check timeout: %s", cfg.Domain.Timeout)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited")
}

func setupRouter(cfg *config.Config, domainHandler *handlers.DomainHandler, wsHandler *handlers.WebSocketHandler) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: true,
	}
	router.Use(cors.New(corsConfig))

	// Request ID middleware
	router.Use(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		c.Header("X-Request-ID", requestID)
		c.Next()
	})

	// Static files for frontend (optional)
	if _, err := os.Stat("./frontend/dist"); err == nil {
		router.Static("/static", "./frontend/dist/static")
		router.LoadHTMLGlob("frontend/dist/*.html")

		// Root route - serve frontend
		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})
	} else {
		// Fallback root route when frontend is not built
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Domain Check API",
				"version": "1.0.0",
				"docs": map[string]string{
					"health":     "GET /api/health",
					"check":      "POST /api/check-domain",
					"history":    "GET /api/domains",
					"extensions": "GET /api/v1/extensions",
					"websocket":  "WS /ws",
				},
			})
		})
	}

	// Setup all API routes
	handlers.SetupRoutes(router, cfg, domainHandler, wsHandler)

	return router
}
