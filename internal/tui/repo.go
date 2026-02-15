package tui

import (
	"codeberg.org/clambin/bubbles/table"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/clambin/expensify/bubbles/statusbar"
	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/statements"
)

type repoView struct {
	tea.Model
	repo.Repo
	RepoKeyMap
	tagRules []statements.TagRule
}

func newRepoView(r repo.Repo, tagRules []statements.TagRule, keyMap RepoKeyMap) repoView {
	return repoView{
		Model: table.New().
			Columns([]table.Column{{Name: "Name"}}).
			Styles(tableStyles),
		Repo:       r,
		tagRules:   tagRules,
		RepoKeyMap: keyMap,
	}
}

func (rv repoView) Init() tea.Cmd {
	return tea.Batch(
		rv.Model.Init(),
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
		rv.Model = rv.Model.(table.Table).Rows(rows)
		return rv, func() tea.Msg { return statusbar.Msg{} }
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, rv.Open):
			return rv, tea.Batch(
				func() tea.Msg { return statusbar.Msg{Text: "Loading statements file ...", Spinner: true} },
				openStatementsFileCmd(rv.Repo, rv.tagRules, rv.Model.(table.Table).SelectedRow()[0].(string)),
				func() tea.Msg { return setActivePaneMsg{summaryPane} },
			)
		case key.Matches(msg, rv.Reload):
			return rv, tea.Batch(
				func() tea.Msg { return statusbar.Msg{Text: "Loading files ...", Spinner: true} },
				loadRepoFilesCmd(rv.Repo),
			)
		default:
			var cmd tea.Cmd
			rv.Model, cmd = rv.Model.Update(msg)
			return rv, cmd
		}
	}
	var cmd tea.Cmd
	rv.Model, cmd = rv.Model.Update(msg)
	return rv, cmd
}

func (rv repoView) SetSize(width, height int) tea.Model {
	rv.Model = rv.Model.(table.Table).Size(width, height)
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
