package services

import (
	"context"
	"fmt"
	"strings"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/errors"
	"github.com/lusis/statusthing/internal/filters"

	"github.com/segmentio/ksuid"
)

// NewItem creates a new item with the provided name
// Supported options:
// - [filters.WithStatusID] set the initial status to the [statusthingv1.Status] with provided status id
// - [filters.WithDescription] sets the optional description of the [statusthingv1.Item]
// - [filters.WithItemID] sets a custom unique id. default is generated via [ksuid.New().String()]
// - [filters.WithNoteText] sets the note text for an initial note to create along with the item
// - [filters.WithStatus] creates a new [statusthingv1.Status] before creating the item and sets the items status to that new status
func (sts *StatusThingService) NewItem(ctx context.Context, name string, opts ...filters.FilterOption) (*statusthingv1.Item, error) {
	if sts.store == nil {
		return nil, fmt.Errorf("store was nil: %w", errors.ErrStoreUnavailable)
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

	// this next bit gets unfortunately a bit convoluted due to flexibility
	// first we check if they have a status id provided
	// if so, we get that status and add it to the new item
	// otherwise we check if there's a status attached
	// and use that
	if f.StatusID() != "" {
		status, err := sts.store.GetStatus(ctx, f.StatusID())
		if err != nil {
			return nil, fmt.Errorf("provided status id: %w", errors.ErrNotFound)
		}
		thing.Status = status
	} else if f.Status() != nil {
		thing.Status = f.Status()
		// we have to create an id if one isn't present
		if f.Status().GetId() == "" {
			thing.Status.Id = ksuid.New().String()
		}
		// same for timestamps
		if f.Status().GetTimestamps() == nil {
			thing.Status.Timestamps = makeTsNow()
		}
	}
	if f.Description() != "" {
		thing.Description = f.Description()
	}
	res, err := sts.store.StoreItem(ctx, thing)
	if err != nil {
		return nil, err
	}

	if f.NoteText() != "" {
		_, nerr := sts.NewNote(ctx, res.GetId(), f.NoteText())
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
		return fmt.Errorf("store was nil: %w", errors.ErrStoreUnavailable)
	}
	return sts.store.UpdateItem(ctx, itemID, opts...)
}

// DeleteItem removes a [statusthingv1.Item] by its unique id
func (sts *StatusThingService) DeleteItem(ctx context.Context, itemID string) error {
	if sts.store == nil {
		return fmt.Errorf("store was nil: %w", errors.ErrStoreUnavailable)
	}
	if strings.TrimSpace(itemID) == "" {
		return fmt.Errorf("itemID: %w", errors.ErrEmptyString)
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
		return nil, fmt.Errorf("store was nil: %w", errors.ErrStoreUnavailable)
	}
	return sts.store.FindItems(ctx, opts...)
}

// GetItem gets a [statusthingv1.Item] by id
func (sts *StatusThingService) GetItem(ctx context.Context, itemID string) (*statusthingv1.Item, error) {
	if sts.store == nil {
		return nil, fmt.Errorf("store was nil: %w", errors.ErrStoreUnavailable)
	}
	return sts.store.GetItem(ctx, itemID)
}
