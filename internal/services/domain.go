package services

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"domaincheck/internal/config"
	"domaincheck/internal/models"
	"domaincheck/internal/utils"
)

// DomainService handles domain checking operations
type DomainService struct {
	cfg             *config.Config
	validExtensions map[string]bool
	checkedDomains  []models.Domain
	mutex           sync.RWMutex
	extensionsMutex sync.RWMutex
	domainIDCounter int
	history         []models.Domain
	historyMutex    sync.RWMutex
}

// NewDomainService creates a new domain service instance
func NewDomainService(cfg *config.Config) (*DomainService, error) {
	service := &DomainService{
		cfg:             cfg,
		validExtensions: make(map[string]bool),
		checkedDomains:  make([]models.Domain, 0),
		domainIDCounter: 1,
	}

	// Load valid extensions from file
	if err := service.loadValidExtensions(); err != nil {
		return nil, fmt.Errorf("failed to load domain extensions: %w", err)
	}

	return service, nil
}

// loadValidExtensions loads valid domain extensions from file
func (s *DomainService) loadValidExtensions() error {
	s.extensionsMutex.Lock()
	defer s.extensionsMutex.Unlock()

	file, err := os.Open(s.cfg.Domain.ExtensionsFile)
	if err != nil {
		return fmt.Errorf("failed to open extensions file: %w", err)
	}
	defer file.Close()

	// Clear existing extensions
	s.validExtensions = make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		extension := strings.TrimSpace(scanner.Text())
		if extension != "" {
			// Ensure extension starts with dot
			if !strings.HasPrefix(extension, ".") {
				extension = "." + extension
			}
			s.validExtensions[strings.ToLower(extension)] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading extensions file: %w", err)
	}

	return nil
}

// ReloadExtensions reloads domain extensions from file
func (s *DomainService) ReloadExtensions() error {
	return s.loadValidExtensions()
}

// GetValidExtensions returns all valid extensions
func (s *DomainService) GetValidExtensions() []string {
	s.extensionsMutex.RLock()
	defer s.extensionsMutex.RUnlock()

	extensions := make([]string, 0, len(s.validExtensions))
	for ext := range s.validExtensions {
		extensions = append(extensions, ext)
	}
	return extensions
}

// IsValidExtension checks if extension is valid
func (s *DomainService) IsValidExtension(extension string) bool {
	s.extensionsMutex.RLock()
	defer s.extensionsMutex.RUnlock()

	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	return s.validExtensions[strings.ToLower(extension)]
}

// CheckDomain performs domain availability check
func (s *DomainService) CheckDomain(ctx context.Context, domainName string) (*models.DomainCheckResponse, error) {
	startTime := time.Now()

	// Sanitize domain
	domainName = utils.SanitizeDomain(domainName)

	// Validate domain format
	if !utils.ValidateDomainFormat(domainName) {
		return nil, fmt.Errorf("invalid domain format: %s", domainName)
	}

	// Extract domain parts
	_, extension := utils.ExtractDomainParts(domainName)
	if extension == "" {
		return nil, fmt.Errorf("domain must have an extension")
	}

	// Check if extension is supported
	isValidTLD := s.IsValidExtension(extension)

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, s.cfg.Domain.Timeout)
	defer cancel()

	// Perform DNS lookup
	domain := &models.Domain{
		ID:        s.getNextDomainID(),
		Name:      domainName,
		Extension: extension,
		CheckedAt: time.Now(),
	}

	// DNS resolution with multiple DNS servers for reliability
	dnsServers := []string{"8.8.8.8:53", "1.1.1.1:53", "208.67.222.222:53"} // Google, Cloudflare, OpenDNS

	var ips []net.IPAddr
	var err error
	var lastError error

	// Try different DNS servers
	for _, dnsServer := range dnsServers {
		resolver := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Second * 2,
				}
				return d.DialContext(ctx, network, dnsServer)
			},
		}

		ips, err = resolver.LookupIPAddr(timeoutCtx, domainName)
		if err == nil {
			break // Success, exit loop
		}
		lastError = err

		// If context is cancelled or timed out, don't try other servers
		if timeoutCtx.Err() != nil {
			break
		}
	}

	if err != nil {
		// Check if it's a real DNS error (domain doesn't exist) or network error
		errorStr := err.Error()
		if strings.Contains(errorStr, "no such host") ||
			strings.Contains(errorStr, "NXDOMAIN") ||
			strings.Contains(errorStr, "server misbehaving") {
			domain.Available = true
			domain.DNSResolved = false
			domain.Status = "Available"
		} else {
			// Network or other error
			domain.Available = false
			domain.DNSResolved = false
			domain.Status = "Error"
		}

		// Use the last error for reporting
		if lastError != nil {
			domain.Error = lastError.Error()
		} else {
			domain.Error = err.Error()
		}
	} else {
		domain.Available = false
		domain.DNSResolved = true
		domain.Status = "Registered"
		if len(ips) > 0 {
			domain.IP = ips[0].IP.String()
		}
	}

	// Calculate response time
	domain.ResponseTime = time.Since(startTime).Milliseconds()

	// Store domain check result
	s.storeDomainResult(*domain)

	// Add to history
	s.AddToHistory(domain)

	response := &models.DomainCheckResponse{
		Domain:       domain,
		IsValidTLD:   isValidTLD,
		SupportedTLD: isValidTLD,
	}

	return response, nil
}

