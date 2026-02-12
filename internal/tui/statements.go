package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"codeberg.org/clambin/bubbles/table"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/clambin/expensify/internal/statements"
	"github.com/clambin/expensify/tcsv"
)

var statementColumns map[string][]table.Column

var formatting = map[string]struct {
	RowStyle table.CellStyle
	Width    int
}{
	"Date":            {Width: 10},
	"Amount":          {Width: 11, RowStyle: table.CellStyle{Style: lipgloss.NewStyle().Align(lipgloss.Right)}},
	"Type":            {Width: 14},
	"Counterpart":     {Width: 17},
	"CounterpartName": {Width: 18},
	"Message":         {Width: 25},
	"Tag":             {Width: 10},
}

func init() {
	statementColumns = make(map[string][]table.Column)
	for name, schema := range statements.Schemas {
		columns := make([]table.Column, 0, len(schema.Columns)+1)
		for _, col := range schema.Columns {
			if col.Label != "" {
				columns = append(columns, table.Column{Name: col.Label, Width: formatting[col.Label].Width, RowStyle: formatting[col.Label].RowStyle})
			}
		}
		columns = append(columns, table.Column{Name: "Tag", Width: formatting["Tag"].Width, RowStyle: formatting["Tag"].RowStyle})
		statementColumns[name] = columns
	}
}

// statementsMode indicates if the statementsView is in list or details mode
type statementsMode int

const (
	statementsList statementsMode = iota
	statementsDetails
)

var (
	_ pane = statementsView{}
	_ pane = statementsListView{}
	_ pane = statementsDetailsView{}
)

// statementsView combines the list of statements with a detailed view of the selected statement
type statementsView struct {
	views map[statementsMode]tea.Model
	mode  statementsMode
}

func newStatementsView(listKeyMap StatementsListKeyMap, detailsKeyMap StatementsDetailsKeyMap) statementsView {
	return statementsView{
		views: map[statementsMode]tea.Model{
			statementsList: &statementsListView{
				Table: table.NewTable(
					"statements",
					statementColumns["bnp-debit"],
					nil,
					tableStyles,
					table.DefaultKeyMap(),
				),
				StatementsListKeyMap: listKeyMap,
			},
			statementsDetails: &statementsDetailsView{
				Table: table.NewTable(
					"statement details",
					table.Columns{
						{Name: "Field", Width: 15},
						{Name: "Value"},
					},
					nil,
					tableStyles, table.DefaultKeyMap(),
				),
				StatementsDetailsKeyMap: detailsKeyMap,
			},
		},
	}
}

func (sv statementsView) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(sv.views))
	for _, v := range sv.views {
		cmds = append(cmds, v.Init())
	}
	return tea.Batch(cmds...)
}

func (sv statementsView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case setStatementsModeMsg:
		sv.mode = msg.mode
		return sv, nil
	case tea.KeyMsg:
		var cmd tea.Cmd
		sv.views[sv.mode], cmd = sv.views[sv.mode].Update(msg)
		return sv, cmd
	default:
		cmds := make([]tea.Cmd, 0, len(sv.views))
		for v := range sv.views {
			var cmd tea.Cmd
			sv.views[v], cmd = sv.views[v].Update(msg)
			cmds = append(cmds, cmd)
		}
		return sv, tea.Batch(cmds...)
	}
}

func (sv statementsView) View() string {
	return sv.views[sv.mode].View()
}

func (sv statementsView) SetSize(width, height int) tea.Model {
	for p := range sv.views {
		sv.views[p] = sv.views[p].(pane).SetSize(width, height)
	}
	return sv
}

func (sv statementsView) ShortHelp() []key.Binding {
	return sv.views[sv.mode].(pane).ShortHelp()
}

func (sv statementsView) FullHelp() [][]key.Binding {
	return [][]key.Binding{sv.ShortHelp()}
}

// statementsListView displays a list of statements
type statementsListView struct {
	*table.Table
	StatementsListKeyMap
	schema tcsv.Schema
}

func (s statementsListView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case populateStatementsMsg:
		s.schema = msg.file.Schema
		s.SetColumns(statementColumns[msg.file.SchemaName])
		return s, func() tea.Msg { return setStatementsModeMsg{statementsList} }
	case showStatementsMsg:
		s.SetRows(buildStatementRows(msg.statements))
		return s, func() tea.Msg { return setStatementsModeMsg{statementsList} }
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.Details):
			return s, tea.Batch(
				func() tea.Msg { return openStatementDetailsMsg{taggedRow: s.SelectedRow, schema: s.schema} },
				func() tea.Msg { return setStatementsModeMsg{statementsDetails} },
			)
		default:
			return s, s.Table.Update(msg)
		}
	default:
		return s, s.Table.Update(msg)
	}
}

func (s statementsListView) ShortHelp() []key.Binding {
	return []key.Binding{s.Details}
}

func (s statementsListView) FullHelp() [][]key.Binding {
	return [][]key.Binding{s.ShortHelp()}
}

func (s statementsListView) SetSize(width, height int) tea.Model {
	s.Table.SetSize(width, height)
	return s
}

func buildStatementRows(statements []statements.TaggedRow) []table.Row {
	rows := make([]table.Row, len(statements))
	for i, row := range statements {
		rows[i] = make(table.Row, len(row.Row), len(row.Row)+1)
		for j, cell := range row.Row {
			var value string
			switch cell := cell.(type) {
			case time.Time:
				value = cell.Format("2006-01-02")
			case float64:
				value = alignFloat(cell, 10, lipgloss.Right)
			default:
				value = fmt.Sprintf("%v", cell)
			}
			rows[i][j] = value
		}
		rows[i] = append(rows[i], row.Tag)
	}
	return rows
}

// statementsDetailsView displays the details of a single statement
type statementsDetailsView struct {
	*table.Table
	StatementsDetailsKeyMap
}

func (s statementsDetailsView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case openStatementDetailsMsg:
		s.SetRows(buildStatementDetailRows(msg.taggedRow, msg.schema))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.Close):
			return s, func() tea.Msg { return setStatementsModeMsg{statementsList} }
		}
	}
	return s, s.Table.Update(msg)
}

func (s statementsDetailsView) ShortHelp() []key.Binding {
	return []key.Binding{s.Close}
}

func (s statementsDetailsView) FullHelp() [][]key.Binding {
	return [][]key.Binding{s.ShortHelp()}
}

func (s statementsDetailsView) SetSize(width, height int) tea.Model {
	s.Table.SetSize(width, height)
	return s
}

func buildStatementDetailRows(taggedRow table.Row, schema tcsv.Schema) []table.Row {
	rows := make([]table.Row, 0, len(taggedRow)+1)
	var idx int
	for _, col := range schema.Columns {
		if col.Label != "" {
			rows = append(rows, table.Row{col.Label, strings.TrimSpace(ansi.Strip(taggedRow[idx].(string)))})
			idx++
		}
	}
	return append(rows, table.Row{"Tag", taggedRow[len(taggedRow)-1]})
}

func alignString(s string, width int, position lipgloss.Position) string {
	return lipgloss.NewStyle().Align(position).Width(width).Render(s)
}

func alignFloat(f float64, width int, position lipgloss.Position) string {
	return alignString(strconv.FormatFloat(f, 'f', 2, 64), width, position)
}
