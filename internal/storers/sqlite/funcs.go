package sqlite

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"

	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers/internal"
	_ "github.com/lusis/statusthing/internal/storers/sqlite/driver" // sql driver

	"google.golang.org/protobuf/proto"

	"modernc.org/sqlite"
)

func (s *Store) storeStruct(ctx context.Context, tableName string, st any) error {
	ds := s.goqudb.Insert(tableName).Prepared(true).Rows(st)
	query, params, qerr := ds.ToSQL()
	if qerr != nil {
		return serrors.NewWrappedError("querybuilder", serrors.ErrUnrecoverable, qerr)
	}
	res, reserr := s.db.ExecContext(ctx, query, params...)
	if reserr != nil {
		if e, ok := reserr.(*sqlite.Error); ok {
			return serrors.NewWrappedError("driver", serrors.ErrStoreUnavailable, e)
		}
		return serrors.NewWrappedError("write", serrors.ErrUnrecoverable, reserr)
	}
	if _, lerr := res.LastInsertId(); lerr != nil {
		return serrors.NewWrappedError("last-insert-id", serrors.ErrUnrecoverable, lerr)
	}
	return nil
}

func (s *Store) getProto(ctx context.Context, idCol string, idVal string, tableName string, dbType internal.DbProtoable) (proto.Message, error) {
	ds := s.goqudb.From(tableName).Prepared(true)
	found, ferr := ds.Where(goqu.C(idCol).Eq(idVal)).ScanStructContext(ctx, dbType)
	if ferr != nil {
		return nil, serrors.NewWrappedError("read", serrors.ErrStoreUnavailable, ferr)
	}
	if found {
		return dbType.ToProto()
	}
	return nil, serrors.NewError("record", serrors.ErrNotFound)
}

func (s *Store) del(ctx context.Context, tableName string, idCol string, idVal string) error {
	res, reserr := s.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE %s = ?", tableName, idCol), idVal)
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

func (s *Store) update(ctx context.Context, tableName string, idCol string, idVal string, cols map[string]any) error {
	query, params, qerr := s.goqudb.Update(tableName).Prepared(true).Where(goqu.C(idCol).Eq(idVal)).Set(cols).ToSQL()
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
