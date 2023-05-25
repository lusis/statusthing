// nolint :revive
// TODO: remove ^ this
package sqlite

import (
	"context"
	"fmt"
	"html"

	"github.com/doug-martin/goqu/v9"
	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers"
	"github.com/lusis/statusthing/internal/validation"
	"google.golang.org/protobuf/proto"
	"modernc.org/sqlite"
)

type dbItem struct {
	ID          string  `db:"id" goqu:"skipupdate"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
	StatusID    *string `db:"status_id"`
	Created     uint64  `db:"created"`
	Updated     uint64  `db:"updated"`
	Deleted     *uint64 `db:"deleted"`
}

// StoreItem stores the provided [statusthingv1.Item]
func (s *Store) StoreItem(ctx context.Context, item *statusthingv1.Item) (*statusthingv1.Item, error) {

	if item.GetStatus() != nil {
		status := proto.Clone(item.GetStatus()).(*statusthingv1.Status)
		// if the status has an id, we check if it exists
		if validation.ValidString(status.GetId()) {
			existing, _ := s.GetStatus(ctx, item.GetStatus().GetId())
			if existing == nil {
				// it's not there so we can attempt to create it
				isres, isreserr := s.StoreStatus(ctx, status)
				if isreserr != nil {
					return nil, isreserr
				}
				if isres != nil {
					status = isres
				}
			}

		}
		item.Status = status
	}
	rec, recerr := dbItemFromProto(item)
	if recerr != nil {
		return nil, recerr
	}

	ds := s.goqudb.Insert(itemsTableName).Prepared(true).Rows(rec)
	query, params, qerr := ds.ToSQL()
	if qerr != nil {
		return nil, serrors.NewWrappedError("querybuilder", serrors.ErrUnrecoverable, qerr)
	}
	res, reserr := s.db.ExecContext(ctx, query, params...)
	if reserr != nil {
		if e, ok := reserr.(*sqlite.Error); ok {
			return nil, serrors.NewWrappedError("driver", serrors.ErrStoreUnavailable, e)
		}
		return nil, serrors.NewWrappedError("write", serrors.ErrUnrecoverable, reserr)
	}
	if _, lerr := res.LastInsertId(); lerr != nil {
		return nil, serrors.NewWrappedError("last-insert-id", serrors.ErrUnrecoverable, lerr)
	}

	return s.GetItem(ctx, rec.ID)
}

// GetItem gets a [statusthingv1.Item] by its id
func (s *Store) GetItem(ctx context.Context, itemID string) (*statusthingv1.Item, error) {
	rec := &dbItem{}
	ds := s.goqudb.From(itemsTableName).Prepared(true)
	found, ferr := ds.Where(goqu.C("id").Eq(itemID)).Order(goqu.C("id").Asc()).ScanStructContext(ctx, rec)
	if ferr != nil {
		return nil, serrors.NewWrappedError("read", serrors.ErrStoreUnavailable, ferr)
	}

	if found {
		pbrec, pberr := rec.toProto()
		if pberr != nil {
			return nil, serrors.NewWrappedError("proto", serrors.ErrUnrecoverable, pberr)
		}
		if rec.StatusID != nil {
			status, serr := s.GetStatus(ctx, *rec.StatusID)
			if serr != nil {
				return nil, serr
			}
			pbrec.Status = status
		}

		return pbrec, nil
	}
	return nil, serrors.NewError("items", serrors.ErrNotFound)
}

// FindItems returns all known [statusthingv1.Item] optionally filtered by the provided [filters.FilterOption]
func (s *Store) FindItems(ctx context.Context, opts ...filters.FilterOption) ([]*statusthingv1.Item, error) {
	/*
		select items.* from items JOIN status on items.status_id = status.id WHERE items.status_id IS NOT NULL AND (items.status_id IN () OR status.kind IN ());
	*/
	panic("not implemented") // TODO: Implement
}

// UpdateItem updates the [statusthingv1.Item] by its id with the provided [filters.FilterOption]
// Supported filters:
// - [filters.WithStatusID]
// - [filters.WithName]
// - [filters.WithDescription]
func (s *Store) UpdateItem(ctx context.Context, itemID string, opts ...filters.FilterOption) error {
	f, ferr := filters.New(opts...)
	if ferr != nil {
		return ferr
	}
	_, eerr := s.GetItem(ctx, itemID)
	if eerr != nil {
		return eerr
	}

	name := f.Name()
	desc := f.Description()
	statusID := f.StatusID()

	columns := map[string]any{}

	if validation.ValidString(statusID) {
		columns[statusIDColumn] = statusID
	}
	if validation.ValidString(name) {
		columns[nameColumn] = name
	}
	if validation.ValidString(desc) {
		columns[descriptionColumn] = desc
	}

	query, params, qerr := s.goqudb.Update(itemsTableName).Prepared(true).Where(goqu.C(idColumn).Eq(itemID)).Set(columns).ToSQL()
	if qerr != nil {
		return serrors.NewWrappedError("driver", serrors.ErrUnrecoverable, qerr)
	}

	res, reserr := s.db.ExecContext(ctx, query, params...)
	if reserr != nil {
		return serrors.NewWrappedError("write", serrors.ErrUnrecoverable, reserr)
	}
	if _, lerr := res.LastInsertId(); lerr != nil {
		return serrors.NewWrappedError("last-insert-id", serrors.ErrUnrecoverable, lerr)
	}
	return nil
}

// DeleteItem deletes the [statusthingv1.Item] by its id
func (s *Store) DeleteItem(ctx context.Context, itemID string) error {
	if !validation.ValidString(itemID) {
		return serrors.NewError("itemid", serrors.ErrEmptyString)
	}

	if _, existserr := s.GetItem(ctx, itemID); existserr != nil {
		return existserr
	}
	res, reserr := s.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = ?", itemsTableName), itemID)
	if reserr != nil {
		return serrors.NewWrappedError("write", serrors.ErrUnrecoverable, reserr)
	}
	affected, aferr := res.RowsAffected()
	if aferr != nil {
		return serrors.NewWrappedError("affected-rows", serrors.ErrUnrecoverable, aferr)
	}
	if affected != 1 {
		// we checked for existence earlier so this should only return if we delete more than one row
		// we don't need to account for zero rows here but we might want to do an optimistic delete instead and handle zero differently
		return serrors.NewError(fmt.Sprintf("%d rows affected", affected), serrors.ErrUnexpectedRows)
	}
	return nil
}

func dbItemFromProto(pbitem *statusthingv1.Item) (*dbItem, error) {
	if pbitem == nil {
		return nil, serrors.NewError("status", serrors.ErrNilVal)
	}
	id := html.EscapeString(pbitem.GetId())
	name := html.EscapeString(pbitem.GetName())
	desc := html.EscapeString(pbitem.GetDescription())
	statusID := pbitem.GetStatus().GetId()
	created := pbitem.GetTimestamps().GetCreated()
	updated := pbitem.GetTimestamps().GetUpdated()
	deleted := pbitem.GetTimestamps().GetDeleted()

	if !validation.ValidString(id) {
		return nil, serrors.NewError("id", serrors.ErrEmptyString)
	}
	if !validation.ValidString(name) {
		return nil, serrors.NewError("name", serrors.ErrEmptyString)
	}
	if pbitem.GetTimestamps() == nil {
		return nil, serrors.NewError("timestamps", serrors.ErrMissingTimestamp)
	}

	if !created.IsValid() {
		return nil, serrors.NewError("created", serrors.ErrMissingTimestamp)
	}
	if !updated.IsValid() {
		return nil, serrors.NewError("updated", serrors.ErrMissingTimestamp)
	}
	dbs := &dbItem{
		ID:      id,
		Name:    name,
		Created: storers.TsToUInt64(created),
		Updated: storers.TsToUInt64(updated),
	}

	if validation.ValidString(statusID) {
		dbs.StatusID = storers.StringPtr(statusID)
	}

	if validation.ValidString(desc) {
		dbs.Description = storers.StringPtr(desc)
	}

	if deleted.IsValid() {
		dbs.Deleted = storers.TsToUInt64Ptr(deleted)
	}
	return dbs, nil
}
func (s *dbItem) toProto() (*statusthingv1.Item, error) {
	res := &statusthingv1.Item{
		Timestamps: &statusthingv1.Timestamps{},
	}
	if !validation.ValidString(s.ID) {
		return nil, serrors.NewError("id", serrors.ErrInvalidData)
	}
	if !validation.ValidString(s.Name) {
		return nil, serrors.NewError("name", serrors.ErrInvalidData)
	}

	// id/name
	res.Id = html.UnescapeString(s.ID)
	res.Name = html.UnescapeString(s.Name)

	//desc
	if s.Description != nil {
		res.Description = html.UnescapeString(*s.Description)
	}

	// timestamps
	pbcreated := storers.Int64ToTs(int64(s.Created))
	pbupdated := storers.Int64ToTs(int64(s.Updated))

	if pbcreated == nil {
		return nil, serrors.NewError("created", serrors.ErrInvalidData)
	}
	if pbupdated == nil {
		return nil, serrors.NewError("updated", serrors.ErrInvalidData)
	}
	res.Timestamps.Created = pbcreated
	res.Timestamps.Updated = pbupdated

	if s.Deleted != nil {
		res.Timestamps.Deleted = storers.Int64ToTs(int64(*s.Deleted))
	}
	// status will be populated outside of here for now
	// but that creates a potential tangled dep down the road
	return res, nil
}
