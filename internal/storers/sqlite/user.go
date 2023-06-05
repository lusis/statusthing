// nolint :unused
package sqlite

import (
	"context"

	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers"
	"github.com/lusis/statusthing/internal/storers/internal"
	_ "github.com/lusis/statusthing/internal/storers/sqlite/driver" // sql driver
	"github.com/lusis/statusthing/internal/validation"
)

// StoreUser stores the provied [v1.User]
func (s *Store) StoreUser(ctx context.Context, user *v1.User) (*v1.User, error) {
	rec, recerr := internal.DbUserFromProto(user)
	if recerr != nil {
		return nil, recerr
	}
	if err := s.storeStruct(ctx, usersTableName, rec); err != nil {
		return nil, err
	}
	return s.GetUser(ctx, rec.Username)
}

// GetUser gets a [v1.User] by id
func (s *Store) GetUser(ctx context.Context, username string) (*v1.User, error) {
	rec := &internal.DbUser{}
	res, err := s.getProto(ctx, "username", username, usersTableName, rec)
	if err != nil {
		return nil, err
	}
	r, ok := res.(*v1.User)
	if !ok {
		return nil, serrors.NewError("casting-user", serrors.ErrInvalidData)
	}
	return r, nil
}

// FindUsers finds users
func (s *Store) FindUsers(ctx context.Context, opts ...filters.FilterOption) ([]*v1.User, error) {
	panic("not implemented") // TODO: Implement
}

// UpdateUser updates a [v1.User]
func (s *Store) UpdateUser(ctx context.Context, username string, opts ...filters.FilterOption) error {
	f, ferr := filters.New(opts...)
	if ferr != nil {
		return ferr
	}
	_, eerr := s.GetUser(ctx, username)
	if eerr != nil {
		return eerr
	}

	fname := f.FirstName()
	lname := f.LastName()
	email := f.EmailAddress()
	lastlogin := f.LastLogin()
	avatarURL := f.AvatarURL()
	password := f.Password()

	columns := map[string]any{}

	if validation.ValidString(fname) {
		columns[fnameColumn] = fname
	}
	if validation.ValidString(lname) {
		columns[lnameColumn] = lname
	}
	if validation.ValidString(email) {
		columns[emailColumn] = email
	}
	if validation.ValidString(avatarURL) {
		columns[avatarURLColumn] = avatarURL
	}
	if validation.ValidString(password) {
		columns[passwordColumn] = password
	}
	if lastlogin != nil {
		columns[lastloginColumn] = storers.TimeToUint64(lastlogin)
	}
	return s.update(ctx, usersTableName, usernameColumn, username, columns)
}

// DeleteUser deletes a [v1.User]
func (s *Store) DeleteUser(ctx context.Context, username string) error {
	if _, existserr := s.GetUser(ctx, username); existserr != nil {
		return existserr
	}
	return s.del(ctx, usersTableName, usernameColumn, username)
}
