package sqlite

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"

	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"

	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers/internal"
	"github.com/lusis/statusthing/internal/validation"

	"modernc.org/sqlite"
)

// StoreNote stores the provided [statusthingv1.Note] associated with the provided [statusthingv1.StatusThing] by its id
func (s *Store) StoreNote(ctx context.Context, note *v1.Note, itemID string) (*v1.Note, error) {
	rec, recerr := internal.DbNoteFromProto(note)
	if recerr != nil {
		return nil, recerr
	}
	if !validation.ValidString(itemID) {
		return nil, serrors.NewError("itemid", serrors.ErrEmptyString)
	}
	rec.ItemID = itemID
	ds := s.goqudb.Insert(notesTableName).Prepared(true).Rows(rec)
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
	return s.GetNote(ctx, rec.ID)
}

// GetNote gets a [statusthingv1.Note] by its id
func (s *Store) GetNote(ctx context.Context, noteID string) (*v1.Note, error) {
	rec := &internal.DbNote{}
	ds := s.goqudb.From(notesTableName).Prepared(true)
	found, ferr := ds.Where(goqu.C("id").Eq(noteID)).Order(goqu.C("id").Asc()).ScanStructContext(ctx, rec)
	if ferr != nil {
		return nil, serrors.NewWrappedError("read", serrors.ErrStoreUnavailable, ferr)
	}
	if found {
		return rec.ToProto()
	}
	return nil, serrors.NewError("status", serrors.ErrNotFound)
}

// FindNotes gets all known [statusthingv1.Note] for the provided item id
// no filters are supported at this time
func (s *Store) FindNotes(ctx context.Context, itemID string, _ ...filters.FilterOption) ([]*v1.Note, error) {
	dbresults := []*internal.DbNote{}
	pbresults := []*v1.Note{}
	if !validation.ValidString(itemID) {
		return nil, serrors.NewError("itemid", serrors.ErrEmptyString)
	}
	dserr := s.goqudb.From(notesTableName).Prepared(true).Where(goqu.C(itemIDColumn).Eq(itemID)).Order(goqu.C(idColumn).Asc()).ScanStructsContext(ctx, &dbresults)
	if dserr != nil {
		return nil, serrors.NewWrappedError("driver", serrors.ErrUnrecoverable, dserr)
	}
	for _, rec := range dbresults {
		pb, pberr := rec.ToProto()
		if pberr != nil {
			return nil, serrors.NewWrappedError("proto", serrors.ErrUnrecoverable, pberr)
		}
		pbresults = append(pbresults, pb)
	}
	return pbresults, nil
}

// UpdateNote updates the [statusthingv1.Note] with the provided [filters.FilterOption]
// supported filters:
// [filters.WithNoteText]
func (s *Store) UpdateNote(ctx context.Context, noteID string, opts ...filters.FilterOption) error {
	if len(opts) == 0 {
		return serrors.NewError("opts", serrors.ErrAtLeastOne)
	}
	f, ferr := filters.New(opts...)
	if ferr != nil {
		return ferr
	}
	_, eerr := s.GetNote(ctx, noteID)
	if eerr != nil {
		return eerr
	}
	if !validation.ValidString(f.NoteText()) {
		return serrors.NewError("text", serrors.ErrEmptyString)
	}

	query, params, qerr := s.goqudb.Update(notesTableName).Prepared(true).Where(goqu.C(idColumn).Eq(noteID)).
		Set(goqu.Record{"note_text": f.NoteText()}).ToSQL()
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

// DeleteNote deletes a [statusthingv1.Note] by its id
func (s *Store) DeleteNote(ctx context.Context, noteID string) error {
	if !validation.ValidString(noteID) {
		return serrors.NewError("noteid", serrors.ErrEmptyString)
	}

	if _, existserr := s.GetNote(ctx, noteID); existserr != nil {
		return existserr
	}
	res, reserr := s.db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = ?", notesTableName), noteID)
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
