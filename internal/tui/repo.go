package tui

import (
	"codeberg.org/clambin/bubbles/frame"
	"codeberg.org/clambin/bubbles/table"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/clambin/expensify/internal/repo"
)

type repoView struct {
	*table.Table
	repo.Repo
	RepoKeyMap
}

func newRepoView(rows []table.Row, repo repo.Repo, keyMap RepoKeyMap) *repoView {
	return &repoView{
		Table: table.NewTable(
			"files",
			[]table.Column{{Name: "Name"}},
			rows,
			tableStyles,
			table.DefaultKeyMap(),
		),
		Repo:       repo,
		RepoKeyMap: keyMap,
	}
}

func (rv *repoView) Init() tea.Cmd {
	return nil
}

func (rv *repoView) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, rv.Open):
			return func() tea.Msg { return openStatementsMsg{name: rv.SelectedRow[0].(string)} }
		default:
			cmd = rv.Table.Update(msg)
			return cmd
		}
	default:
		cmd = rv.Table.Update(msg)
	}
	return cmd
}

func (rv *repoView) View() string {
	return frame.Draw("files", lipgloss.Center, rv.Table.View(), frameStyles)
}

func (rv *repoView) SetSize(width, height int) {
	width -= frameStyles.Border.GetHorizontalBorderSize()
	height -= frameStyles.Border.GetVerticalBorderSize()
	rv.Table.SetSize(width, height)
}
