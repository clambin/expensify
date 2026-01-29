package tui

import (
	"strings"

	"codeberg.org/clambin/bubbles/frame"
	"codeberg.org/clambin/bubbles/table"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/clambin/expensify/tcsv"
)

type detailsView struct {
	frame.Styles
	*table.Table
	DetailsKeyMap
}

func newDetailsView(keyMap DetailsKeyMap) *detailsView {
	return &detailsView{
		Table: table.NewTable("statement details", table.Columns{
			{Name: "Field", Width: 15},
			{Name: "Value"},
		}, nil, tableStyles, table.DefaultKeyMap()),
		Styles:        frameStyles,
		DetailsKeyMap: keyMap,
	}
}

func (d *detailsView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.Close):
			return func() tea.Msg { return closeDetailsMsg{} }
		}
	}
	return d.Table.Update(msg)
}

func (d *detailsView) View() string {
	return frame.Draw("statemenet details", lipgloss.Center, d.Table.View(), frameStyles)

}

func (d *detailsView) SetSize(width, height int) {
	width -= frameStyles.Border.GetHorizontalBorderSize()
	height -= frameStyles.Border.GetVerticalBorderSize()
	d.Table.SetSize(width, height)
}

func (d *detailsView) load(taggedRow table.Row, schema tcsv.Schema) tea.Cmd {
	rows := make([]table.Row, 0, len(taggedRow)-1)
	var idx int
	for _, col := range schema.Columns {
		if col.Label != "" {
			rows = append(rows, table.Row{col.Label, strings.TrimSpace(ansi.Strip(taggedRow[idx].(string)))})
			idx++
		}
	}
	rows = append(rows, table.Row{"Tag", taggedRow[len(taggedRow)-1]})
	d.SetRows(rows)
	return nil
}
