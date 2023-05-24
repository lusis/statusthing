package serrors

import (
	"fmt"
	"strings"
)

func safeField(f string) string {
	if strings.TrimSpace(f) == "" {
		f = "field-not-provided"
	}
	return f
}

// ValidString provides a shorter consistent check for what constitutes a valid string
// to reduce boilerplate
func ValidString(f string) bool {
	return strings.TrimSpace(f) != ""
}

// NewError returns a statusthing specific error with a consistent error message
func NewError(field string, err error) error {
	return fmt.Errorf("%s: %w", safeField(field), err)
}
