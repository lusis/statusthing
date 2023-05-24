package filters

import (
	"fmt"
	"strings"

	"github.com/lusis/statusthing/internal/serrors"
)

// WithColor provides a custom color value
func WithColor(color string) FilterOption {
	return func(f *Filters) error {
		if strings.TrimSpace(color) == "" {
			return fmt.Errorf("thing id %w", serrors.ErrEmptyString)
		}
		if f.color != nil {
			return fmt.Errorf("thing id: %w", serrors.ErrAlreadySet)
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
