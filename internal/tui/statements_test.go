package tui

import (
	"testing"

	"codeberg.org/clambin/bubbles/table"
	"github.com/clambin/expensify/csvt"
	"github.com/clambin/expensify/internal/statements"
	"github.com/stretchr/testify/assert"
)

func TestSummarizeStatements(t *testing.T) {
	rows := []statements.TaggedRow{
		{Tag: "Food", Row: csvt.Row{nil, nil, 10.0}},
		{Tag: "Rent", Row: csvt.Row{nil, nil, 100.0}},
		{Tag: "Food", Row: csvt.Row{nil, nil, 5.0}},
	}
	schema := csvt.Schema{
		Columns: []csvt.Column{
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
	sv := newStatementsView(defaultStatementsListKeyMap(), defaultStatementsDetailsKeyMap())
	assert.Equal(t, statementsList, sv.mode)

	// Open details
	msg := openStatementDetailsMsg{
		taggedRow: table.Row{"2024-01-01", "10.0", "Food"},
		schema:    csvt.Schema{},
	}
	sv, _ = sv.Update(msg)
	assert.Equal(t, statementsList, sv.mode)

	// Close details
	sv, _ = sv.Update(setStatementsModeMsg{statementsDetails})
	assert.Equal(t, statementsDetails, sv.mode)
}
