package tui

import (
	"maps"
	"slices"

	"codeberg.org/clambin/bubbles/table"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/clambin/expensify/bubbles/statusbar"
	"github.com/clambin/expensify/internal/statements"
	"github.com/clambin/expensify/tcsv"
)

var _ pane = summaryView{}

type summaryView struct {
	*table.Table
	SummaryKeyMap
	taggedStatements []statements.TaggedRow
	schema           tcsv.Schema
}

func newSummaryView(keyMap SummaryKeyMap) tea.Model {
	return summaryView{
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

func (sv summaryView) Init() tea.Cmd {
	return nil
}

func (sv summaryView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case populateStatementsMsg:
		sv.taggedStatements = msg.taggedStatements
		sv.schema = msg.file.Schema
		sv.SetRows(buildRows(sv))
		return sv, func() tea.Msg { return statusbar.Msg{} }
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, sv.Open):
			return sv, tea.Batch(
				func() tea.Msg {
					return showStatementsMsg{filterStatements(sv.SelectedRow[0].(string), sv.taggedStatements)}
				},
				func() tea.Msg { return setActivePaneMsg{statementsPane} },
			)
		default:
			return sv, sv.Table.Update(msg)
		}
	default:
		return sv, sv.Table.Update(msg)
	}
}

func (sv summaryView) SetSize(width, height int) tea.Model {
	sv.Table.SetSize(width, height)
	return sv
}

func buildRows(sv summaryView) []table.Row {
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

func filterStatements(tag string, taggedStatements []statements.TaggedRow) []statements.TaggedRow {
	filteredStatements := make([]statements.TaggedRow, 0, len(taggedStatements))
	for _, statement := range taggedStatements {
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
