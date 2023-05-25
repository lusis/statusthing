package sqlite

import (
	"context"
	"fmt"
	"html"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers"
	"github.com/lusis/statusthing/internal/validation"

	"modernc.org/sqlite"
)

type dbStatus struct {
	ID          string  `db:"id" goqu:"skipupdate"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
	Color       *string `db:"color"`
	Kind        *string `db:"kind"`
	Created     uint64  `db:"created"`
	Updated     uint64  `db:"updated"`
	Deleted     *uint64 `db:"deleted"`
}

// StoreStatus stores the provided [statusthingv1.Status]
func (s *Store) StoreStatus(ctx context.Context, status *statusthingv1.Status) (*statusthingv1.Status, error) {
	rec, recerr := dbStatusFromProto(status)
	if recerr != nil {
		return nil, recerr
	}

	ds := s.goqudb.Insert(statusTableName).Prepared(true).Rows(rec)
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
	return s.GetStatus(ctx, rec.ID)
}

// GetStatus gets a [statusthingv1.Status] by its unique id
func (s *Store) GetStatus(ctx context.Context, statusID string) (*statusthingv1.Status, error) {
	rec := &dbStatus{}
	ds := s.goqudb.From(statusTableName).Prepared(true)
	found, ferr := ds.Where(goqu.C("id").Eq(statusID)).Order(goqu.C("id").Asc()).ScanStructContext(ctx, rec)
	if ferr != nil {
		return nil, serrors.NewWrappedError("read", serrors.ErrStoreUnavailable, ferr)
	}
	if found {
		return rec.toProto()
	}
	return nil, serrors.NewError("status", serrors.ErrNotFound)
}

// FindStatus returns all know [statusthingv1.Status] optionally filtered by the provided [filters.FilterOption]
func (s *Store) FindStatus(ctx context.Context, opts ...filters.FilterOption) ([]*statusthingv1.Status, error) {
	f, ferr := filters.New(opts...)
	if ferr != nil {
		return nil, ferr
	}

	dbResults := []*dbStatus{}
	pbResults := []*statusthingv1.Status{}

	ds := s.goqudb.From(statusTableName).Prepared(true)
	exprs := []exp.Expression{}
	// add status kinds
	if len(f.StatusKinds()) != 0 {
		kindStrings := func() []string {
			ss := []string{}
			for _, k := range f.StatusKinds() {
				ss = append(ss, k.String())
			}
			return ss
		}()
		exprs = append(exprs, goqu.C(kindColumn).In(kindStrings))
	}

	if len(f.StatusIDs()) != 0 {
		exprs = append(exprs, goqu.C(idColumn).In(f.StatusIDs()))
	}
	where := ds.Where(goqu.Or(exprs...)).Order(goqu.C(idColumn).Asc())
	werr := where.ScanStructsContext(ctx, &dbResults)
	if werr != nil {
		return nil, serrors.NewWrappedError("driver", serrors.ErrUnrecoverable, werr)
	}
	for _, rec := range dbResults {
		pb, pberr := rec.toProto()
		if pberr != nil {
			return nil, serrors.NewWrappedError("proto", serrors.ErrUnrecoverable, pberr)
		}
		pbResults = append(pbResults, pb)
	}
	return pbResults, nil
}

// UpdateStatus updates the [statusthingv1.Status] by id with the provided [filters.FilterOption]
func (s *Store) UpdateStatus(ctx context.Context, statusID string, opts ...filters.FilterOption) error {
	f, ferr := filters.New(opts...)
	if ferr != nil {
		return ferr
	}
	_, eerr := s.GetStatus(ctx, statusID)
	if eerr != nil {
		return eerr
	}

	kind := f.StatusKind()
	color := f.Color()
	name := f.Name()
	desc := f.Description()

	columns := map[string]any{}
	if kind != statusthingv1.StatusKind_STATUS_KIND_UNKNOWN {
		columns[kindColumn] = kind.String()
	}
	if validation.ValidString(color) {
		columns[colorColumn] = color
	}
	if validation.ValidString(name) {
		columns[nameColumn] = name
	}
	if validation.ValidString(desc) {
		columns[descriptionColumn] = desc
	}

	query, params, qerr := s.goqudb.Update(statusTableName).Prepared(true).Where(goqu.C(idColumn).Eq(statusID)).Set(columns).ToSQL()
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

// DeleteStatus deletes a [statusthingv1.Status] by its id
func (s *Store) DeleteStatus(ctx context.Context, statusID string) error {
	if !validation.ValidString(statusID) {
		return serrors.NewError("statusid", serrors.ErrEmptyString)
	}

	if _, existserr := s.GetStatus(ctx, statusID); existserr != nil {
		return existserr
	}
	res, reserr := s.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = ?", statusTableName), statusID)
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

func dbStatusFromProto(pbstatus *statusthingv1.Status) (*dbStatus, error) {
	if pbstatus == nil {
		return nil, serrors.NewError("status", serrors.ErrNilVal)
	}
	id := html.EscapeString(pbstatus.GetId())
	name := html.EscapeString(pbstatus.GetName())
	desc := html.EscapeString(pbstatus.GetDescription())
	color := html.EscapeString(pbstatus.GetColor())
	created := pbstatus.GetTimestamps().GetCreated()
	updated := pbstatus.GetTimestamps().GetUpdated()
	deleted := pbstatus.GetTimestamps().GetDeleted()
	kind := pbstatus.GetKind()

	if !validation.ValidString(id) {
		return nil, serrors.NewError("id", serrors.ErrEmptyString)
	}
	if !validation.ValidString(name) {
		return nil, serrors.NewError("name", serrors.ErrEmptyString)
	}
	if kind == statusthingv1.StatusKind_STATUS_KIND_UNKNOWN {
		return nil, serrors.NewError("kind", serrors.ErrEmptyEnum)
	}

	if pbstatus.GetTimestamps() == nil {
		return nil, serrors.NewError("timestamps", serrors.ErrMissingTimestamp)
	}

	if !created.IsValid() {
		return nil, serrors.NewError("created", serrors.ErrMissingTimestamp)
	}
	if !updated.IsValid() {
		return nil, serrors.NewError("updated", serrors.ErrMissingTimestamp)
	}
	dbs := &dbStatus{
		ID:      id,
		Name:    name,
		Kind:    storers.StringPtr(kind.String()),
		Created: storers.TsToUInt64(created),
		Updated: storers.TsToUInt64(updated),
	}
	if validation.ValidString(desc) {
		dbs.Description = storers.StringPtr(desc)
	}
	if validation.ValidString(color) {
		dbs.Color = storers.StringPtr(color)
	}
	if deleted.IsValid() {
		dbs.Deleted = storers.TsToUInt64Ptr(deleted)
	}
	return dbs, nil
}
func (s *dbStatus) toProto() (*statusthingv1.Status, error) {
	res := &statusthingv1.Status{
		Timestamps: &statusthingv1.Timestamps{},
	}
	if !validation.ValidString(s.ID) {
		return nil, serrors.NewError("id", serrors.ErrInvalidData)
	}
	if !validation.ValidString(s.Name) {
		return nil, serrors.NewError("name", serrors.ErrInvalidData)
	}
	if s.Kind == nil || s.Kind == storers.StringPtr(statusthingv1.StatusKind_STATUS_KIND_UNKNOWN.String()) {
		return nil, serrors.NewError("kind", serrors.ErrInvalidData)
	}

	// id/name
	res.Id = html.UnescapeString(s.ID)
	res.Name = html.UnescapeString(s.Name)

	// kind
	res.Kind = statusthingv1.StatusKind(statusthingv1.StatusKind_value[*s.Kind])

	//desc/color
	if s.Description != nil {
		res.Description = html.UnescapeString(*s.Description)
	}
	if s.Color != nil {
		res.Color = html.UnescapeString(*s.Color)
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
	return res, nil
}
