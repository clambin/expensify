package tui

import (
	"codeberg.org/clambin/bubbles/table"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/clambin/expensify/bubbles/statusbar"
	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/statements"
)

var _ pane = repoView{}

type repoView struct {
	*table.Table
	repo.Repo
	RepoKeyMap
	tagRules []statements.TagRule
}

func newRepoView(r repo.Repo, tagRules []statements.TagRule, keyMap RepoKeyMap) tea.Model {
	return repoView{
		Table: table.NewTable(
			"files",
			[]table.Column{{Name: "Name"}},
			nil,
			tableStyles,
			table.DefaultKeyMap(),
		),
		Repo:       r,
		tagRules:   tagRules,
		RepoKeyMap: keyMap,
	}
}

func (rv repoView) Init() tea.Cmd {
	return tea.Batch(
		rv.Table.Init(),
		func() tea.Msg { return statusbar.Msg{Text: "Loading files ...", Spinner: true} },
		loadRepoFilesCmd(rv.Repo),
	)
}

func (rv repoView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case populateRepoFilesMsg:
		rows := make([]table.Row, len(msg.files))
		for i, f := range msg.files {
			rows[i] = table.Row{f}
		}
		rv.SetRows(rows)
		return rv, func() tea.Msg { return statusbar.Msg{} }
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, rv.Open):
			return rv, tea.Batch(
				func() tea.Msg { return statusbar.Msg{Text: "Loading statements file ...", Spinner: true} },
				openStatementsFileCmd(rv.Repo, rv.tagRules, rv.SelectedRow[0].(string)),
				func() tea.Msg { return setActivePaneMsg{summaryPane} },
			)
		case key.Matches(msg, rv.Reload):
			return rv, tea.Batch(
				func() tea.Msg { return statusbar.Msg{Text: "Loading files ...", Spinner: true} },
				loadRepoFilesCmd(rv.Repo),
			)
		default:
			return rv, rv.Table.Update(msg)
		}
	default:
		return rv, rv.Table.Update(msg)
	}
}

func (rv repoView) SetSize(width, height int) tea.Model {
	rv.Table.SetSize(width, height)
	return rv
}

func loadRepoFilesCmd(r repo.Repo) tea.Cmd {
	return func() tea.Msg {
		files, err := r.List()
		if err != nil {
			return statusbar.Msg{Text: "Error loading files: " + err.Error(), Warn: true}
		}
		return populateRepoFilesMsg{files: files}
	}
}

func openStatementsFileCmd(r repo.Repo, tagRules statements.TagRules, key string) tea.Cmd {
	return func() tea.Msg {
		taggedStatements, file, err := statements.Open(r, key, tagRules)
		if err != nil {
			return statusbar.Msg{Text: "Error loading file " + key + ": " + err.Error(), Warn: true}
		}
		return populateStatementsMsg{taggedStatements: taggedStatements, file: file}
	}
}
