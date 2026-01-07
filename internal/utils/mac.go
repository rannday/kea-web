package utils

import (
	"errors"
	"regexp"
	"strings"
)

// IsValidMAC validates a MAC address and returns it in an unformatted form (AAAAAAAAAAAA).
func IsValidMAC(mac string) (string, error) {
	// Strict format validation for common MAC address formats
	validMACRegex := regexp.MustCompile(`^(?:[0-9A-Fa-f]{4}\.[0-9A-Fa-f]{4}\.[0-9A-Fa-f]{4}|[0-9A-Fa-f]{2}[:-]?[0-9A-Fa-f]{2}[:-]?[0-9A-Fa-f]{2}[:-]?[0-9A-Fa-f]{2}[:-]?[0-9A-Fa-f]{2}[:-]?[0-9A-Fa-f]{2}|[0-9A-Fa-f]{12})$`)

	// Check if the input matches valid MAC formats
	if !validMACRegex.MatchString(mac) {
		return "", errors.New("invalid MAC address format")
	}

	// Use CleanMAC to normalize the MAC address
	cleanMac := CleanMAC(mac)

	// Enforce strict length of 12 characters after sanitization
	if len(cleanMac) != 12 {
		return "", errors.New("invalid MAC address length")
	}

	return cleanMac, nil
}

// CleanMAC removes separators and ensures uppercase (for SQL storage)
func CleanMAC(mac string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(mac, ":", ""), "-", ""), ".", ""))
}

// FormatMAC ensures the MAC address is in aa:aa:aa:aa:aa:aa format
func FormatMAC(mac string) string {
	// Check if it's already correctly formatted
	validMACRegex := regexp.MustCompile(`^(?:[0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)
	if validMACRegex.MatchString(mac) {
		return strings.ToLower(mac) // Ensure lowercase
	}

	// Use CleanMAC to remove invalid characters
	cleanMac := CleanMAC(mac)
	if len(cleanMac) != 12 {
		return "" // Invalid MAC length after stripping
	}

	// Format into aa:aa:aa:aa:aa:aa
	return strings.ToLower(strings.Join([]string{
		cleanMac[0:2], cleanMac[2:4], cleanMac[4:6], cleanMac[6:8], cleanMac[8:10], cleanMac[10:12],
	}, ":"))
}
