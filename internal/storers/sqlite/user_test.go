package sqlite

import (
	"context"
	"testing"
	"time"

	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/testutils"

	"github.com/stretchr/testify/require"
)

func TestUserLifeCycle(t *testing.T) {
	db, err := makeTestdb(t, ":memory:")
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()
	user := &v1.User{
		Id:           t.Name() + "_id",
		Username:     t.Name() + "_username",
		Password:     t.Name() + "_password",
		FirstName:    t.Name() + "_fname",
		LastName:     t.Name() + "_lname",
		EmailAddress: t.Name() + "_email",
		Timestamps:   testutils.MakeTimestamps(false),
	}
	store, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, store)
	ctx := context.TODO()
	res, err := store.StoreUser(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, res)
	now := time.Now()
	uerr := store.UpdateUser(ctx, user.Username,
		filters.WithFirstName("bob"),
		filters.WithLastName("smith"),
		filters.WithAvatarURL("avatar"),
		filters.WithEmailAddress("new-email"),
		filters.WithLastLogin(&now),
	)
	require.NoError(t, uerr)
	delerr := store.DeleteUser(ctx, res.Username)
	require.NoError(t, delerr)
	cres, cerr := store.GetUser(ctx, res.Username)
	require.ErrorIs(t, cerr, serrors.ErrNotFound)
	require.Nil(t, cres)
}
