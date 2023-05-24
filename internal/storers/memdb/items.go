package memdb

import (
	"context"
	"fmt"
	"strings"

	hcmemdb "github.com/hashicorp/go-memdb"
	"google.golang.org/protobuf/types/known/timestamppb"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
)

const (
	itemTableName    = "items"
	itemPrimaryIndex = "id"
)

type dbItem struct {
	ID          string
	Name        string
	Description string
	StatusID    string
	Created     int
	Updated     int
	Deleted     int
}

func memItemFromProto(pn *statusthingv1.Item) (*dbItem, error) {
	if pn == nil {
		return nil, fmt.Errorf("item: %w", serrors.ErrNilVal)
	}
	if pn.GetId() == "" {
		return nil, fmt.Errorf("id: %w", serrors.ErrEmptyString)
	}
	if pn.GetName() == "" {
		return nil, fmt.Errorf("name: %w", serrors.ErrEmptyString)
	}
	if pn.GetTimestamps() == nil {
		return nil, fmt.Errorf("timestamps: %w", serrors.ErrNilVal)
	}
	if !pn.GetTimestamps().GetCreated().IsValid() {
		return nil, fmt.Errorf("created: %w", serrors.ErrInvalidData)
	}
	if !pn.GetTimestamps().GetUpdated().IsValid() {
		return nil, fmt.Errorf("updated: %w", serrors.ErrInvalidData)
	}
	n := &dbItem{
		ID:          pn.GetId(),
		Name:        pn.GetName(),
		Description: pn.GetDescription(),
		Created:     tsToInt(pn.GetTimestamps().GetCreated()),
		Updated:     tsToInt(pn.GetTimestamps().GetUpdated()),
	}
	if pn.GetTimestamps().GetDeleted().IsValid() {
		n.Deleted = tsToInt(pn.GetTimestamps().GetDeleted())
	}
	if pn.GetStatus() != nil {
		n.StatusID = pn.GetStatus().GetId()
	}
	if pn.GetDescription() == "" {
		n.Description = pn.GetDescription()
	}

	return n, nil
}

func (mi *dbItem) toProto() (*statusthingv1.Item, error) {
	if mi.ID == "" {
		return nil, fmt.Errorf("id: %w", serrors.ErrInvalidData)
	}
	if mi.Name == "" {
		return nil, fmt.Errorf("name: %w", serrors.ErrInvalidData)
	}
	item := &statusthingv1.Item{
		Id:         mi.ID,
		Name:       mi.Name,
		Timestamps: &v1.Timestamps{},
	}
	if mi.Description != "" {
		item.Description = mi.Description
	}
	// timestamps
	created, updated, deleted := intToTs(mi.Created), intToTs(mi.Updated), intToTs(mi.Deleted)
	if !created.IsValid() {
		return nil, fmt.Errorf("created: %w", serrors.ErrInvalidData)
	}
	item.Timestamps.Created = created
	if !updated.IsValid() {
		return nil, fmt.Errorf("updated: %w", serrors.ErrInvalidData)
	}
	item.Timestamps.Updated = updated
	if mi.Description != "" {
		item.Description = mi.Description
	}
	if deleted.IsValid() {
		item.Timestamps.Deleted = deleted
	}
	// we don't process status here. that's done at call site
	return item, nil
}

var itemSchema = &hcmemdb.TableSchema{
	Name: itemTableName,
	Indexes: map[string]*hcmemdb.IndexSchema{
		"id": {
			Name:    itemPrimaryIndex,
			Unique:  true,
			Indexer: &hcmemdb.StringFieldIndex{Field: "ID"},
		},
		"name": {
			Name:    "name",
			Unique:  true,
			Indexer: &hcmemdb.StringFieldIndex{Field: "Name"},
		},
		"description": {
			Name:         "description",
			AllowMissing: true,
			Indexer:      &hcmemdb.StringFieldIndex{Field: "Description"},
		},
		"status_id": {
			Name:         "status_id",
			AllowMissing: true,
			Indexer:      &hcmemdb.StringFieldIndex{Field: "StatusID"},
		},
		"created": {
			Name:    "created",
			Indexer: &hcmemdb.IntFieldIndex{Field: "Created"},
		},
		"updated": {
			Name:    "updated",
			Indexer: &hcmemdb.IntFieldIndex{Field: "Updated"},
		},
		"deleted": {
			Name:         "deleted",
			AllowMissing: true,
			Indexer:      &hcmemdb.IntFieldIndex{Field: "Deleted"},
		},
	},
}

// StoreItem stores the provided [statusthingv1.Item]
func (s *StatusThingStore) StoreItem(ctx context.Context, item *v1.Item) (*v1.Item, error) {
	if item.GetStatus() != nil {
		// we have to create the status we're given if provided
		sres, serr := s.StoreStatus(ctx, item.GetStatus())
		if serr != nil {
			return nil, serr
		}
		item.Status = sres
	}
	txn := s.db.Txn(true)
	mItem, err := memItemFromProto(item)
	if err != nil {
		txn.Abort()
		return nil, err
	}
	if err := txn.Insert(itemTableName, mItem); err != nil {
		txn.Abort()
		return nil, err
	}
	txn.Commit()
	return s.GetItem(ctx, mItem.ID)
}

