package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"codeberg.org/clambin/bubbles/table"
	"github.com/clambin/expensify/internal/statements"
	"github.com/clambin/expensify/tcsv"
	"github.com/stretchr/testify/assert"
)

func TestSummarizeStatements(t *testing.T) {
	rows := []statements.TaggedRow{
		{Tag: "Food", Row: tcsv.Row{nil, nil, 10.0}},
		{Tag: "Rent", Row: tcsv.Row{nil, nil, 100.0}},
		{Tag: "Food", Row: tcsv.Row{nil, nil, 5.0}},
	}
	schema := tcsv.Schema{
		Columns: []tcsv.Column{
			{Label: "Date"},
			{Label: "Other"},
			{Label: "Amount"},
		},
	}

	summary, tags := summarizeStatements(rows, schema)

	assert.Equal(t, []string{"Food", "Rent"}, tags)
	assert.Equal(t, 15.0, summary["Food"])
	assert.Equal(t, 100.0, summary["Rent"])
}

func TestStatementsView_DetailsToggle(t *testing.T) {
	var sv tea.Model = newStatementsView(defaultStatementsListKeyMap(), defaultStatementsDetailsKeyMap())
	assert.Equal(t, statementsList, sv.(statementsView).mode)

	// Open details
	msg := openStatementDetailsMsg{
		taggedRow: table.Row{"2024-01-01", "10.0", "Food"},
		schema:    tcsv.Schema{},
	}
	sv, _ = sv.Update(msg)
	assert.Equal(t, statementsList, sv.(statementsView).mode)

	// Close details
	sv, _ = sv.Update(setStatementsModeMsg{statementsDetails})
	assert.Equal(t, statementsDetails, sv.(statementsView).mode)
}
