package sqlite

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers/internal"
	_ "github.com/lusis/statusthing/internal/storers/sqlite/driver" // sql driver
	"github.com/lusis/statusthing/internal/validation"
)

type dbStatus = internal.DbStatus

// StoreStatus stores the provided [statusthingv1.Status]
func (s *Store) StoreStatus(ctx context.Context, status *statusthingv1.Status) (*statusthingv1.Status, error) {
	rec, recerr := internal.DbStatusFromProto(status)
	if recerr != nil {
		return nil, recerr
	}
	if err := s.storeStruct(ctx, statusTableName, rec); err != nil {
		return nil, err
	}

	return s.GetStatus(ctx, rec.ID)
}

// GetStatus gets a [statusthingv1.Status] by its unique id
func (s *Store) GetStatus(ctx context.Context, statusID string) (*statusthingv1.Status, error) {
	rec := &dbStatus{}
	res, err := s.getProto(ctx, "id", statusID, statusTableName, rec)
	if err != nil {
		return nil, err
	}
	r, ok := res.(*statusthingv1.Status)
	if !ok {
		return nil, serrors.NewError("casting-status", serrors.ErrInvalidData)
	}
	return r, nil
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
		pb, pberr := rec.ToProto()
		if pberr != nil {
			return nil, serrors.NewWrappedError("proto", serrors.ErrUnrecoverable, pberr)
		}
		pbResults = append(pbResults, pb.(*statusthingv1.Status))
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
	return s.update(ctx, statusTableName, idColumn, statusID, columns)
}

// DeleteStatus deletes a [statusthingv1.Status] by its id
func (s *Store) DeleteStatus(ctx context.Context, statusID string) error {
	if !validation.ValidString(statusID) {
		return serrors.NewError("statusid", serrors.ErrEmptyString)
	}

	if _, existserr := s.GetStatus(ctx, statusID); existserr != nil {
		return existserr
	}
	return s.del(ctx, statusTableName, "id", statusID)
}
