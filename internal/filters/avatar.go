package filters

import (
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/validation"
)

// WithAvatarURL provides a custom avatar url
func WithAvatarURL(aurl string) FilterOption {
	return func(f *Filters) error {
		if !validation.ValidString(aurl) {
			return serrors.NewError("aurl", serrors.ErrEmptyString)
		}
		if f.avatarURL != nil {
			return serrors.NewError("lname", serrors.ErrAlreadySet)
		}
		f.avatarURL = &aurl
		return nil
	}
}

// AvatarURL returns the configured avatar url
func (f *Filters) AvatarURL() string {
	f.l.RLock()
	defer f.l.RUnlock()
	if f.avatarURL == nil {
		return ""
	}
	return *f.avatarURL
}
