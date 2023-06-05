package services

import (
	"context"
	"fmt"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/validation"

	"github.com/segmentio/ksuid"
)

// AddItem creates a new item with the provided name
// Supported options:
// - [filters.WithStatusID] set the initial status to the [statusthingv1.Status] with provided status id
// - [filters.WithDescription] sets the optional description of the [statusthingv1.Item]
// - [filters.WithItemID] sets a custom unique id. default is generated via [ksuid.New().String()]
// - [filters.WithNoteText] sets the note text for an initial note to create along with the item
// - [filters.WithStatus] creates a new [statusthingv1.Status] before creating the item and sets the items status to that new status
func (sts *StatusThingService) AddItem(ctx context.Context, name string, opts ...filters.FilterOption) (*statusthingv1.Item, error) {
	if sts.store == nil {
		return nil, serrors.NewError("store", serrors.ErrStoreUnavailable)
	}
	f, err := filters.New(opts...)
	if err != nil {
		return nil, err
	}
	id := ksuid.New().String()
	if f.ItemID() != "" {
		id = f.ItemID()
	}

	thing := &statusthingv1.Item{
		Id:         id,
		Name:       name,
		Timestamps: makeTsNow(),
	}
	statusID := f.StatusID()
	status := f.Status()
	desc := f.Description()
	noteText := f.NoteText()

	// this next bit gets unfortunately a bit convoluted due to flexibility
	// first we check if they have a status id provided
	// if so, we get that status and add it to the new item
	// otherwise we check if there's a status attached
	// and use that
	if validation.ValidString(statusID) {
		status, err := sts.store.GetStatus(ctx, statusID)
		if err != nil {
			return nil, serrors.NewWrappedError("provided-statusid", serrors.ErrNotFound, err)
		}
		thing.Status = status
	} else if status != nil {
		thing.Status = status
		// we have to create an id if one isn't present
		if f.Status().GetId() == "" {
			thing.Status.Id = ksuid.New().String()
		}
		// same for timestamps
		if f.Status().GetTimestamps() == nil {
			thing.Status.Timestamps = makeTsNow()
		}
	}
	if desc != "" {
		thing.Description = desc
	}
	res, err := sts.store.StoreItem(ctx, thing)
	if err != nil {
		return nil, err
	}

	if validation.ValidString(noteText) {
		_, nerr := sts.AddNote(ctx, res.GetId(), f.NoteText())
		if nerr != nil {
			return nil, nerr
		}
		// do a fresh get to include the notes
		return sts.GetItem(ctx, res.GetId())
	}
	return res, nil
}

// EditItem updates the [statusthingv1.Item] with the provided id
func (sts *StatusThingService) EditItem(ctx context.Context, itemID string, opts ...filters.FilterOption) error {
	if sts.store == nil {
		return serrors.NewError("store", serrors.ErrStoreUnavailable)
	}
	return sts.store.UpdateItem(ctx, itemID, opts...)
}

// RemoveItem removes a [statusthingv1.Item] by its unique id
func (sts *StatusThingService) RemoveItem(ctx context.Context, itemID string) error {
	if sts.store == nil {
		return serrors.NewError("store", serrors.ErrStoreUnavailable)
	}
	if !validation.ValidString(itemID) {
		return serrors.NewError("itemID", serrors.ErrEmptyString)
	}
	return sts.store.DeleteItem(ctx, itemID)
}

// FindItems returns all known [statusthingv1.Item]
// supported filters:
// - [filters.WithStatusIDs]: only return results having the provided status ids
// - [filters.WithStatusKinds]: only return restuls having the provided status kinds
// StatusIDs and StatusKinds are mutually exclusive
func (sts *StatusThingService) FindItems(ctx context.Context, opts ...filters.FilterOption) ([]*statusthingv1.Item, error) {
	if sts.store == nil {
		return nil, fmt.Errorf("store was nil: %w", serrors.ErrStoreUnavailable)
	}
	return sts.store.FindItems(ctx, opts...)
}

// GetItem gets a [statusthingv1.Item] by id
func (sts *StatusThingService) GetItem(ctx context.Context, itemID string) (*statusthingv1.Item, error) {
	if sts.store == nil {
		return nil, fmt.Errorf("store was nil: %w", serrors.ErrStoreUnavailable)
	}
	return sts.store.GetItem(ctx, itemID)
}
