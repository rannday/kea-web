package utils

import (
	"errors"
	"net"
)

// IsValidIP checks if a given string is a valid IP address and returns the IP or an error
func IsValidIP(ip string) (string, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "", errors.New("invalid IP address")
	}
	return ip, nil
}