// CheckMultipleDomains checks multiple domains concurrently
func (s *DomainService) CheckMultipleDomains(ctx context.Context, domains []string) ([]*models.DomainCheckResponse, error) {
	if len(domains) == 0 {
		return []*models.DomainCheckResponse{}, nil
	}

	// Limit concurrent checks
	maxConcurrent := s.cfg.Domain.MaxConcurrentChecks
	if len(domains) < maxConcurrent {
		maxConcurrent = len(domains)
	}

	// Create channels for coordination
	domainChan := make(chan string, len(domains))
	resultChan := make(chan *models.DomainCheckResponse, len(domains))
	errorChan := make(chan error, len(domains))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < maxConcurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for domain := range domainChan {
				result, err := s.CheckDomain(ctx, domain)
				if err != nil {
					errorChan <- fmt.Errorf("failed to check %s: %w", domain, err)
					continue
				}
				resultChan <- result
			}
		}()
	}

	// Send domains to workers
	go func() {
		defer close(domainChan)
		for _, domain := range domains {
			domainChan <- domain
		}
	}()

	// Close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Collect results
	var results []*models.DomainCheckResponse
	var errors []error

	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				resultChan = nil
			} else {
				results = append(results, result)
			}
		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
			} else {
				errors = append(errors, err)
			}
		}

		if resultChan == nil && errorChan == nil {
			break
		}
	}

	// Return first error if any
	if len(errors) > 0 {
		return results, errors[0]
	}

	return results, nil
}

// GetDomainHistory returns checked domain history
func (s *DomainService) GetDomainHistory() []models.Domain {
	s.historyMutex.RLock()
	defer s.historyMutex.RUnlock()

	// Return a copy to avoid race conditions
	result := make([]models.Domain, len(s.history))
	copy(result, s.history)
	return result
}

// ClearHistory clears domain check history
func (s *DomainService) ClearHistory() {
	s.historyMutex.Lock()
	defer s.historyMutex.Unlock()
	s.history = []models.Domain{}
}

// storeDomainResult stores domain check result
func (s *DomainService) storeDomainResult(domain models.Domain) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.checkedDomains = append(s.checkedDomains, domain)
}

// getNextDomainID returns next domain ID
func (s *DomainService) getNextDomainID() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := s.domainIDCounter
	s.domainIDCounter++
	return id
}

