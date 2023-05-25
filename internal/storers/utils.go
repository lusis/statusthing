package storers

import (
	"time"

	"github.com/lusis/statusthing/internal/validation"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TsToInt64 - Consistent timestamp -> datastore mapping
func TsToInt64(ts *timestamppb.Timestamp) int64 {
	if ts == nil {
		return 0
	}
	return ts.AsTime().UnixNano()
}

// Int64ToTs - Consistent datastore -> timestamp mapping
func Int64ToTs(i int64) *timestamppb.Timestamp {
	// zero int val is a nil for us
	if i == 0 {
		return nil
	}
	return timestamppb.New(time.Unix(0, i).UTC())
}

// TsToUInt64 returns a timestamppb.Timestamp as a uint64
func TsToUInt64(ts *timestamppb.Timestamp) uint64 {
	def := uint64(0)
	if ts != nil {
		def = uint64(ts.AsTime().UnixNano())
	}
	return def
}

// TsToUInt64Ptr returns a timestamppb.Timestamp as a pointer uint64
// this is used in place of sql.Null* types
func TsToUInt64Ptr(ts *timestamppb.Timestamp) *uint64 {
	def := uint64(0)
	if ts != nil {
		def = uint64(ts.AsTime().UnixNano())
	}
	return &def
}

// StringPtr returns a string ptr as pointers must be to addressables
// and we don't want to have to create variables for this kind of stuff
// these are used in place of sql.Null* types
func StringPtr(s string) *string {
	if validation.ValidString(s) {
		return &s
	}
	return nil
}
