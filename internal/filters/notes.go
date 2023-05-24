package filters

import (
	"github.com/lusis/statusthing/internal/errors"

	"fmt"
	"strings"
)

// NoteText gets the value of [Filters.noteText]
func (f *Filters) NoteText() string {
	f.l.RLock()
	defer f.l.RUnlock()
	return safeString(f.noteText)
}

// WithNoteText provides a custom note text for things like updates
func WithNoteText(text string) FilterOption {
	return func(f *Filters) error {
		if strings.TrimSpace(text) == "" {
			return fmt.Errorf("text: %w", errors.ErrEmptyString)
		}
		if f.noteText != nil {
			return fmt.Errorf("noteText: %w", errors.ErrAlreadySet)
		}
		f.noteText = &text
		return nil
	}
}
