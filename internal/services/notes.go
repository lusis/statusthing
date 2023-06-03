package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/segmentio/ksuid"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/validation"
)

// AddNote adds the provided text as a [statusthingv1.Note] to the [statusthingv1.Item] with the provided id
func (sts *StatusThingService) AddNote(ctx context.Context, itemID, noteText string) (*statusthingv1.Note, error) {
	if sts.store == nil {
		return nil, fmt.Errorf("store was nil: %w", serrors.ErrStoreUnavailable)
	}
	if !validation.ValidString(itemID) {
		return nil, serrors.NewError("itemID", serrors.ErrEmptyString)
	}
	if !validation.ValidString(noteText) {
		return nil, serrors.NewError("noteText", serrors.ErrEmptyString)
	}
	if _, err := sts.GetItem(ctx, itemID); err != nil {
		return nil, err
	}
	id := ksuid.New().String()
	note := &statusthingv1.Note{
		Timestamps: makeTsNow(),
		Id:         id,
		Text:       noteText,
	}
	return sts.store.StoreNote(ctx, note, itemID)
}

// EditNote edits the [statusthingv1.Note] with provided id to set the text to provided text
// supported opts:
// [filters.WithTimestamps] to override the timestamps (generally for testing)
func (sts *StatusThingService) EditNote(ctx context.Context, noteID, noteText string, opts ...filters.FilterOption) error {
	if sts.store == nil {
		return fmt.Errorf("store was nil: %w", serrors.ErrStoreUnavailable)
	}
	if strings.TrimSpace(noteID) == "" {
		return fmt.Errorf("noteID: %w", serrors.ErrEmptyString)
	}
	if strings.TrimSpace(noteText) == "" {
		return fmt.Errorf("noteText: %w", serrors.ErrEmptyString)
	}
	f, err := filters.New(opts...)
	if err != nil {
		return err
	}

	allowedOpts := []filters.FilterOption{
		filters.WithNoteText(noteText),
	}
	if f.Timestamps() != nil {
		allowedOpts = append(allowedOpts, filters.WithTimestamps(f.Timestamps()))
	}

	return sts.store.UpdateNote(ctx, noteID, allowedOpts...)
}

// RemoveNote removes the [statusthingv1.Note] with the provided id from the [statusthingv1.Item] with the provided id
func (sts *StatusThingService) RemoveNote(ctx context.Context, noteID string) error {
	if sts.store == nil {
		return fmt.Errorf("store was nil: %w", serrors.ErrStoreUnavailable)
	}
	if strings.TrimSpace(noteID) == "" {
		return fmt.Errorf("noteID: %w", serrors.ErrEmptyString)
	}
	return sts.store.DeleteNote(ctx, noteID)
}

// FindNotes returns all [statusthingv1.Note] belonging to [statusthingv1.Item] with provided id
func (sts *StatusThingService) FindNotes(ctx context.Context, itemID string) ([]*statusthingv1.Note, error) {
	if sts.store == nil {
		return nil, serrors.NewError("store", serrors.ErrStoreUnavailable)
	}
	if !validation.ValidString(itemID) {
		return nil, serrors.NewError("itemID", serrors.ErrEmptyString)
	}
	if _, err := sts.GetItem(ctx, itemID); err != nil {
		return nil, err
	}
	return sts.store.FindNotes(ctx, itemID)
}

// GetNote gets a [statusthingv1.Note] by id
func (sts *StatusThingService) GetNote(ctx context.Context, noteID string) (*statusthingv1.Note, error) {
	if sts.store == nil {
		return nil, fmt.Errorf("store was nil: %w", serrors.ErrStoreUnavailable)
	}

	if strings.TrimSpace(noteID) == "" {
		return nil, fmt.Errorf("noteID: %w", serrors.ErrEmptyString)
	}
	return sts.store.GetNote(ctx, noteID)
}