// GetItem gets a [statusthingv1.Item] by its id
func (s *StatusThingStore) GetItem(ctx context.Context, itemID string) (*v1.Item, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	res, err := firstWithTxn(ctx, txn, itemTableName, itemPrimaryIndex, itemID)
	if err != nil {
		return nil, err
	}
	item, ok := res.(*dbItem)
	if !ok {
		return nil, serrors.ErrInvalidData
	}
	pitem, err := item.toProto()
	if err != nil {
		return nil, err
	}

	if item.StatusID != "" {
		status, err := getStatusWithTxn(ctx, txn, item.StatusID)
		if err != nil {
			if err == serrors.ErrNotFound {
				return nil, fmt.Errorf("status not found: %w", serrors.ErrInvalidData)
			}
			return nil, err
		}
		pitem.Status = status
	}

	notes, err := getNotesWithTxn(ctx, txn, item.ID)
	if err != nil {
		return nil, err
	}
	pitem.Notes = notes

	return pitem, nil
}

// FindItems returns all known [statusthingv1.Item] optionally filtered by the provided [filters.FilterOption]
func (s *StatusThingStore) FindItems(ctx context.Context, opts ...filters.FilterOption) ([]*v1.Item, error) {
	// we can't do the kind of joins we would normally do with another db here so
	// we'll have to filter after getting results (which is fine)
	f, err := filters.New(opts...)
	if err != nil {
		return nil, err
	}

	txn := s.db.Txn(false)
	defer txn.Abort()

	res, err := getWithTxn(ctx, txn, itemTableName, itemPrimaryIndex)
	if err != nil {
		return nil, err
	}
	items := []*statusthingv1.Item{}
	for obj := res.Next(); obj != nil; obj = res.Next() {
		entry, ok := obj.(*dbItem)
		if !ok {
			return nil, serrors.ErrInvalidData
		}
		// short-circuit any results that can never match when filters are provided
		if (len(f.StatusIDs()) != 0 || len(f.StatusKinds()) != 0) && entry.StatusID == "" {
			continue
		}
		if len(f.StatusIDs()) != 0 {
			var hasSid bool
			for _, sid := range f.StatusIDs() {
				if sid == entry.StatusID {
					hasSid = true
				}
			}
			if !hasSid {
				continue
			}
		}
		pb, err := entry.toProto()
		if err != nil {
			return nil, err
		}
		if entry.StatusID != "" {
			status, err := getStatusWithTxn(ctx, txn, entry.StatusID)
			if err != nil {
				return nil, fmt.Errorf("getting status: %w", err)
			}
			// we have to filter out if we were given status kinds
			if len(f.StatusKinds()) != 0 {
				var hasKind bool
				for _, kind := range f.StatusKinds() {
					if kind == status.GetKind() {
						hasKind = true
					}
				}
				if !hasKind {
					continue
				}
			}
			pb.Status = status
		}
		items = append(items, pb)
	}

	return items, nil
}

// UpdateItem updates the [statusthingv1.Item] by its id with the provided [filters.FilterOption]
// supported filters are
// - [filters.WithName]
// - [filters.WithDescription]
// - [filters.WithStatus]
func (s *StatusThingStore) UpdateItem(ctx context.Context, itemID string, opts ...filters.FilterOption) error {
	if len(opts) == 0 {
		return serrors.ErrAtLeastOne
	}
	if strings.TrimSpace(itemID) == "" {
		return serrors.ErrEmptyString
	}
	f, err := filters.New(opts...)
	if err != nil {
		return err
	}
	existing, err := s.GetItem(ctx, itemID)
	if err != nil {
		return err
	}

	if f.Name() != "" {
		existing.Name = f.Name()
	}
	if f.Description() != "" {
		existing.Description = f.Description()
	}
	// statusid/status are mutually exclusive
	// filter logic should disallow this but let's play it safe
	if f.StatusID() != "" {
		sres, err := s.GetStatus(ctx, f.StatusID())
		if err != nil {
			return err
		}
		existing.Status = sres
	} else if f.Status() != nil {
		existing.Status = f.Status()
	}

	existing.Timestamps.Updated = timestamppb.Now()
	// updates in memdb are just inserts to the same id so we can just call our own store func here
	if _, err := s.StoreItem(ctx, existing); err != nil {
		return err
	}

	return nil
}

// DeleteItem deletes the [statusthingv1.Item] by its id
func (s *StatusThingStore) DeleteItem(ctx context.Context, itemID string) error {
	txn := s.db.Txn(true)
	if err := deleteWithTxn(ctx, txn, itemTableName, &dbItem{ID: itemID}); err != nil {
		txn.Abort()
		return err
	}
	txn.Commit()
	return nil
}
