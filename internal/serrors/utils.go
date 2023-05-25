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

// NewError returns a statusthing specific error with a consistent error message
func NewError(field string, err error) error {
	return fmt.Errorf("%s: %w", safeField(field), err)
}

// NewWrappedError returns an serror wrapping an original error
func NewWrappedError(field string, statusThingErr error, originalErr error) error {
	return fmt.Errorf("%s: %w (%w)", field, statusThingErr, originalErr)
}
