package utils

import (
	"encoding/hex"
	"errors"
	"net"
	"strings"
)

// IsValidIPv4 validates IPv4 and returns canonical dotted-quad.
func IsValidIPv4(s string) (string, error) {
	ip := net.ParseIP(strings.TrimSpace(s))
	if ip == nil {
		return "", errors.New("invalid IP address")
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return "", errors.New("not an IPv4 address")
	}
	return ip4.String(), nil
}

// IsValidIPv6 validates IPv6 and returns canonical form.
// Rejects IPv4-mapped/compatible addresses (anything that To4() would treat as IPv4).
// Also rejects zone identifiers (e.g. "fe80::1%eth0") by design.
func IsValidIPv6(s string) (string, error) {
	s = strings.TrimSpace(s)
	if strings.Contains(s, "%") {
		return "", errors.New("invalid IPv6 address")
	}

	ip := net.ParseIP(s)
	if ip == nil {
		return "", errors.New("invalid IPv6 address")
	}
	if ip.To4() != nil {
		return "", errors.New("not an IPv6 address")
	}

	return ip.String(), nil
}

// MACPolicy controls additional policy checks beyond basic parsing.
type MACPolicy struct {
	AllowAllZero         bool // 00:00:00:00:00:00
	AllowBroadcast       bool // FF:FF:FF:FF:FF:FF
	AllowMulticast       bool // I/G bit set (LSB of first octet)
	AllowLocallyAdmin    bool // U/L bit set (second LSB of first octet)
}

// DefaultMACPolicy is a sane “real device” policy: rejects the common bogus cases.
var DefaultMACPolicy = MACPolicy{
	AllowAllZero:      false,
	AllowBroadcast:    false,
	AllowMulticast:    false,
	AllowLocallyAdmin: true, // set false if you want only globally-administered OUIs
}

// IsValidMAC validates a MAC-48 address and returns it as 12 uppercase hex chars.
func IsValidMAC(s string) (string, error) {
	return IsValidMACWithPolicy(s, DefaultMACPolicy)
}

// IsValidMACWithPolicy validates MAC-48 and applies additional policy checks.
func IsValidMACWithPolicy(s string, p MACPolicy) (string, error) {
	hw, err := net.ParseMAC(strings.TrimSpace(s))
	if err != nil {
		return "", errors.New("invalid MAC address")
	}
	// net.ParseMAC returns 6 bytes for MAC-48; EUI-64 would be 8 bytes.
	if len(hw) != 6 {
		return "", errors.New("invalid MAC length")
	}

	// Policy checks
	allZero := true
	allFF := true
	for _, b := range hw {
		if b != 0x00 {
			allZero = false
		}
		if b != 0xFF {
			allFF = false
		}
	}

	if allZero && !p.AllowAllZero {
		return "", errors.New("invalid MAC address")
	}
	if allFF && !p.AllowBroadcast {
		return "", errors.New("invalid MAC address")
	}

	// Multicast if I/G bit is 1 (LSB of first octet)
	isMulticast := (hw[0] & 0x01) != 0
	if isMulticast && !p.AllowMulticast {
		return "", errors.New("invalid MAC address")
	}

	// Locally administered if U/L bit is 1 (second LSB of first octet)
	isLocallyAdmin := (hw[0] & 0x02) != 0
	if isLocallyAdmin && !p.AllowLocallyAdmin {
		return "", errors.New("invalid MAC address")
	}

	// 12 hex uppercase, no separators
	return strings.ToUpper(hex.EncodeToString(hw)), nil
}

// FormatMAC returns aa:aa:aa:aa:aa:aa (lowercase) for a valid MAC.
func FormatMAC(s string) (string, error) {
	clean, err := IsValidMAC(s)
	if err != nil {
		return "", err
	}
	return strings.ToLower(strings.Join([]string{
		clean[0:2], clean[2:4], clean[4:6], clean[6:8], clean[8:10], clean[10:12],
	}, ":")), nil
}
