package internal

import (
	"testing"

	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers"
	"github.com/lusis/statusthing/internal/testutils"
	"github.com/lusis/statusthing/internal/validation"

	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type userTestCase struct {
	db      *DbUser
	pb      *v1.User
	errtext string
	err     error
}

func TestUserFromProto(t *testing.T) {
	t.Parallel()

	t.Run("nil-check", func(t *testing.T) {
		s, serr := DbUserFromProto(nil)
		require.ErrorIs(t, serr, serrors.ErrNilVal)
		require.Nil(t, s)
	})

	// t.Name inside this map refers to the name of the parent test not the iteration.
	// It's intentional since we're not worried about conflict here
	testcases := map[string]userTestCase{
		"happy-path": {},
		"missing-id": {
			pb:      &v1.User{Id: "", Timestamps: testutils.MakeTimestamps(false)},
			err:     serrors.ErrEmptyString,
			errtext: "id",
		},
		"missing-username": {
			pb:      &v1.User{Id: t.Name(), Password: t.Name() + "password", Timestamps: testutils.MakeTimestamps(false)},
			err:     serrors.ErrEmptyString,
			errtext: "username",
		},
		"missing-password": {
			pb:      &v1.User{Id: t.Name(), Username: t.Name() + "username", Timestamps: testutils.MakeTimestamps(false)},
			err:     serrors.ErrEmptyString,
			errtext: "password",
		},
		"missing-timestamps": {
			pb:      &v1.User{Id: t.Name(), Username: t.Name() + "username", Password: t.Name() + "password"},
			err:     serrors.ErrMissingTimestamp,
			errtext: "timestamps",
		},
		"missing-created": {
			pb: &v1.User{Id: t.Name(), Username: t.Name() + "username", Password: t.Name() + "password", Timestamps: &v1.Timestamps{
				Updated: timestamppb.Now(),
			}},
			err:     serrors.ErrMissingTimestamp,
			errtext: "created",
		},
		"missing-updated": {
			pb: &v1.User{Id: t.Name(), Username: t.Name() + "username", Password: t.Name() + "password", Timestamps: &v1.Timestamps{
				Created: timestamppb.Now(),
			}},
			err:     serrors.ErrMissingTimestamp,
			errtext: "updated",
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			pb := tc.pb
			if pb == nil {
				pb = testutils.MakeUser(t.Name())
				pb.FirstName = t.Name() + "_firstname"
				pb.LastName = t.Name() + "_lastname"
				pb.Timestamps = testutils.MakeTimestamps(true)
				pb.LastLogin = timestamppb.Now()
				pb.EmailAddress = t.Name() + "_emailaddress"
			}
			s, serr := DbUserFromProto(pb)
			if tc.err != nil {
				require.ErrorIs(t, serr, tc.err)
				require.Nil(t, s)
				if validation.ValidString(tc.errtext) {
					require.ErrorContains(t, serr, tc.errtext)
				}
			} else {
				require.NoError(t, serr)
				require.NotNil(t, s)
				require.Equal(t, pb.GetId(), s.ID)
				require.Equal(t, pb.GetUsername(), s.Username)
				require.Equal(t, pb.GetPassword(), s.Password)
				require.Equal(t, pb.GetFirstName(), *s.FirstName)
				require.Equal(t, pb.GetLastName(), *s.LastName)
				require.Equal(t, pb.GetEmailAddress(), *s.EmailAddress)
				require.NotZero(t, s.Created)
				require.NotZero(t, s.Updated)
				if tc.pb.GetTimestamps().GetDeleted().IsValid() {
					require.NotNil(t, s.Deleted)
				}
				if tc.pb.GetLastLogin().IsValid() {
					require.NotNil(t, s.LastLogin)
				}
			}
		})
	}
}

func TestUserToProto(t *testing.T) {
	t.Parallel()
	// t.Name inside this map refers to the name of the parent test not the iteration.
	// It's intentional since we're not worried about conflict here
	testcases := map[string]userTestCase{
		"happy-path": {},
		"missing-id": {
			db:      &DbUser{ID: ""},
			err:     serrors.ErrInvalidData,
			errtext: "id",
		},
		"missing-username": {
			db:      &DbUser{ID: t.Name()},
			err:     serrors.ErrInvalidData,
			errtext: "username",
		},
		"missing-password": {
			db:      &DbUser{ID: t.Name(), Username: t.Name()},
			err:     serrors.ErrInvalidData,
			errtext: "password",
		},
		"missing-created": {
			db: &DbUser{
				ID:           t.Name(),
				Username:     t.Name(),
				Password:     t.Name(),
				DbTimestamps: &DbTimestamps{Updated: uint64(storers.TsToInt64(timestamppb.Now()))},
			},
			err:     serrors.ErrInvalidData,
			errtext: "created",
		},
		"missing-updated": {
			db: &DbUser{
				ID:           t.Name(),
				Username:     t.Name(),
				Password:     t.Name(),
				DbTimestamps: &DbTimestamps{Created: uint64(storers.TsToInt64(timestamppb.Now()))},
			},
			err:     serrors.ErrInvalidData,
			errtext: "updated",
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			db := tc.db
			if db == nil {
				firstname := t.Name() + "firstname"
				lastname := t.Name() + "lastname"
				email := t.Name() + "email"
				db = &DbUser{
					ID:           t.Name() + "id",
					Username:     t.Name() + "username",
					Password:     t.Name() + "password",
					FirstName:    &firstname,
					LastName:     &lastname,
					EmailAddress: &email,
					LastLogin:    storers.TsToUInt64Ptr(timestamppb.Now()),
					DbTimestamps: &DbTimestamps{
						Created: storers.TsToUInt64(timestamppb.Now()),
						Updated: storers.TsToUInt64(timestamppb.Now()),
						Deleted: storers.TsToUInt64Ptr(timestamppb.Now())},
				}
			}
			res, serr := db.ToProto()

			if tc.err != nil {
				require.ErrorIs(t, serr, tc.err)
				require.Nil(t, res)
				if validation.ValidString(tc.errtext) {
					require.ErrorContains(t, serr, tc.errtext)
				}
			} else {
				s, ok := res.(*v1.User)
				require.True(t, ok)
				require.NoError(t, serr)
				require.NotNil(t, s)
				require.Equal(t, db.ID, s.GetId())
				require.Equal(t, db.Username, s.GetUsername())
				require.Equal(t, db.Password, s.GetPassword())
				require.Equal(t, *db.FirstName, s.GetFirstName())
				require.Equal(t, *db.LastName, s.GetLastName())
				require.Equal(t, *db.EmailAddress, s.GetEmailAddress())
				require.True(t, s.GetTimestamps().GetCreated().IsValid())
				require.True(t, s.GetTimestamps().GetUpdated().IsValid())
				if db.Deleted != nil {
					require.True(t, s.GetTimestamps().GetDeleted().IsValid())
				}
				if db.LastLogin != nil {
					require.True(t, s.GetLastLogin().IsValid())
				}
			}
		})
	}
}
