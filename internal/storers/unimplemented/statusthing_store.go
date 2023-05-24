package unimplemented

// StatusThingStore stores [statusthingv1.StatusThing]
type StatusThingStore struct {
	*NoteStorer
	*StatusStore
	*ItemStore
}
