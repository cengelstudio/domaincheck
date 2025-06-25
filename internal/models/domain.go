package models

import (
	"time"
)

// Domain represents a domain check result
type Domain struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Extension    string    `json:"extension"`
	Available    bool      `json:"available"`
	Status       string    `json:"status"` // "Available", "Registered", "Error"
	IP           string    `json:"ip,omitempty"`
	DNSResolved  bool      `json:"dns_resolved"`
	CheckedAt    time.Time `json:"checked_at"`
	ResponseTime int64     `json:"response_time_ms"`
	Error        string    `json:"error,omitempty"`
}

// DomainCheckRequest represents the request payload for domain checking
type DomainCheckRequest struct {
	Domain string `json:"domain" binding:"required" validate:"required,min=1"`
}

// DomainCheckResponse represents the response for domain checking
type DomainCheckResponse struct {
	Domain       *Domain `json:"domain"`
	IsValidTLD   bool    `json:"is_valid_tld"`
	SupportedTLD bool    `json:"supported_tld"`
}

// AllExtensionsCheckResult represents the result for checking all extensions
type AllExtensionsCheckResult struct {
	DomainName         string                `json:"domain_name"`
	TotalExtensions    int                   `json:"total_extensions"`
	AvailableCount     int                   `json:"available_count"`
	UnavailableCount   int                   `json:"unavailable_count"`
	ErrorCount         int                   `json:"error_count"`
	CheckedAt          time.Time             `json:"checked_at"`
	TotalTime          int64                 `json:"total_time_ms"`
	AvailableDomains   []DomainCheckResponse `json:"available_domains"`
	UnavailableDomains []DomainCheckResponse `json:"unavailable_domains"`
	ErrorDomains       []DomainCheckResponse `json:"error_domains"`
	AllResults         []DomainCheckResponse `json:"all_results"`
	Summary            ExtensionCheckSummary `json:"summary"`
}

// ExtensionCheckSummary provides a summary of the extension check
type ExtensionCheckSummary struct {
	PopularAvailable       []string `json:"popular_available"`       // Popular extensions that are available
	RecommendedDomains     []string `json:"recommended_domains"`     // Recommended domains to register
	AlternativeSuggestions []string `json:"alternative_suggestions"` // Alternative domain suggestions
	FastestResponse        *Domain  `json:"fastest_response"`        // Domain with fastest DNS response
	SlowestResponse        *Domain  `json:"slowest_response"`        // Domain with slowest DNS response
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta represents metadata for API responses
type Meta struct {
	Total       int    `json:"total,omitempty"`
	Page        int    `json:"page,omitempty"`
	PerPage     int    `json:"per_page,omitempty"`
	TotalPages  int    `json:"total_pages,omitempty"`
	RequestID   string `json:"request_id,omitempty"`
	ProcessTime int64  `json:"process_time_ms,omitempty"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status      string    `json:"status"`
	Version     string    `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
	Uptime      string    `json:"uptime"`
	Environment string    `json:"environment"`
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// WebSocketDomainCheck represents a domain check result for WebSocket
type WebSocketDomainCheck struct {
	Domain       string `json:"domain"`
	Status       string `json:"status"`
	IP           string `json:"ip,omitempty"`
	ResponseTime int64  `json:"response_time_ms"`
	CheckedAt    string `json:"checked_at"`
}

// WebSocketBulkProgress represents bulk check progress
type WebSocketBulkProgress struct {
	DomainName         string                 `json:"domain_name"`
	TotalExtensions    int                    `json:"total_extensions"`
	CheckedCount       int                    `json:"checked_count"`
	AvailableCount     int                    `json:"available_count"`
	UnavailableCount   int                    `json:"unavailable_count"`
	ErrorCount         int                    `json:"error_count"`
	CurrentDomain      *WebSocketDomainCheck  `json:"current_domain,omitempty"`
	AvailableDomains   []WebSocketDomainCheck `json:"available_domains"`
	UnavailableDomains []WebSocketDomainCheck `json:"unavailable_domains"`
	IsComplete         bool                   `json:"is_complete"`
	TotalTime          int64                  `json:"total_time_ms"`
}

// WhoisInfo represents WHOIS information for a domain
type WhoisInfo struct {
	Domain         string    `json:"domain"`
	Registrar      string    `json:"registrar,omitempty"`
	CreationDate   string    `json:"creation_date,omitempty"`
	ExpirationDate string    `json:"expiration_date,omitempty"`
	UpdatedDate    string    `json:"updated_date,omitempty"`
	NameServers    []string  `json:"name_servers,omitempty"`
	Status         []string  `json:"status,omitempty"`
	AdminContact   string    `json:"admin_contact,omitempty"`
	TechContact    string    `json:"tech_contact,omitempty"`
	RawData        string    `json:"raw_data,omitempty"`
	CheckedAt      time.Time `json:"checked_at"`
}
