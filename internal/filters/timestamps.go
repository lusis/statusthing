package filters

import (
	"fmt"
	"time"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/serrors"

	"google.golang.org/protobuf/proto"
)

// Timestamps gets the [statusthingv1.Timestamps] that was provided
// when working with any protobuf message always use the getters
// and they can be safely chained
// i.e.
// `if f.Timestamps().GetCreated() == nil `
// or
// `mything := &statusthingv1.StatusThing{}; created := mything.GetTimestamps().GetCreated()`
// is the only reliable way to check a value
func (f *Filters) Timestamps() *statusthingv1.Timestamps {
	f.l.RLock()
	defer f.l.RUnlock()
	return f.timestamps
}

// WithTimestamps provides a custom [statusthingv1.Timestamps]
// Note that this is unconcerned with any individual timestamp value
// You are not required to set any timestamp field and this only checks that the
// current value is not equal to the existing one
func WithTimestamps(ts *statusthingv1.Timestamps) FilterOption {
	return func(f *Filters) error {
		if ts == nil {
			return fmt.Errorf("status kind: %w", serrors.ErrNilVal)
		}
		curval := f.timestamps

		if curval != nil && !proto.Equal(ts, curval) {
			return fmt.Errorf("timestamps: %w", serrors.ErrAlreadySet)
		}
		f.timestamps = ts
		return nil
	}
}

// WithLastLogin sets the last login
func (f *Filters) WithLastLogin(lastlogin *time.Time) FilterOption {
	return func(f *Filters) error {
		if lastlogin == nil {
			return serrors.NewError("lastlogin", serrors.ErrNilVal)
		}
		if f.lastlogin != nil {
			return serrors.NewError("lastlogin", serrors.ErrAlreadySet)
		}
		f.lastlogin = lastlogin
		return nil
	}
}

// LastLogin gets the last login
func (f *Filters) LastLogin() *time.Time {
	f.l.RLock()
	defer f.l.RUnlock()
	return f.lastlogin
}
