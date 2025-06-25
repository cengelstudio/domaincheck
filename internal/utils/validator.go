package utils

import (
	"regexp"
	"strings"
)

// ValidateDomainFormat validates domain name format
func ValidateDomainFormat(domain string) bool {
	if domain == "" {
		return false
	}

	// Remove protocol if present
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "www.")

	// Domain regex pattern
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)

	return domainRegex.MatchString(domain)
}

// ExtractDomainParts extracts domain name and extension
func ExtractDomainParts(domain string) (string, string) {
	// Clean domain
	domain = strings.ToLower(strings.TrimSpace(domain))
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "www.")

	// Find last dot
	lastDot := strings.LastIndex(domain, ".")
	if lastDot == -1 {
		return domain, ""
	}

	name := domain[:lastDot]
	extension := domain[lastDot:]

	return name, extension
}

// SanitizeDomain cleans and formats domain input
func SanitizeDomain(domain string) string {
	// Remove whitespace and convert to lowercase
	domain = strings.ToLower(strings.TrimSpace(domain))

	// Remove common prefixes
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "www.")

	// Remove trailing slash
	domain = strings.TrimSuffix(domain, "/")

	return domain
}
