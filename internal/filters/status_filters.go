package filters

import (
	"fmt"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/serrors"
)

// WithStatus provides a custom [statusthingv1.Status]
func WithStatus(s *statusthingv1.Status) FilterOption {
	return func(f *Filters) error {
		if s == nil {
			return fmt.Errorf("status: %w", serrors.ErrNilVal)
		}
		if f.status != nil {
			return fmt.Errorf("status: %w", serrors.ErrAlreadySet)
		}
		if f.statusID != nil {
			return fmt.Errorf("status id already set: %w", serrors.ErrAlreadySet)
		}
		f.status = s
		return nil
	}
}

// Status returns the [statusthingv1.Status] that was provided
func (f *Filters) Status() *statusthingv1.Status {
	f.l.RLock()
	defer f.l.RUnlock()
	return f.status
}
