package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"codeberg.org/clambin/bubbles/table"
	"github.com/charmbracelet/x/ansi"
	"github.com/clambin/expensify/csvt"
	"github.com/clambin/expensify/internal/statements"
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
				columns = append(columns, table.Column{
					Name:      col.Label,
					Width:     formatting[col.Label].Width,
					CellStyle: formatting[col.Label].RowStyle,
				})
			}
		}
		columns = append(columns, table.Column{
			Name:      "Tag",
			Width:     formatting["Tag"].Width,
			CellStyle: formatting["Tag"].RowStyle})
		statementColumns[name] = columns
	}
}

// statementsMode indicates if the statementsView is in list or details mode
type statementsMode int

const (
	statementsList statementsMode = iota
	statementsDetails
)

// statementsView combines the list of statements with a detailed view of the selected statement
type statementsView struct {
	statementsDetails statementsDetailsView
	statementsList    statementsListView
	mode              statementsMode
}

func newStatementsView(listKeyMap StatementsListKeyMap, detailsKeyMap StatementsDetailsKeyMap) statementsView {
	return statementsView{
		statementsList: statementsListView{
			Table: table.New().
				Columns(statementColumns["bnp-debit"]).
				Styles(tableStyles),
			StatementsListKeyMap: listKeyMap,
		},
		statementsDetails: statementsDetailsView{
			Table: table.New().
				Columns(table.Columns{{Name: "Field", Width: 15}, {Name: "Value"}}).
				Styles(tableStyles),
			StatementsDetailsKeyMap: detailsKeyMap,
		},
	}
}

func (sv statementsView) Init() tea.Cmd {
	return tea.Batch(
		sv.statementsList.Init(),
		sv.statementsDetails.Init(),
	)
}

func (sv statementsView) Update(msg tea.Msg) (statementsView, tea.Cmd) {
	switch msg := msg.(type) {
	case setStatementsModeMsg:
		sv.mode = msg.mode
		return sv, nil
	case tea.KeyPressMsg:
		var cmd tea.Cmd
		switch sv.mode {
		case statementsList:
			sv.statementsList, cmd = sv.statementsList.Update(msg)
		case statementsDetails:
			sv.statementsDetails, cmd = sv.statementsDetails.Update(msg)
		}
		return sv, cmd
	default:
		cmds := make([]tea.Cmd, 2)
		sv.statementsList, cmds[0] = sv.statementsList.Update(msg)
		sv.statementsDetails, cmds[1] = sv.statementsDetails.Update(msg)
		return sv, tea.Batch(cmds...)
	}
}

func (sv statementsView) View() string {
	switch sv.mode {
	case statementsList:
		return sv.statementsList.View()
	case statementsDetails:
		return sv.statementsDetails.View()
	default:
		return ""
	}
}

func (sv statementsView) SetSize(width, height int) statementsView {
	sv.statementsList = sv.statementsList.SetSize(width, height)
	sv.statementsDetails = sv.statementsDetails.SetSize(width, height)
	return sv
}

func (sv statementsView) ShortHelp() []key.Binding {
	switch sv.mode {
	case statementsList:
		return sv.statementsList.ShortHelp()
	case statementsDetails:
		return sv.statementsDetails.ShortHelp()
	default:
		return nil
	}
}

func (sv statementsView) FullHelp() [][]key.Binding {
	return [][]key.Binding{sv.ShortHelp()}
}

// statementsListView displays a list of statements
type statementsListView struct {
	StatementsListKeyMap
	schema csvt.Schema
	table.Table
}

func (s statementsListView) Update(msg tea.Msg) (statementsListView, tea.Cmd) {
	switch msg := msg.(type) {
	case populateStatementsMsg:
		s.schema = msg.file.Schema
		s.Table = s.Columns(statementColumns[msg.file.SchemaName])
		return s, func() tea.Msg { return setStatementsModeMsg{statementsList} }
	case showStatementsMsg:
		s.Table = s.Rows(buildStatementRows(msg.statements))
		return s, func() tea.Msg { return setStatementsModeMsg{statementsList} }
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, s.Details):
			return s, tea.Batch(
				func() tea.Msg {
					return openStatementDetailsMsg{taggedRow: s.SelectedRow(), schema: s.schema}
				},
				func() tea.Msg { return setStatementsModeMsg{statementsDetails} },
			)
		default:
			var cmd tea.Cmd
			s.Table, cmd = s.Table.Update(msg)
			return s, cmd
		}
	default:
		var cmd tea.Cmd
		s.Table, cmd = s.Table.Update(msg)
		return s, cmd
	}
}

func (s statementsListView) ShortHelp() []key.Binding {
	return []key.Binding{s.Details}
}

func (s statementsListView) FullHelp() [][]key.Binding {
	return [][]key.Binding{s.ShortHelp()}
}

func (s statementsListView) SetSize(width, height int) statementsListView {
	s.Table = s.Size(width, height)
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
	StatementsDetailsKeyMap
	table.Table
}

func (s statementsDetailsView) Update(msg tea.Msg) (statementsDetailsView, tea.Cmd) {
	switch msg := msg.(type) {
	case openStatementDetailsMsg:
		s.Table = s.Rows(buildStatementDetailRows(msg.taggedRow, msg.schema))
		return s, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.Close):
			return s, func() tea.Msg { return setStatementsModeMsg{statementsList} }

		}
	}
	var cmd tea.Cmd
	s.Table, cmd = s.Table.Update(msg)
	return s, cmd
}

func (s statementsDetailsView) ShortHelp() []key.Binding {
	return []key.Binding{s.Close}
}

func (s statementsDetailsView) FullHelp() [][]key.Binding {
	return [][]key.Binding{s.ShortHelp()}
}

func (s statementsDetailsView) SetSize(width, height int) statementsDetailsView {
	s.Table = s.Size(width, height)
	return s
}

func buildStatementDetailRows(taggedRow table.Row, schema csvt.Schema) []table.Row {
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
