package filters

import (
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/validation"
)

// WithPassword provides a password to a filter
func WithPassword(password string) FilterOption {
	return func(f *Filters) error {
		if !validation.ValidString(password) {
			return serrors.NewError("password", serrors.ErrEmptyString)
		}
		if f.password != nil {
			return serrors.NewError("password", serrors.ErrAlreadySet)
		}
		f.password = &password
		return nil
	}
}

// Password gets the set password
func (f *Filters) Password() string {
	f.l.RLock()
	defer f.l.RUnlock()
	return safeString(f.password)
}
