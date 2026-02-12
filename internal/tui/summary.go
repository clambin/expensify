package tui

import (
	"maps"
	"slices"

	"codeberg.org/clambin/bubbles/table"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/clambin/expensify/internal/statements"
	"github.com/clambin/expensify/tcsv"
)

type summaryView struct {
	*table.Table
	SummaryKeyMap
	taggedStatements []statements.TaggedRow
	schema           tcsv.Schema
}

func newSummaryView(keyMap SummaryKeyMap) *summaryView {
	return &summaryView{
		Table: table.NewTable(
			"summary",
			table.Columns{
				{Name: "Tag"},
				{Name: "Amount", Width: 10, RowStyle: table.CellStyle{Style: lipgloss.NewStyle().Align(lipgloss.Right)}},
			},
			nil,
			tableStyles,
			table.DefaultKeyMap(),
		),
		SummaryKeyMap: keyMap,
	}
}

func (sv *summaryView) Init() tea.Cmd {
	return nil
}

func (sv *summaryView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case populateStatementsMsg:
		sv.taggedStatements = msg.taggedStatements
		sv.schema = msg.file.Schema
		sv.SetRows(sv.buildRows())
		return func() tea.Msg { return statusMsg{} }
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, sv.Open):
			filteredStatements := sv.filterStatements(sv.SelectedRow[0].(string))
			return tea.Batch(
				func() tea.Msg {
					return showStatementsMsg{filteredStatements}
				},
				func() tea.Msg { return setActivePaneMsg{statementsPane} },
			)
		default:
			return sv.Table.Update(msg)
		}
	default:
		return sv.Table.Update(msg)
	}
}

func (sv *summaryView) buildRows() []table.Row {
	vals, keys := summarizeStatements(sv.taggedStatements, sv.schema)
	var total float64

	rows := make([]table.Row, 0, 1+len(keys))
	for _, k := range keys {
		rows = append(rows, table.Row{k, alignFloat(vals[k], 10, lipgloss.Right)})
		total += vals[k]
	}
	rows = append(rows, table.Row{
		tableStyles.Header.Render("Total"),
		tableStyles.Header.Render(alignFloat(total, 10, lipgloss.Right)),
	})
	return rows
}

func (sv *summaryView) filterStatements(tag string) []statements.TaggedRow {
	filteredStatements := make([]statements.TaggedRow, 0, len(sv.taggedStatements))
	for _, statement := range sv.taggedStatements {
		if statement.Tag == tag {
			filteredStatements = append(filteredStatements, statement)
		}
	}
	return filteredStatements
}

func summarizeStatements(statements []statements.TaggedRow, schema tcsv.Schema) (map[string]float64, []string) {
	summary := make(map[string]float64)
	tags := make(map[string]struct{})
	for _, statement := range statements {
		summary[statement.Tag] += schema.ToMap(statement.Row)["Amount"].(float64)
		tags[statement.Tag] = struct{}{}
	}
	sortedTags := slices.Collect(maps.Keys(tags))
	slices.Sort(sortedTags)

	return summary, sortedTags
}
