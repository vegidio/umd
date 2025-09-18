package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// HasHost checks if the host part of the given URL starts with the specified prefix.
//
// It returns true if the host ends with the suffix, otherwise false. If the URL is invalid, it prints an error message
// and returns false.
func HasHost(urlStr string, suffix string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Invalid URL:", err)
		return false
	}

	// Remove port if present
	host := parsedURL.Hostname()
	return strings.HasSuffix(host, suffix)
}
