package services

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"domaincheck/internal/config"
	"domaincheck/internal/models"
	"domaincheck/internal/utils"
)

// DomainService handles domain checking operations
type DomainService struct {
	cfg              *config.Config
	validExtensions  map[string]bool
	checkedDomains   []models.Domain
	mutex            sync.RWMutex
	extensionsMutex  sync.RWMutex
	domainIDCounter  int
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
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Return copy of slice
	history := make([]models.Domain, len(s.checkedDomains))
	copy(history, s.checkedDomains)
	return history
}

// ClearHistory clears domain check history
func (s *DomainService) ClearHistory() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.checkedDomains = make([]models.Domain, 0)
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

// CheckAllExtensions checks domain name with all available extensions
func (s *DomainService) CheckAllExtensions(ctx context.Context, domainName string) (*models.AllExtensionsCheckResult, error) {
	startTime := time.Now()

	// Sanitize domain name (remove any existing extension)
	domainName = utils.SanitizeDomain(domainName)

	// Remove any extension if already provided
	if strings.Contains(domainName, ".") {
		parts := strings.Split(domainName, ".")
		domainName = parts[0]
	}

	// Validate domain name format (without extension)
	if domainName == "" || !regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`).MatchString(domainName) {
		return nil, fmt.Errorf("invalid domain name format: %s", domainName)
	}

	// Get all valid extensions
	extensions := s.GetValidExtensions()

	// Create domains list
	domains := make([]string, 0, len(extensions))
	for _, ext := range extensions {
		domains = append(domains, domainName+ext)
	}

	// Check all domains concurrently
	results, err := s.CheckMultipleDomains(ctx, domains)
	if err != nil {
		return nil, fmt.Errorf("failed to check domains: %w", err)
	}

	// Convert pointer slice to value slice
	allResults := make([]models.DomainCheckResponse, len(results))
	for i, result := range results {
		allResults[i] = *result
	}

	// Analyze results
	result := &models.AllExtensionsCheckResult{
		DomainName:         domainName,
		TotalExtensions:    len(extensions),
		CheckedAt:         time.Now(),
		TotalTime:         time.Since(startTime).Milliseconds(),
		AvailableDomains:   make([]models.DomainCheckResponse, 0),
		UnavailableDomains: make([]models.DomainCheckResponse, 0),
		ErrorDomains:       make([]models.DomainCheckResponse, 0),
		AllResults:         allResults,
	}

	// Categorize results
	var fastestDomain, slowestDomain *models.Domain
	popularExtensions := []string{".com", ".net", ".org", ".io", ".co", ".app", ".dev"}
	popularAvailable := make([]string, 0)

	for _, domainResult := range results {
		if domainResult.Domain.Error != "" {
			result.ErrorDomains = append(result.ErrorDomains, *domainResult)
			result.ErrorCount++
		} else if domainResult.Domain.Available {
			result.AvailableDomains = append(result.AvailableDomains, *domainResult)
			result.AvailableCount++

			// Check if it's a popular extension
			for _, popularExt := range popularExtensions {
				if domainResult.Domain.Extension == popularExt {
					popularAvailable = append(popularAvailable, domainResult.Domain.Name)
					break
				}
			}
		} else {
			result.UnavailableDomains = append(result.UnavailableDomains, *domainResult)
			result.UnavailableCount++
		}

		// Track fastest and slowest responses
		if domainResult.Domain.Error == "" {
			if fastestDomain == nil || domainResult.Domain.ResponseTime < fastestDomain.ResponseTime {
				fastestDomain = domainResult.Domain
			}
			if slowestDomain == nil || domainResult.Domain.ResponseTime > slowestDomain.ResponseTime {
				slowestDomain = domainResult.Domain
			}
		}
	}

	// Generate recommendations
	recommendedDomains := make([]string, 0)
	alternativeSuggestions := make([]string, 0)

	// Add top available popular domains to recommendations
	for _, domain := range popularAvailable {
		if len(recommendedDomains) < 5 {
			recommendedDomains = append(recommendedDomains, domain)
		}
	}

	// Add other available domains to recommendations if needed
	for _, domainResult := range result.AvailableDomains {
		if len(recommendedDomains) < 10 {
			alreadyAdded := false
			for _, rec := range recommendedDomains {
				if rec == domainResult.Domain.Name {
					alreadyAdded = true
					break
				}
			}
			if !alreadyAdded {
				recommendedDomains = append(recommendedDomains, domainResult.Domain.Name)
			}
		}
	}

	// Generate alternative suggestions
	if len(result.AvailableDomains) < 5 {
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

		for _, suggestion := range suggestions {
			if len(alternativeSuggestions) < 5 {
				alternativeSuggestions = append(alternativeSuggestions, suggestion+".com")
			}
		}
	}

	// Build summary
	result.Summary = models.ExtensionCheckSummary{
		PopularAvailable:       popularAvailable,
		RecommendedDomains:     recommendedDomains,
		AlternativeSuggestions: alternativeSuggestions,
		FastestResponse:        fastestDomain,
		SlowestResponse:        slowestDomain,
	}

	return result, nil
}