// CheckAllExtensions checks a domain name with all available extensions
func (s *DomainService) CheckAllExtensions(ctx context.Context, domainName string) (*models.AllExtensionsCheckResult, error) {
	startTime := time.Now()

	// Get all extensions
	extensions := s.GetValidExtensions()
	totalExtensions := len(extensions)

	// Create result structure
	result := &models.AllExtensionsCheckResult{
		DomainName:         domainName,
		TotalExtensions:    totalExtensions,
		AvailableCount:     0,
		UnavailableCount:   0,
		ErrorCount:         0,
		TotalTime:          0,
		AllResults:         []models.DomainCheckResponse{},
		AvailableDomains:   []models.DomainCheckResponse{},
		UnavailableDomains: []models.DomainCheckResponse{},
		Summary: models.ExtensionCheckSummary{
			RecommendedDomains:     []string{},
			AlternativeSuggestions: []string{},
		},
	}

	// Create channels for concurrent processing
	results := make(chan models.DomainCheckResponse, totalExtensions)
	errors := make(chan error, totalExtensions)

	// Start concurrent domain checks
	semaphore := make(chan struct{}, s.cfg.Domain.MaxConcurrentChecks)
	var wg sync.WaitGroup

	for _, ext := range extensions {
		wg.Add(1)
		go func(extension string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			fullDomain := domainName + "." + extension
			start := time.Now()

			checkResult, err := s.CheckDomain(ctx, fullDomain)
			responseTime := time.Since(start).Milliseconds()

			if err != nil {
				errors <- err
				return
			}

			// Update response time
			checkResult.Domain.ResponseTime = responseTime

			// Add to history
			s.AddToHistory(checkResult.Domain)

			results <- *checkResult
		}(ext)
	}

	// Close channels when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Process results
	for resultItem := range results {
		result.AllResults = append(result.AllResults, resultItem)

		if resultItem.Domain.Status == "Available" {
			result.AvailableCount++
			result.AvailableDomains = append(result.AvailableDomains, resultItem)
			result.Summary.RecommendedDomains = append(result.Summary.RecommendedDomains, resultItem.Domain.Name)
		} else {
			result.UnavailableCount++
			result.UnavailableDomains = append(result.UnavailableDomains, resultItem)
		}
	}

	// Process errors
	for range errors {
		result.ErrorCount++
	}

	// Calculate total time
	result.TotalTime = time.Since(startTime).Milliseconds()

	// Generate alternative suggestions
	result.Summary.AlternativeSuggestions = s.generateAlternativeSuggestions(domainName)

	return result, nil
}

// GetWhoisInfo retrieves WHOIS information for a domain
func (s *DomainService) GetWhoisInfo(ctx context.Context, domain string) (*models.WhoisInfo, error) {
	// Simple WHOIS lookup using net package
	// Note: This is a basic implementation. For production, consider using a dedicated WHOIS library
	whoisInfo := &models.WhoisInfo{
		Domain:    domain,
		CheckedAt: time.Now(),
	}

	// Try to get basic domain information
	ips, err := net.LookupIP(domain)
	if err == nil && len(ips) > 0 {
		whoisInfo.Status = []string{"Active"}
	}

	// Get nameservers
	ns, err := net.LookupNS(domain)
	if err == nil {
		for _, n := range ns {
			whoisInfo.NameServers = append(whoisInfo.NameServers, n.Host)
		}
	}

	// For a more comprehensive WHOIS lookup, you would need to:
	// 1. Connect to the appropriate WHOIS server for the TLD
	// 2. Send WHOIS query
	// 3. Parse the response
	// This is a simplified version for demonstration

	return whoisInfo, nil
}

// AddToHistory adds a domain check result to history
func (s *DomainService) AddToHistory(domain *models.Domain) {
	s.historyMutex.Lock()
	defer s.historyMutex.Unlock()

	// Add to beginning of slice
	s.history = append([]models.Domain{*domain}, s.history...)

	// Keep only last 1000 entries
	if len(s.history) > 1000 {
		s.history = s.history[:1000]
	}
}

// generateAlternativeSuggestions generates alternative domain suggestions
func (s *DomainService) generateAlternativeSuggestions(domainName string) []string {
	suggestions := []string{
		domainName + "app",
		domainName + "pro",
		domainName + "online",
		domainName + "digital",
		domainName + "tech",
		"get" + domainName,
		"my" + domainName,
		domainName + "hub",
		domainName + "zone",
		domainName + "lab",
	}

	// Return first 5 suggestions
	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	// Add .com extension
	for i, suggestion := range suggestions {
		suggestions[i] = suggestion + ".com"
	}

	return suggestions
}
