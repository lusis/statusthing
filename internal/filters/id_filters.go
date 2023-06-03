package filters

import (
	"fmt"
	"strings"

	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/validation"
)

// WithItemID provides a custom [statusthingv1.StatusThing] id
func WithItemID(id string) FilterOption {
	return func(f *Filters) error {
		if strings.TrimSpace(id) == "" {
			return fmt.Errorf("item id %w", serrors.ErrEmptyString)
		}
		if f.itemID != nil {
			return fmt.Errorf("item id: %w", serrors.ErrAlreadySet)
		}
		f.itemID = &id
		return nil
	}
}

// ItemID returns the configured [statusthingv1.StatusThing] id
func (f *Filters) ItemID() string {
	f.l.RLock()
	defer f.l.RUnlock()
	return safeString(f.itemID)
}

// WithStatusID provides a custom [statusthingv1.CustomStatus] id
func WithStatusID(id string) FilterOption {
	return func(f *Filters) error {
		if strings.TrimSpace(id) == "" {
			return fmt.Errorf("status id %w", serrors.ErrEmptyString)
		}
		if f.statusID != nil {
			return fmt.Errorf("status id: %w", serrors.ErrAlreadySet)
		}
		if f.status != nil {
			return fmt.Errorf("status already set: %w", serrors.ErrAlreadySet)
		}
		f.statusID = &id
		return nil
	}
}

// StatusID returns the configured [statusthingv1.CustomStatus] id
func (f *Filters) StatusID() string {
	f.l.RLock()
	defer f.l.RUnlock()
	return safeString(f.statusID)
}

// WithNoteID provides a custom [statusthingv1.Note] id
func WithNoteID(id string) FilterOption {
	return func(f *Filters) error {
		if strings.TrimSpace(id) == "" {
			return fmt.Errorf("note id: %w", serrors.ErrEmptyString)
		}
		if f.noteID != nil {
			return fmt.Errorf("note id: %w", serrors.ErrAlreadySet)
		}
		f.noteID = &id
		return nil
	}
}

// NoteID returns the configured [statusthingv1.Note] id
func (f *Filters) NoteID() string {
	f.l.RLock()
	defer f.l.RUnlock()
	if f.noteID == nil {
		return ""
	}
	return *f.noteID
}

// WithStatusIDs provides a custom slice of [statusthingv1.CustomStatus] ids
func WithStatusIDs(statusIDs ...string) FilterOption {
	return func(f *Filters) error {
		if len(statusIDs) == 0 {
			return fmt.Errorf("statusIDs %w", serrors.ErrAtLeastOne)
		}
		if len(f.statusIDs) != 0 {
			return fmt.Errorf("statusIDs: %w", serrors.ErrAlreadySet)
		}
		f.statusIDs = statusIDs

		return nil
	}
}

// StatusIDs returns the configured slice of [statusthingv1.CustomStatus] ids
func (f *Filters) StatusIDs() []string {
	f.l.RLock()
	defer f.l.RUnlock()
	return f.statusIDs
}

// WithUserID provides a custom [v1.User] id
func WithUserID(id string) FilterOption {
	return func(f *Filters) error {
		if !validation.ValidString(id) {
			return serrors.NewError("userid", serrors.ErrEmptyString)
		}
		if f.userid != nil {
			return serrors.NewError("userid", serrors.ErrAlreadySet)
		}
		f.userid = &id
		return nil
	}
}

// UserID returns the configured [v1.User] id
func (f *Filters) UserID() string {
	f.l.RLock()
	defer f.l.RUnlock()
	if f.userid == nil {
		return ""
	}
	return *f.userid
}
