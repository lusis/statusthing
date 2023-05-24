package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/segmentio/ksuid"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/errors"
	"github.com/lusis/statusthing/internal/filters"
)

// NewNote adds the provided text as a [statusthingv1.Note] to the [statusthingv1.Item] with the provided id
func (sts *StatusThingService) NewNote(ctx context.Context, itemID, noteText string) (*statusthingv1.Note, error) {
	if sts.store == nil {
		return nil, fmt.Errorf("store was nil: %w", errors.ErrStoreUnavailable)
	}
	if strings.TrimSpace(itemID) == "" {
		return nil, fmt.Errorf("itemID: %w", errors.ErrEmptyString)
	}
	if strings.TrimSpace(noteText) == "" {
		return nil, fmt.Errorf("noteText: %w", errors.ErrEmptyString)
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
		return fmt.Errorf("store was nil: %w", errors.ErrStoreUnavailable)
	}
	if strings.TrimSpace(noteID) == "" {
		return fmt.Errorf("noteID: %w", errors.ErrEmptyString)
	}
	if strings.TrimSpace(noteText) == "" {
		return fmt.Errorf("noteText: %w", errors.ErrEmptyString)
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

// DeleteNote removes the [statusthingv1.Note] with the provided id from the [statusthingv1.Item] with the provided id
func (sts *StatusThingService) DeleteNote(ctx context.Context, noteID string) error {
	if sts.store == nil {
		return fmt.Errorf("store was nil: %w", errors.ErrStoreUnavailable)
	}
	if strings.TrimSpace(noteID) == "" {
		return fmt.Errorf("noteID: %w", errors.ErrEmptyString)
	}
	return sts.store.DeleteNote(ctx, noteID)
}

// AllNotes returns all [statusthingv1.Note] belonging to [statusthingv1.Item] with provided id
func (sts *StatusThingService) AllNotes(ctx context.Context, itemID string) ([]*statusthingv1.Note, error) {
	if sts.store == nil {
		return nil, fmt.Errorf("store was nil: %w", errors.ErrStoreUnavailable)
	}
	if strings.TrimSpace(itemID) == "" {
		return nil, fmt.Errorf("itemID: %w", errors.ErrEmptyString)
	}
	return sts.store.FindNotes(ctx, itemID)
}

// GetNote gets a [statusthingv1.Note] by id
func (sts *StatusThingService) GetNote(ctx context.Context, noteID string) (*statusthingv1.Note, error) {
	if sts.store == nil {
		return nil, fmt.Errorf("store was nil: %w", errors.ErrStoreUnavailable)
	}
	if strings.TrimSpace(noteID) == "" {
		return nil, fmt.Errorf("noteID: %w", errors.ErrEmptyString)
	}
	return sts.store.GetNote(ctx, noteID)
}
