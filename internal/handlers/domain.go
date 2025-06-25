package handlers

import (
	"net/http"
	"strconv"
	"time"

	"domaincheck/internal/models"
	"domaincheck/internal/services"

	"github.com/gin-gonic/gin"
)

// DomainHandler handles domain-related HTTP requests
type DomainHandler struct {
	domainService *services.DomainService
	startTime     time.Time
}

// NewDomainHandler creates a new domain handler
func NewDomainHandler(domainService *services.DomainService) *DomainHandler {
	return &DomainHandler{
		domainService: domainService,
		startTime:     time.Now(),
	}
}

// HealthCheck handles health check requests
func (h *DomainHandler) HealthCheck(c *gin.Context) {
	uptime := time.Since(h.startTime)

	response := models.HealthResponse{
		Status:      "healthy",
		Version:     "1.0.0",
		Timestamp:   time.Now(),
		Uptime:      uptime.String(),
		Environment: "development",
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
		Message: "Service is healthy",
	})
}

// CheckDomain handles single domain check requests
func (h *DomainHandler) CheckDomain(c *gin.Context) {
	startTime := time.Now()

	var request models.DomainCheckRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	// Check domain
	result, err := h.domainService.CheckDomain(c.Request.Context(), request.Domain)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Domain check failed",
			Error:   err.Error(),
		})
		return
	}

	// Calculate process time
	processTime := time.Since(startTime).Milliseconds()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    result,
		Message: "Domain check completed successfully",
		Meta: &models.Meta{
			ProcessTime: processTime,
			RequestID:   c.GetHeader("X-Request-ID"),
		},
	})
}

// CheckAllExtensions handles domain name check for all available extensions
func (h *DomainHandler) CheckAllExtensions(c *gin.Context) {
	startTime := time.Now()

	var request struct {
		DomainName string `json:"domain_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	// Check domain name with all extensions
	result, err := h.domainService.CheckAllExtensions(c.Request.Context(), request.DomainName)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Domain extensions check failed",
			Error:   err.Error(),
		})
		return
	}

	// Calculate process time
	processTime := time.Since(startTime).Milliseconds()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    result,
		Message: "Domain extensions check completed successfully",
		Meta: &models.Meta{
			Total:       len(result.AllResults),
			ProcessTime: processTime,
			RequestID:   c.GetHeader("X-Request-ID"),
		},
	})
}

// CheckMultipleDomains handles multiple domain check requests
func (h *DomainHandler) CheckMultipleDomains(c *gin.Context) {
	startTime := time.Now()

	var request struct {
		Domains []string `json:"domains" binding:"required,min=1,max=50"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	// Check domains
	results, err := h.domainService.CheckMultipleDomains(c.Request.Context(), request.Domains)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Domain check failed",
			Error:   err.Error(),
		})
		return
	}

	// Calculate process time
	processTime := time.Since(startTime).Milliseconds()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    results,
		Message: "Domain checks completed successfully",
		Meta: &models.Meta{
			Total:       len(results),
			ProcessTime: processTime,
			RequestID:   c.GetHeader("X-Request-ID"),
		},
	})
}

// GetDomainHistory returns domain check history
func (h *DomainHandler) GetDomainHistory(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Get history
	history := h.domainService.GetDomainHistory()
	total := len(history)

	// Calculate pagination
	totalPages := (total + perPage - 1) / perPage
	start := (page - 1) * perPage
	end := start + perPage

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	// Get page data
	var pageData []models.Domain
	if start < end {
		pageData = history[start:end]
	} else {
		pageData = []models.Domain{}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    pageData,
		Message: "Domain history retrieved successfully",
		Meta: &models.Meta{
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
			RequestID:  c.GetHeader("X-Request-ID"),
		},
	})
}

// ClearHistory clears domain check history
func (h *DomainHandler) ClearHistory(c *gin.Context) {
	h.domainService.ClearHistory()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Domain history cleared successfully",
	})
}

// GetValidExtensions returns list of valid domain extensions
func (h *DomainHandler) GetValidExtensions(c *gin.Context) {
	extensions := h.domainService.GetValidExtensions()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    extensions,
		Message: "Valid extensions retrieved successfully",
		Meta: &models.Meta{
			Total:     len(extensions),
			RequestID: c.GetHeader("X-Request-ID"),
		},
	})
}

// ReloadExtensions reloads domain extensions from file
func (h *DomainHandler) ReloadExtensions(c *gin.Context) {
	if err := h.domainService.ReloadExtensions(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to reload extensions",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Extensions reloaded successfully",
	})
}

// GetWhoisInfo handles WHOIS information requests
func (h *DomainHandler) GetWhoisInfo(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Domain parameter is required",
		})
		return
	}

	// Get WHOIS information
	whoisInfo, err := h.domainService.GetWhoisInfo(c.Request.Context(), domain)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Failed to get WHOIS information",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    whoisInfo,
		Message: "WHOIS information retrieved successfully",
	})
}
