package sqlite

const (
	// tables
	itemsTableName  = "items"
	statusTableName = "status"
	notesTableName  = "notes"
	usersTableName  = "users"

	// columns
	idColumn          = "id"
	statusIDColumn    = "status_id"
	itemIDColumn      = "item_id"
	descriptionColumn = "description"
	nameColumn        = "name"
	colorColumn       = "color"
	createdColumn     = "created"
	updatedColumn     = "updated"
	deleteColumn      = "deleted"
	kindColumn        = "kind"
	noteColumn        = "note_text"
	fnameColumn       = "first_name"
	lnameColumn       = "last_name"
	lastloginColumn   = "last_login"
	usernameColumn    = "username"
	passwordColumn    = "password"
	emailColumn       = "email_address"
	avatarURLColumn   = "avatar_url"
)
