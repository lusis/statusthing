// nolint: revive
package unimplemented

import (
	"context"

	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"

	"github.com/lusis/statusthing/internal/filters"
)

type UserStore struct{}

// StoreUser stores the provied [v1.User]
func (us *UserStore) StoreUser(ctx context.Context, user *v1.User) (*v1.User, error) {
	panic("not implemented") // TODO: Implement
}

// GetUser gets a [v1.User] by username
func (us *UserStore) GetUser(ctx context.Context, username string) (*v1.User, error) {
	panic("not implemented") // TODO: Implement
}

// FindUsers finds users
func (us *UserStore) FindUsers(ctx context.Context, opts ...filters.FilterOption) ([]*v1.User, error) {
	panic("not implemented") // TODO: Implement
}

// UpdateUser updates a [v1.User]
func (us *UserStore) UpdateUser(ctx context.Context, userID string, opts ...filters.FilterOption) error {
	panic("not implemented") // TODO: Implement
}

// DeleteUser deletes a [v1.User]
func (us *UserStore) DeleteUser(ctx context.Context, userID string) error {
	panic("not implemented") // TODO: Implement
}
