package services

import (
	"context"
	"testing"

	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers/memdb"
	"github.com/stretchr/testify/require"
)

func TestUsers(t *testing.T) {
	store, storerr := memdb.New()
	require.NoError(t, storerr)
	require.NotNil(t, store)
	svc, serr := NewStatusThingService(store)
	require.NoError(t, serr)
	require.NotNil(t, svc)
	ctx := context.TODO()

	u, uerr := svc.AddUser(ctx, t.Name(), "password1", "test@test.com")
	require.NoError(t, uerr)
	require.NotNil(t, u)
	require.Equal(t, t.Name(), u.GetUsername())
	require.Equal(t, "test@test.com", u.GetEmailAddress())
	require.NotEmpty(t, u.GetPassword())

	checkpass, checkerr := svc.CheckPassword(ctx, t.Name(), "password1")
	require.NoError(t, checkerr)
	require.NotNil(t, checkpass)
	require.Equal(t, u, checkpass)

	changerr := svc.ChangePassword(ctx, u.GetUsername(), "password1", "otherpassword")
	require.NoError(t, changerr)

	againpass, againerr := svc.CheckPassword(ctx, u.GetUsername(), "otherpassword")
	require.NoError(t, againerr)
	require.NotNil(t, againpass)
	require.NotEqual(t, u.GetPassword(), againpass.GetPassword())

	gres, gerr := svc.GetUser(ctx, againpass.GetUsername())
	require.NoError(t, gerr)
	require.NotNil(t, gres)
	require.Equal(t, againpass, gres)

	badpass, badpasserr := svc.CheckPassword(ctx, gres.GetUsername(), "not the right password yall")
	require.ErrorIs(t, badpasserr, serrors.ErrInvalidPassword)
	require.Nil(t, badpass)

	delerr := svc.RemoveUser(ctx, gres.GetUsername())
	require.NoError(t, delerr)

	fres, ferr := svc.GetUser(ctx, gres.GetUsername())
	require.ErrorIs(t, ferr, serrors.ErrNotFound)
	require.Nil(t, fres)
}
