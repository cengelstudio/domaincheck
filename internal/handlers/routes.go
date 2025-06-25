package handlers

import (
	"domaincheck/internal/config"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, cfg *config.Config, domainHandler *DomainHandler, wsHandler *WebSocketHandler) {
	// Health check route
	router.GET("/api/v1/health", domainHandler.HealthCheck)
	router.GET("/api/health", domainHandler.HealthCheck) // Backward compatibility

	// WebSocket route
	router.GET("/ws", wsHandler.HandleWebSocket)

	// Domain routes
	setupDomainRoutes(router, domainHandler)

	// Extension routes
	setupExtensionRoutes(router, domainHandler)
}

// setupDomainRoutes configures domain-related routes
func setupDomainRoutes(router *gin.Engine, domainHandler *DomainHandler) {
	// v1 API routes
	domainsV1 := router.Group("/api/v1/domains")
	{
		domainsV1.POST("/check", domainHandler.CheckDomain)
		domainsV1.POST("/check-all-extensions", domainHandler.CheckAllExtensions)
		domainsV1.POST("/check-multiple", domainHandler.CheckMultipleDomains)
		domainsV1.GET("/history", domainHandler.GetDomainHistory)
		domainsV1.DELETE("/history", domainHandler.ClearHistory)
		domainsV1.GET("/whois/:domain", domainHandler.GetWhoisInfo)
	}

	// Backward compatibility routes (v0)
	domainsV0 := router.Group("/api")
	{
		domainsV0.POST("/check-domain", domainHandler.CheckDomain)
		domainsV0.POST("/check-all-extensions", domainHandler.CheckAllExtensions)
		domainsV0.GET("/domains", domainHandler.GetDomainHistory)
	}
}

// setupExtensionRoutes configures extension-related routes
func setupExtensionRoutes(router *gin.Engine, domainHandler *DomainHandler) {
	extensions := router.Group("/api/v1/extensions")
	{
		extensions.GET("/", domainHandler.GetValidExtensions)
		extensions.POST("/reload", domainHandler.ReloadExtensions)
	}
}
