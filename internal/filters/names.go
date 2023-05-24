package filters

import (
	"fmt"
	"strings"

	"github.com/lusis/statusthing/internal/errors"
)

// WithName provides a custom name value
func WithName(name string) FilterOption {
	return func(f *Filters) error {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("name: %w", errors.ErrEmptyString)
		}
		if f.name != nil {
			return fmt.Errorf("name: %w", errors.ErrAlreadySet)
		}
		f.name = &name
		return nil
	}
}

// Name returns the custom name value
func (f *Filters) Name() string {
	f.l.RLock()
	defer f.l.RUnlock()

	return safeString(f.name)
}
