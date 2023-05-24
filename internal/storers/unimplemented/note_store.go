package unimplemented

import (
	"context"

	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/errors"
	"github.com/lusis/statusthing/internal/filters"
)

// NoteStorer stores [statusthingv1.Note]
type NoteStorer struct{}

// GetNote gets a [statusthingv1.Note] by its id
func (ns *NoteStorer) GetNote(ctx context.Context, noteID string) (*v1.Note, error) { // nolint: revive
	return nil, errors.ErrNotImplemented
}

// StoreNote stores the provided [statusthingv1.Note] associated with the provided [statusthingv1.StatusThing] by its id
func (ns *NoteStorer) StoreNote(ctx context.Context, note *v1.Note, statusThingID string) (*v1.Note, error) { // nolint: revive
	return nil, errors.ErrNotImplemented
}

// FindNotes gets all known [statusthingv1.Note]
func (ns *NoteStorer) FindNotes(ctx context.Context, itemID string, opts ...filters.FilterOption) ([]*v1.Note, error) { // nolint: revive
	return nil, errors.ErrNotImplemented
}

// UpdateNote updates the [statusthingv1.Note] with the provided [filters.FilterOption]
func (ns *NoteStorer) UpdateNote(ctx context.Context, noteID string, opts ...filters.FilterOption) error { // nolint: revive
	return errors.ErrNotImplemented
}

// DeleteNote deletes a [statusthingv1.Note] by its id
func (ns *NoteStorer) DeleteNote(ctx context.Context, noteID string) error { // nolint: revive
	return errors.ErrNotImplemented
}
