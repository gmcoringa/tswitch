package util

import "strings"

// Check if a string is blank
func IsBlank(value string) bool {
	return strings.TrimSpace(value) == ""
}

// Check if a string is not blank
func IsNotBlank(value string) bool {
	return !IsBlank(value)
}
