package tui

import (
	"codeberg.org/clambin/bubbles/table"
	"github.com/clambin/expensify/internal/statements"
	"github.com/clambin/expensify/tcsv"
)

type errorMsg struct {
	err error
}

// openStatementsMsg opens a statement file
type openStatementsMsg struct {
	name string
}

// loadStatementsMsg loads the statements from the file requested by openStatementsMsg
type loadStatementsMsg struct {
	taggedStatements []statements.TaggedRow
	file             tcsv.File
}

// closeStatementsMsg indicates the user has closed the statement file
type closeStatementsMsg struct{}

// openDetailsMsg opens the details view for the given row
type openDetailsMsg struct {
	taggedRow table.Row
	schema    tcsv.Schema
}

// closeDetailsMsg indicates the user has closed the details view
type closeDetailsMsg struct{}
