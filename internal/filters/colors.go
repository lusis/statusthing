package filters

import (
	"fmt"
	"strings"

	"github.com/lusis/statusthing/internal/errors"
)

// WithColor provides a custom color value
func WithColor(color string) FilterOption {
	return func(f *Filters) error {
		if strings.TrimSpace(color) == "" {
			return fmt.Errorf("thing id %w", errors.ErrEmptyString)
		}
		if f.color != nil {
			return fmt.Errorf("thing id: %w", errors.ErrAlreadySet)
		}
		f.color = &color
		return nil
	}
}

// Color returns the custom color value
func (f *Filters) Color() string {
	f.l.RLock()
	defer f.l.RUnlock()
	return safeString(f.color)
}
