package sqlite

import "fmt"

const (
	// tables
	itemsTableName  = "items"
	statusTableName = "status"
	notesTableName  = "notes"

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
)

// notes:
// we're going to use FKs here with sqlite
// I have .... opinions ... on FKs but those mostly deal with MySQL as well as FKs with db migrations in general
// so we'll use them here
var (
	// sqlite doesn't require escaped column names with backticks like mysql so we win the string literal game
	tmplCreateStatusTable = `
	CREATE TABLE IF NOT EXISTS %s 
	(
		id VARCHAR(191) PRIMARY KEY,
		name VARCHAR(191) NOT NULL UNIQUE,
		kind VARCHAR(191) NOT NULL,
		description VARCHAR(191) DEFAULT NULL,
		color VARCHAR(191) DEFAULT NULL,
		created INT NOT NULL,
		updated INT NOT NULL,
		deleted INT DEFAULT NULL
	)`
	tmplCreateItemsTable = `
	CREATE TABLE IF NOT EXISTS %s
	(
		id VARCHAR(191) PRIMARY KEY,
		name VARCHAR(191) NOT NULL UNIQUE,
		description VARCHAR(191) DEFAULT NULL,
		status_id VARCHAR(191) DEFAULT NULL,
		created INT NOT NULL,
		updated INT NOT NULL,
		deleted INT DEFAULT NULL,
		FOREIGN KEY(status_id) REFERENCES status(id)
	)
	`
	tmplCreateNotesTable = `
	CREATE TABLE IF NOT EXISTS %s
	(
		id VARCHAR(191) PRIMARY KEY,
		note_text VARCHAR(191) NOT NULL UNIQUE,
		item_id VARCHAR(191) NOT NULL,
		created INT NOT NULL,
		updated INT NOT NULL,
		deleted INT DEFAULT NULL,
		FOREIGN KEY(item_id) REFERENCES items(id)
	)
	`
	stmtCreateStatusTable = fmt.Sprintf(tmplCreateStatusTable, statusTableName)
	stmtCreateItemsTable  = fmt.Sprintf(tmplCreateItemsTable, itemsTableName)
	stmtCreateNotesTable  = fmt.Sprintf(tmplCreateNotesTable, notesTableName)
)
