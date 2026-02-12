package tui

import (
	"codeberg.org/clambin/bubbles/table"
	"github.com/clambin/expensify/internal/statements"
	"github.com/clambin/expensify/tcsv"
)

// setActivePaneMsg sets the active pane in the UI
type setActivePaneMsg struct {
	paneID
}

// populateRepoFilesMsg populates the list of statement files
type populateRepoFilesMsg struct {
	files []string
}

// populateStatementsMsg contains the loaded statements and is used to populate the summary and statements views
type populateStatementsMsg struct {
	taggedStatements []statements.TaggedRow
	file             tcsv.File
}

// showStatementsMsg shows the statements in the statements view
type showStatementsMsg struct {
	statements []statements.TaggedRow
}

// setStatementsModeMsg sets the mode of the statements view.
// currently, either list or details is supported
type setStatementsModeMsg struct {
	mode statementsMode
}

// openStatementDetailsMsg opens the details view for the given row
type openStatementDetailsMsg struct {
	taggedRow table.Row
	schema    tcsv.Schema
}
