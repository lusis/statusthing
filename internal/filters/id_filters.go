package filters

import (
	"fmt"
	"strings"

	"github.com/lusis/statusthing/internal/errors"
)

// WithItemID provides a custom [statusthingv1.StatusThing] id
func WithItemID(id string) FilterOption {
	return func(f *Filters) error {
		if strings.TrimSpace(id) == "" {
			return fmt.Errorf("item id %w", errors.ErrEmptyString)
		}
		if f.itemID != nil {
			return fmt.Errorf("item id: %w", errors.ErrAlreadySet)
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
			return fmt.Errorf("status id %w", errors.ErrEmptyString)
		}
		if f.statusID != nil {
			return fmt.Errorf("status id: %w", errors.ErrAlreadySet)
		}
		if f.status != nil {
			return fmt.Errorf("status already set: %w", errors.ErrAlreadySet)
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
			return fmt.Errorf("note id: %w", errors.ErrEmptyString)
		}
		if f.noteID != nil {
			return fmt.Errorf("note id: %w", errors.ErrAlreadySet)
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
			return fmt.Errorf("statusIDs %w", errors.ErrAtLeastOne)
		}
		if len(f.statusIDs) != 0 {
			return fmt.Errorf("statusIDs: %w", errors.ErrAlreadySet)
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
