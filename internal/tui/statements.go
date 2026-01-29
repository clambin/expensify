package tui

import (
	"bytes"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"time"

	"codeberg.org/clambin/bubbles/frame"
	"codeberg.org/clambin/bubbles/table"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/statements"
	"github.com/clambin/expensify/tcsv"
)

var (
	borderStyles = map[bool]frame.Styles{
		true:  selectedFrameStyles,
		false: frameStyles,
	}
)

var _ help.KeyMap = (*bodyView)(nil)

type bodyView struct {
	statementsView *statementsView
	summaryView    *summaryView
	keyMap         BodyKeyMap
	activeTab      int
}

func newBodyView(keyMap BodyKeyMap) *bodyView {
	return &bodyView{keyMap: keyMap}
}

func (b *bodyView) Init() tea.Cmd {
	return tea.Batch(b.statementsView.Init(), b.summaryView.Init())
}

func (b *bodyView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keyMap.NextPane):
			b.activeTab = (b.activeTab + 1) % 2
			return nil
		case key.Matches(msg, b.keyMap.PrevPane):
			b.activeTab = (b.activeTab - 1 + 2) % 2
			return nil
		case key.Matches(msg, b.keyMap.Close):
			return func() tea.Msg { return closeStatementsMsg{} }
		default:
			return b.updateActivePane(msg)
		}
	default:
		return b.updateActivePane(msg)
	}
}

func (b *bodyView) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Left,
		frame.Draw("statements", lipgloss.Center, b.statementsView.View(), borderStyles[b.activeTab == 0]),
		frame.Draw("summary", lipgloss.Center, b.summaryView.View(), borderStyles[b.activeTab == 1]),
	)
}

func (b *bodyView) updateActivePane(msg tea.Msg) tea.Cmd {
	switch b.activeTab {
	case 0:
		return b.statementsView.Update(msg)
	case 1:
		return b.summaryView.Update(msg)
	default:
		return nil
	}
}

func (b *bodyView) open(file tcsv.File, taggedStatements []statements.TaggedRow) {
	// create the statements table & view
	b.statementsView = newStatementsView(file, taggedStatements, b.keyMap.StatementsKeyMap)

	// create the summary
	// currently operating in the tea update loop.
	// if statements get too big, may need to move this to a command.
	summary, sortedTags := summarizeStatements(taggedStatements, file.Schema)

	rows := make([]table.Row, len(sortedTags)+1)
	var total float64
	for i, sortedTag := range sortedTags {
		total += summary[sortedTag]
		rows[i] = table.Row{
			sortedTag,
			alignFloat(summary[sortedTag], 10, lipgloss.Right),
		}
	}
	rows[len(sortedTags)] = table.Row{
		tableStyles.Header.Render("Total"),
		tableStyles.Header.Render(alignFloat(total, 10, lipgloss.Right)),
	}

	// create the summary table & view
	b.summaryView = newSummaryView(rows)

	// activate the body
	b.activeTab = 0
}

func (b *bodyView) setSize(width, height int) {
	width -= 2 * frameStyles.Border.GetHorizontalBorderSize()
	height -= frameStyles.Border.GetVerticalBorderSize()

	summaryWidth := width / 5
	statementsWidth := width - summaryWidth

	b.statementsView.SetSize(statementsWidth, height)
	b.summaryView.SetSize(summaryWidth, height)
}

func (b *bodyView) ShortHelp() []key.Binding {
	keys := b.keyMap.ShortHelp()
	switch b.activeTab {
	case 0:
		keys = append(keys, b.statementsView.ShortHelp()...)
	case 1:
		//keys = append(keys, b.summaryView.ShortHelp()...)
	}
	return keys
}

func (b *bodyView) FullHelp() [][]key.Binding {
	return [][]key.Binding{b.ShortHelp()}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var statementColumns = map[string][]table.Column{
	"bnp-debit": {
		{Name: "Date", Width: 10},
		{Name: "Amount", Width: 11, RowStyle: table.CellStyle{Style: lipgloss.NewStyle().Align(lipgloss.Right)}},
		{Name: "Type", Width: 14},
		{Name: "Counterpart", Width: 17},
		{Name: "CounterpartName", Width: 18},
		{Name: "Message", Width: 25},
		{Name: "Details"},
		{Name: "Tag", Width: 10},
	},
	"bnp-visa": {
		{Name: "Date", Width: 10},
		{Name: "Amount", Width: 11, RowStyle: table.CellStyle{Style: lipgloss.NewStyle().Align(lipgloss.Right)}},
		{Name: "Vendor"},
		{Name: "Tag", Width: 10},
	},
}

type statementsView struct {
	*table.Table
	StatementsKeyMap
	schema tcsv.Schema
}

func newStatementsView(f tcsv.File, r []statements.TaggedRow, keyMap StatementsKeyMap) *statementsView {
	// rows
	rows := make([]table.Row, len(r))
	for i, row := range r {
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
	// model
	tbl := table.NewTable(
		"statemenets",
		statementColumns[f.SchemaName],
		rows,
		tableStyles,
		table.DefaultKeyMap(),
	)

	return &statementsView{
		Table:            tbl,
		schema:           f.Schema,
		StatementsKeyMap: keyMap,
	}
}

func (sv *statementsView) Init() tea.Cmd {
	return nil
}

func (sv *statementsView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, sv.Open):
			return func() tea.Msg { return openDetailsMsg{taggedRow: sv.SelectedRow, schema: sv.schema} }
		}
	}
	return sv.Table.Update(msg)
}

func alignString(s string, width int, position lipgloss.Position) string {
	return lipgloss.NewStyle().Align(position).Width(width).Render(s)
}

func alignFloat(f float64, width int, position lipgloss.Position) string {
	return alignString(strconv.FormatFloat(f, 'f', 2, 64), width, position)
}

func loadStatementsCmd(r repo.Repo, tagRules []statements.TagRule, key string) tea.Cmd {
	return func() tea.Msg {
		body, err := r.Get(key)
		if err != nil {
			return errorMsg{err: fmt.Errorf("read: %w", err)}
		}
		f, err := statements.Schemas.Parse(bytes.NewBuffer(body))
		if err != nil {
			return errorMsg{err: fmt.Errorf("parse: %w", err)}
		}
		taggedStatements, err := statements.Tag(f.Rows, f.Schema, tagRules)
		if err != nil {
			return errorMsg{err: fmt.Errorf("tag: %w", err)}
		}

		return loadStatementsMsg{taggedStatements: taggedStatements, file: f}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type summaryView struct {
	*table.Table
}

func newSummaryView(rows []table.Row) *summaryView {
	return &summaryView{
		Table: table.NewTable(
			"summary",
			table.Columns{
				{Name: "Tag"},
				{Name: "Amount", RowStyle: table.CellStyle{Style: lipgloss.NewStyle().Align(lipgloss.Right)}},
			},
			rows,
			tableStyles,
			table.DefaultKeyMap(),
		),
	}
}

func (sv *summaryView) Init() tea.Cmd {
	return nil
}

//func (sv *summaryView) Update(msg tea.Msg) tea.Cmd {
//	return sv.Table.Update(msg)
//}

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
