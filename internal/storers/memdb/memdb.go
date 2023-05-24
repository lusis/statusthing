package memdb

import (
	"context"
	"time"

	hcmemdb "github.com/hashicorp/go-memdb"

	"github.com/lusis/statusthing/internal/serrors"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var schema = &hcmemdb.DBSchema{
	Tables: map[string]*hcmemdb.TableSchema{
		"statuses": statusSchema,
		"items":    itemSchema,
		"notes":    notesSchema,
	},
}

// StatusThingStore ...
type StatusThingStore struct {
	db *hcmemdb.MemDB
}

// New ...
func New() (*StatusThingStore, error) {
	db, err := hcmemdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}

	return &StatusThingStore{db: db}, nil
}

func deleteWithTxn(_ context.Context, txn *hcmemdb.Txn, tableName string, item any) error {
	err := txn.Delete(tableName, item)
	if err != nil {
		if err == hcmemdb.ErrNotFound {
			return serrors.ErrNotFound
		}
		return err
	}
	return nil
}

func getWithTxn(_ context.Context, txn *hcmemdb.Txn, tableName string, index string, args ...interface{}) (hcmemdb.ResultIterator, error) {
	res, err := txn.Get(tableName, index, args...)
	if err != nil {
		if err == hcmemdb.ErrNotFound {
			return nil, serrors.ErrNotFound
		}
		return nil, err
	}
	if res == nil {
		return nil, serrors.ErrNotFound
	}
	return res, nil
}

func firstWithTxn(_ context.Context, txn *hcmemdb.Txn, tableName string, index string, args ...interface{}) (interface{}, error) { // nolint: unparam
	// unparam due to index string always being the same
	res, err := txn.First(tableName, index, args...)
	if err != nil {
		if err == hcmemdb.ErrNotFound {
			return nil, serrors.ErrNotFound
		}
		return nil, err
	}
	if res == nil {
		return nil, serrors.ErrNotFound
	}
	return res, nil
}

func insertWithTxn(_ context.Context, txn *hcmemdb.Txn, tableName string, item any) error {
	return txn.Insert(tableName, item)
}

// tsToInt ...
// Deprecated: move to [storers.tsToInt]
func tsToInt(ts *timestamppb.Timestamp) int {
	if ts == nil {
		return 0
	}
	return int(ts.AsTime().UnixNano())
}

// intToTs ...
// Deprecated: move to [storers.intToTs]
func intToTs(i int) *timestamppb.Timestamp {
	// zero int val is a nil for us
	if i == 0 {
		return nil
	}
	return timestamppb.New(time.Unix(0, int64(i)).UTC())
}
