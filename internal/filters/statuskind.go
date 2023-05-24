package filters

import (
	"fmt"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/serrors"
)

// StatusKind gets the [statusthingv1.StatusKind] that was provided
func (f *Filters) StatusKind() statusthingv1.StatusKind {
	f.l.RLock()
	defer f.l.RUnlock()
	return f.statusKind
}

// WithStatusKind provides a custom [statusthingv1.StatusKind]
func WithStatusKind(k statusthingv1.StatusKind) FilterOption {
	return func(f *Filters) error {
		if k == statusthingv1.StatusKind_STATUS_KIND_UNKNOWN {
			return fmt.Errorf("statusKind: %w", serrors.ErrEmptyEnum)
		}
		curval := f.statusKind
		if curval != statusthingv1.StatusKind_STATUS_KIND_UNKNOWN &&
			curval != k {
			return fmt.Errorf("statusKind: %w", serrors.ErrAlreadySet)
		}
		f.statusKind = k
		return nil
	}
}

// StatusKinds gets the slice of [statusthingv1.StatusKind] that was provided
func (f *Filters) StatusKinds() []statusthingv1.StatusKind {
	f.l.RLock()
	defer f.l.RUnlock()
	return f.statusKinds
}

// WithStatusKinds provides a custom slice of [statusthingv1.StatusKind]
func WithStatusKinds(kinds ...statusthingv1.StatusKind) FilterOption {
	return func(f *Filters) error {
		if len(kinds) == 0 {
			return fmt.Errorf("kinds: %w", serrors.ErrAtLeastOne)
		}
		// check for invalid values
		for _, k := range kinds {
			if k == statusthingv1.StatusKind_STATUS_KIND_UNKNOWN {
				return fmt.Errorf("status kind: %w", serrors.ErrEmptyEnum)
			}
		}

		if f.statusKinds != nil {
			return fmt.Errorf("statusKinds: %w", serrors.ErrAlreadySet)
		}
		f.statusKinds = kinds
		return nil
	}
}
