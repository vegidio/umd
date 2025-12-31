package shared

import urlpkg "net/url"

// IsValidURL checks if a given string is a valid URL. It verifies that the URL can be parsed and has both a scheme
// (e.g., "http", "https") and a host.
//
// # Parameters:
//   - url: The URL string to validate.
//
// # Returns:
//   - bool: true if the URL is valid (has scheme and host), false otherwise.
//
// # Example:
//
//	IsValidURL("https://example.com")         // returns true
//	IsValidURL("http://example.com/path?q=1") // returns true
//	IsValidURL("example.com")                 // returns false (no scheme)
//	IsValidURL("https://")                    // returns false (no host)
//	IsValidURL("not a url")                   // returns false
func IsValidURL(url string) bool {
	parsedUrl, err := urlpkg.Parse(url)
	if err != nil {
		return false
	}

	// A valid URL should have both a scheme (http, https, etc.) and a host
	return parsedUrl.Scheme != "" && parsedUrl.Host != ""
}

// GetHost extracts the hostname from a URL string. It parses the given URL and returns the hostname portion, excluding
// the port number if present.
//
// # Parameters:
//   - url: The URL string to parse. Should be a valid URL format (e.g., "https://example.com:8080/path").
//
// # Returns:
//   - string: The hostname extracted from the URL (e.g., "example.com").
//     Returns an empty string if the URL is invalid or cannot be parsed.
//
// # Example:
//
//	GetHost("https://example.com:8080/path") // returns "example.com"
//	GetHost("https://api.example.com/path")  // returns "api.example.com"
//	GetHost("invalid-url")
func GetHost(url string) string {
	parsedUrl, err := urlpkg.Parse(url)
	if err != nil {
		return ""
	}

	return parsedUrl.Hostname()
}
