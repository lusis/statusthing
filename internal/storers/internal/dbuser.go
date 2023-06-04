package internal

import (
	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers"
	"github.com/lusis/statusthing/internal/validation"

	"google.golang.org/protobuf/proto"
)

// DbUser represents a user stored in a database
type DbUser struct {
	ID           string  `db:"id" goqu:"skipupdate"`
	Username     string  `db:"username"`
	Password     string  `db:"password"`
	FirstName    *string `db:"first_name"`
	LastName     *string `db:"last_name"`
	EmailAddress *string `db:"email_address"`
	LastLogin    *uint64 `db:"last_login"`
	AvatarURL    *string `db:"avatar_url"`
	*DbTimestamps
}

// DbUserFromProto creates a [DBUser] from a [v1.User]
func DbUserFromProto(pbuser *v1.User) (*DbUser, error) {
	if pbuser == nil {
		return nil, serrors.NewError("status", serrors.ErrNilVal)
	}
	id := pbuser.GetId()
	username := pbuser.GetUsername()
	password := pbuser.GetPassword()
	fname := pbuser.GetFirstName()
	lname := pbuser.GetLastName()
	email := pbuser.GetEmailAddress()
	lastlogin := pbuser.GetLastLogin()
	timestamps, err := MakeDbTimestamps(pbuser.GetTimestamps())
	if err != nil {
		return nil, err
	}
	if !validation.ValidString(id) {
		return nil, serrors.NewError("id", serrors.ErrEmptyString)
	}
	if !validation.ValidString(username) {
		return nil, serrors.NewError("username", serrors.ErrEmptyString)
	}
	if !validation.ValidString(password) {
		return nil, serrors.NewError("password", serrors.ErrEmptyString)
	}

	res := &DbUser{
		ID:           id,
		Username:     username,
		Password:     password,
		DbTimestamps: timestamps,
	}
	if validation.ValidString(fname) {
		res.FirstName = &fname
	}
	if validation.ValidString(lname) {
		res.LastName = &lname
	}
	if validation.ValidString(email) {
		res.EmailAddress = &email
	}
	if err := lastlogin.CheckValid(); err == nil {
		res.LastLogin = storers.TsToUInt64Ptr(lastlogin)
	}
	return res, nil
}

// ToProto converts a [v1.User] to a [DBUser]
func (u *DbUser) ToProto() (proto.Message, error) {
	if !validation.ValidString(u.ID) {
		return nil, serrors.NewError("id", serrors.ErrInvalidData)
	}
	if !validation.ValidString(u.Username) {
		return nil, serrors.NewError("username", serrors.ErrInvalidData)
	}
	if !validation.ValidString(u.Password) {
		return nil, serrors.NewError("password", serrors.ErrInvalidData)
	}
	res := &v1.User{
		Id:         u.ID,
		Username:   u.Username,
		Password:   u.Password,
		Timestamps: &v1.Timestamps{},
	}
	// timestamps
	pbcreated := storers.Int64ToTs(int64(u.Created))
	pbupdated := storers.Int64ToTs(int64(u.Updated))

	if pbcreated == nil {
		return nil, serrors.NewError("created", serrors.ErrInvalidData)
	}
	if pbupdated == nil {
		return nil, serrors.NewError("updated", serrors.ErrInvalidData)
	}
	res.Timestamps.Created = pbcreated
	res.Timestamps.Updated = pbupdated

	if u.Deleted != nil {
		res.Timestamps.Deleted = storers.Int64ToTs(int64(*u.Deleted))
	}
	if u.LastLogin != nil {
		res.LastLogin = storers.Int64ToTs(int64(*u.LastLogin))
	}
	if u.FirstName != nil {
		res.FirstName = *u.FirstName
	}
	if u.LastName != nil {
		res.LastName = *u.LastName
	}
	if u.EmailAddress != nil {
		res.EmailAddress = *u.EmailAddress
	}
	return res, nil
}
