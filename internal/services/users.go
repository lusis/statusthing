// nolint: revive
package services

import (
	"context"

	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/validation"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/alexedwards/argon2id"
	"github.com/segmentio/ksuid"
)

// AddUser creates a new user
func (sts *StatusThingService) AddUser(ctx context.Context, username string, password string, emailAddress string, opts ...filters.FilterOption) (*v1.User, error) {
	if sts.store == nil {
		return nil, serrors.NewError("store", serrors.ErrStoreUnavailable)
	}

	if !validation.ValidString(username) {
		return nil, serrors.NewError("username", serrors.ErrEmptyString)
	}
	if !validation.ValidString(password) {
		return nil, serrors.NewError("password", serrors.ErrEmptyString)
	}
	if !validation.ValidString(emailAddress) {
		return nil, serrors.NewError("email_address", serrors.ErrEmptyString)
	}
	f, err := filters.New(opts...)
	if err != nil {
		return nil, err
	}
	id := ksuid.New().String()
	if validation.ValidString(f.UserID()) {
		id = f.UserID()
	}
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, serrors.NewWrappedError("hashed-password", serrors.ErrUnrecoverable, err)
	}
	u := &v1.User{
		Id:           id,
		Username:     username,
		Password:     hash,
		EmailAddress: emailAddress,
		Timestamps:   makeTsNow(),
	}

	if validation.ValidString(f.FirstName()) {
		u.FirstName = f.FirstName()
	}
	if validation.ValidString(f.LastName()) {
		u.LastName = f.LastName()
	}
	if f.LastLogin() != nil {
		u.LastLogin = timestamppb.New(*f.LastLogin())
		if !u.LastLogin.IsValid() {
			return nil, serrors.NewError("lastlogin", serrors.ErrUnrecoverable)
		}
	}
	return sts.store.StoreUser(ctx, u)
}

// GetUser gets a user by username
func (sts *StatusThingService) GetUser(ctx context.Context, username string) (*v1.User, error) {
	if sts.store == nil {
		return nil, serrors.NewError("store", serrors.ErrStoreUnavailable)
	}
	if !validation.ValidString(username) {
		return nil, serrors.NewError("username", serrors.ErrEmptyString)
	}
	return sts.store.GetUser(ctx, username)
}

// CheckPassword validates the password for the supplied userid
func (sts *StatusThingService) CheckPassword(ctx context.Context, username string, password string) (*v1.User, error) {
	if sts.store == nil {
		return nil, serrors.NewError("store", serrors.ErrStoreUnavailable)
	}
	if !validation.ValidString(username) {
		return nil, serrors.NewError("username", serrors.ErrEmptyString)
	}
	if !validation.ValidString(password) {
		return nil, serrors.NewError("password", serrors.ErrEmptyEnum)
	}
	return sts.checkPassword(ctx, username, password)
}

// ChangePassword changes the password
func (sts *StatusThingService) ChangePassword(ctx context.Context, username string, currPass string, newPass string) error {
	if sts.store == nil {
		return serrors.NewError("store", serrors.ErrStoreUnavailable)
	}
	if !validation.ValidString(username) {
		return serrors.NewError("username", serrors.ErrEmptyString)
	}
	if !validation.ValidString(currPass) {
		return serrors.NewError("current password", serrors.ErrEmptyEnum)
	}
	if !validation.ValidString(newPass) {
		return serrors.NewError("new password", serrors.ErrEmptyEnum)
	}

	if _, err := sts.checkPassword(ctx, username, currPass); err != nil {
		return err
	}
	hash, err := argon2id.CreateHash(newPass, argon2id.DefaultParams)
	if err != nil {
		return serrors.NewWrappedError("hashed-password", serrors.ErrUnrecoverable, err)
	}
	return sts.EditUser(ctx, username, filters.WithPassword(hash))
}

// EditUser edits the user
func (sts *StatusThingService) EditUser(ctx context.Context, username string, opts ...filters.FilterOption) error {
	return sts.store.UpdateUser(ctx, username, opts...)
}

// RemoveUser removes the user
func (sts *StatusThingService) RemoveUser(ctx context.Context, username string) error {
	return sts.store.DeleteUser(ctx, username)
}

func (sts *StatusThingService) checkPassword(ctx context.Context, username, providedPassword string) (*v1.User, error) {
	if sts.store == nil {
		return nil, serrors.NewError("store", serrors.ErrStoreUnavailable)
	}

	res, err := sts.store.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	match, _, err := argon2id.CheckHash(providedPassword, res.GetPassword())
	if err != nil {
		return nil, serrors.NewWrappedError("password-check", serrors.ErrInvalidPassword, err)
	}
	if !match {
		return nil, serrors.ErrInvalidPassword
	}
	return res, nil
}
