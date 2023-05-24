package filters

import (
	"fmt"
	"strings"

	"github.com/lusis/statusthing/internal/errors"
)

// WithDescription provides a custom description
func WithDescription(description string) FilterOption {
	return func(f *Filters) error {
		if strings.TrimSpace(description) == "" {
			return fmt.Errorf("description: %w", errors.ErrEmptyString)
		}
		if f.description != nil {
			return fmt.Errorf("description: %w", errors.ErrAlreadySet)
		}
		f.description = &description
		return nil
	}
}

// Description returns the custom description value
func (f *Filters) Description() string {
	f.l.RLock()
	defer f.l.RUnlock()
	return safeString(f.description)
}
