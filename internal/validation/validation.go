// Package validation contains consistent validation helpers
package validation

import "strings"

// ValidString provides a shorter consistent check for what constitutes a valid string
// to reduce boilerplate
func ValidString(f string) bool {
	return strings.TrimSpace(f) != ""
}
