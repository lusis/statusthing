// Package serrors contains all custom error types
package serrors

import "fmt"

// ErrEmptyString is a custom error when an empty string is passed and is not valid
// this error is generally returned when a string parameter cannot be empty
var ErrEmptyString = fmt.Errorf("string cannot be empty")

// ErrAlreadySet is a custom error when a value is already set and cannot be overwritten
var ErrAlreadySet = fmt.Errorf("value already set")

// ErrEmptyEnum is a custom error when an enum value is the zero value
// this error is generally returned when an enum provided is that enum zero value and not allowed
var ErrEmptyEnum = fmt.Errorf("unknown not allowed")

// ErrNilVal is a custom error when a nil value is not allowed
// this error is generally returned when a function disallows a nil value
var ErrNilVal = fmt.Errorf("nil not allowed")

// ErrNotFound is a custom error when a value is not found
// this error is generally returned when the system is unable to find some resource
var ErrNotFound = fmt.Errorf("not found")

// ErrStoreUnavailable is a custom error when a store is unavailable
// this error is generally returned when a data store is not available at request time
var ErrStoreUnavailable = fmt.Errorf("store unavailable")

// ErrAtLeastOne is a custom error when a slice requires at least one entry
// this error may be returned when at least filter is required for a query
var ErrAtLeastOne = fmt.Errorf("at least one value is required")

// ErrNotImplemented is a custom error when an interface function has not been implemented
var ErrNotImplemented = fmt.Errorf("not implemented")

// ErrInvalidData is a custom error when data is in a corrupt state
// this error is generally returned when data is a store is invalid or corrupt for some reason
// this should be treated with urgency
var ErrInvalidData = fmt.Errorf("invalid data")

// ErrInUse is a custom error when a resource is in use
// this error is generally returned when attempting to delete a Status entry when it is in use by an Item
var ErrInUse = fmt.Errorf("in use")

// ErrMissingTimestamp is the error when a Timestamp specific field is required
var ErrMissingTimestamp = fmt.Errorf("timestamp required")

// ErrUnrecoverable is the error when something has ABENDed in an unsafe to continue way
var ErrUnrecoverable = fmt.Errorf("unrecoverable")

// ErrUnexpectedRows is the error when a db query affects more rows than expected
var ErrUnexpectedRows = fmt.Errorf("unexpected rows affect")

// ErrMissingCredentials is the error when something expects credentials (i.e. a database connection string or api call)
var ErrMissingCredentials = fmt.Errorf("missing credentials")
