// nolint :revive
// TODO: remove ^ this
package sqlite

import (
	"context"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers/internal"
	"github.com/lusis/statusthing/internal/validation"
	"google.golang.org/protobuf/proto"
	"modernc.org/sqlite"
)

type dbItem = internal.DbItem

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
	rec, recerr := internal.DbItemFromProto(item)
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
		pbrec, pberr := rec.ToProto()
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

		notes, nerr := s.FindNotes(ctx, rec.ID)
		if nerr != nil {
			return nil, serrors.NewWrappedError("notes", serrors.ErrInvalidData, nerr)
		}
		pbrec.Notes = notes
		return pbrec, nil
	}
	return nil, serrors.NewError("items", serrors.ErrNotFound)
}

// FindItems returns all known [statusthingv1.Item] optionally filtered by the provided [filters.FilterOption]
func (s *Store) FindItems(ctx context.Context, opts ...filters.FilterOption) ([]*statusthingv1.Item, error) {
	f, ferr := filters.New(opts...)
	if ferr != nil {
		return nil, ferr
	}
	dbitems := []*dbItem{}
	pbitems := []*statusthingv1.Item{}

	// this one gets a little weirder due to the join
	// note the use of goqu.I instead of just the column name
	// this indicates we want the identifier table.col not a column called "table.col"
	ds := s.goqudb.From(itemsTableName).Prepared(true).
		Select("items.*").
		LeftJoin(goqu.T(statusTableName), goqu.On(goqu.I("items.status_id").Eq(goqu.I("status.id"))))
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
		exprs = append(exprs, goqu.I("status.kind").In(kindStrings)) // TOOD: cleanup column names
	}

	if len(f.StatusIDs()) != 0 {
		exprs = append(exprs, goqu.I("status.id").In(f.StatusIDs())) // TOOD: cleanup column names
	}

	where := ds.Where(goqu.Or(exprs...)).Order(goqu.C(idColumn).Asc())
	werr := where.ScanStructsContext(ctx, &dbitems)
	if werr != nil {
		return nil, serrors.NewWrappedError("driver", serrors.ErrUnrecoverable, werr)
	}
	for _, i := range dbitems {
		// we're going to reuse our GetItem here so status gets populated properly
		pbitem, err := s.GetItem(ctx, i.ID)
		if err != nil {
			return nil, err
		}
		pbitems = append(pbitems, pbitem)
	}
	return pbitems, nil
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

	query, params, qerr := s.goqudb.Update(itemsTableName).Prepared(true).Where(goqu.I(idColumn).Eq(itemID)).Set(columns).ToSQL()
	if qerr != nil {
		return serrors.NewWrappedError("sqlbuilder", serrors.ErrUnrecoverable, qerr)
	}

	res, reserr := s.db.ExecContext(ctx, query, params...)
	if reserr != nil {
		var sqliteErr *sqlite.Error
		if errors.As(reserr, &sqliteErr) {
			if sqliteErr.Code() == 787 {
				return serrors.NewWrappedError("item", serrors.ErrNotFound, sqliteErr)
			}
		}
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
