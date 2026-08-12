package tui

import (
	"sync/atomic"

	tea "charm.land/bubbletea/v2"
	"codeberg.org/clambin/bubbles/statusbar"
	"codeberg.org/clambin/bubbles/table"
	"github.com/clambin/expensify/csvt"
	"github.com/clambin/expensify/internal/statements"
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
	file             csvt.File
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
	schema    csvt.Schema
}

// setStatusBar updates the status bar
type setStatusBarMsg struct {
	text    string
	options []statusbar.SetStatusOption
	seqNo   uint64
}

func setStatusBarCmd(text string, options ...statusbar.SetStatusOption) tea.Cmd {
	return func() tea.Msg {
		return setStatusBarMsg{seqNo: nextSeqNo(), text: text, options: options}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// statusBarSeqNo keeps track of the current status bar sequence number
// so that Application can ignore messages that arrive out of order.
var statusBarSeqNo atomic.Uint64

func nextSeqNo() uint64 {
	return statusBarSeqNo.Add(1)
}

func isCurrentSeqNo(seqNo uint64) bool {
	return seqNo == statusBarSeqNo.Load()
}
