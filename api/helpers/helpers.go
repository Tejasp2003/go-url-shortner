package helpers

import (
	"os"
	"strings"
)

// EnforceHTTP ensures that the URL starts with "http://".
// If the URL does not start with "http", it prepends "http://" to the URL.
func EnforceHTTP(url string) string {
	// Check if the URL starts with "http".
	// If not, prepend "http://".
	if len(url) < 4 || url[:4] != "http" {
		return "http://" + url
	}
	return url
}

// RemoveDomainError checks if the provided URL matches the DOMAIN environment variable.
// It normalizes the URL by removing the protocol and "www".
// It returns false if the URL matches the DOMAIN environment variable, true otherwise.
func RemoveDomainError(url string) bool {
	// Get the DOMAIN environment variable.
	domain := os.Getenv("DOMAIN")

	// Check if the URL exactly matches the DOMAIN environment variable.
	if url == domain {
		return false
	}

	// Remove "http://", "https://", and "www." from the URL.
	// Split the URL at the first "/" and take the first part (the domain).
	newURL := strings.TrimPrefix(url, "http://") //TrimPrefix removes the prefix from the string ex: "http://www.google.com" -> "www.google.com"
	newURL = strings.TrimPrefix(newURL, "https://")
	newURL = strings.TrimPrefix(newURL, "www.")
	newURL = strings.Split(newURL, "/")[0]

	// Check if the normalized URL matches the DOMAIN environment variable.
	return newURL != domain
}
