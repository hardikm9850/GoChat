package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
)

// HashPhoneNumber creates a deterministic hash for contact matching
func HashPhoneNumber(countryCode, phone string) string {
	normalized := countryCode + strings.ReplaceAll(phone, " ", "")
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:])
}

// NormalizePhoneNumber ensures consistent format
func NormalizePhoneNumber(countryCode, phoneNumber string) string {
    // Remove any spaces, dashes, parentheses
    cleaned := regexp.MustCompile(`[^0-9+]`).ReplaceAllString(phoneNumber, "")
    return countryCode + cleaned
}