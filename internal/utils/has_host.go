package utils

import (
	"fmt"
	urlpkg "net/url"
	"strings"

	"github.com/samber/lo"
)

// HasHost checks if the given URL's hostname matches or ends with any of the provided hostnames. It parses the URL,
// extracts the hostname (excluding port), and performs suffix matching against the provided hostnames.
//
// # Parameters:
//   - url: The URL string to check. Should be a valid URL format (e.g., "https://example.com/path").
//   - hostnames: One or more hostname strings to match against. The function checks if the URL's
//     hostname ends with any of these values, allowing for subdomain matching.
//
// # Returns:
//   - bool: true if the URL's hostname ends with any of the provided hostnames, false otherwise.
//     Also returns false if the URL is invalid (an error message will be printed to stdout).
//
// # Example:
//
//	HasHost("https://example.com:8080/path", "example.com") // returns true
//	HasHost("https://api.example.com/path", "example.com") // returns true
//	HasHost("https://other.com/path", "example.com") // returns false
func HasHost(url string, hostnames ...string) bool {
	parsedUrl, err := urlpkg.Parse(url)
	if err != nil {
		fmt.Println("Invalid URL:", err)
		return false
	}

	// Remove port if present
	host := parsedUrl.Hostname()

	return lo.ContainsBy(hostnames, func(hostname string) bool {
		return strings.HasSuffix(host, hostname)
	})
}
