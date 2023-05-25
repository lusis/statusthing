// Package internal contains internal storer code
package internal

// DbCommon is a struct that maps to common fields we use
type DbCommon struct {
	ID          string  `db:"id" goqu:"skipupdate"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
	Created     uint64  `db:"created"`
	Updated     uint64  `db:"updated"`
	Deleted     *uint64 `db:"deleted"`
}
