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
	checkres, checkerr := store.GetUser(ctx, user.Username)
	require.NoError(t, checkerr)
	require.NotNil(t, checkres)
	require.Equal(t, user.GetId(), checkres.GetId())
	require.Equal(t, user.GetUsername(), checkres.GetUsername())
	require.Equal(t, user.GetPassword(), checkres.GetPassword())
	require.Equal(t, user.GetLastName(), checkres.GetLastName())
	require.Equal(t, user.GetFirstName(), checkres.GetFirstName())
	require.Equal(t, user.GetEmailAddress(), checkres.GetEmailAddress())
	require.True(t, checkres.GetTimestamps().GetCreated().IsValid())
	require.True(t, checkres.GetTimestamps().GetUpdated().IsValid())
	now := time.Now()
	uerr := store.UpdateUser(ctx, user.Username,
		filters.WithFirstName("bob"),
		filters.WithLastName("smith"),
		filters.WithAvatarURL("avatar"),
		filters.WithEmailAddress("new-email"),
		filters.WithLastLogin(&now),
		filters.WithPassword("newpass"),
	)
	require.NoError(t, uerr)
	gres, gerr := store.GetUser(ctx, user.Username)
	require.NoError(t, gerr)
	require.NotNil(t, gres)
	require.Equal(t, "bob", gres.GetFirstName())
	require.Equal(t, "smith", gres.GetLastName())
	require.Equal(t, "avatar", gres.GetAvatarUrl())
	require.Equal(t, "new-email", gres.GetEmailAddress())
	require.True(t, gres.GetLastLogin().IsValid())
	require.Equal(t, "newpass", gres.GetPassword())

	delerr := store.DeleteUser(ctx, res.Username)
	require.NoError(t, delerr)
	cres, cerr := store.GetUser(ctx, res.Username)
	require.ErrorIs(t, cerr, serrors.ErrNotFound)
	require.Nil(t, cres)
}
