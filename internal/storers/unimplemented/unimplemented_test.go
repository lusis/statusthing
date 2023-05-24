package unimplemented

import (
	"testing"

	"github.com/lusis/statusthing/internal/storers"
	"github.com/stretchr/testify/require"
)

func TestImplements(t *testing.T) {
	require.Implements(t, (*storers.StatusStorer)(nil), new(StatusStore), "unimplemented custom status store should satisfy interface")
	require.Implements(t, (*storers.NoteStorer)(nil), new(NoteStorer), "unimplemented note store should satisfy interface")
	require.Implements(t, (*storers.ItemStorer)(nil), new(ItemStore), "unimplemented status thing store should sastify interface")
	require.Implements(t, (*storers.StatusThingStorer)(nil), new(StatusThingStore), "unimplemented status thing store should sastify interface")
}
