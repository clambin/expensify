package tui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"codeberg.org/clambin/bubbles/table"
	"github.com/clambin/expensify/bubbles/statusbar"
	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/statements"
)

type repoView struct {
	repo.Repo
	RepoKeyMap
	tagRules []statements.TagRule
	table.Table
}

func newRepoView(r repo.Repo, tagRules []statements.TagRule, keyMap RepoKeyMap) repoView {
	return repoView{
		Table: table.New().
			Columns([]table.Column{{Name: "Name"}}).
			Styles(tableStyles),
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

func (rv repoView) Update(msg tea.Msg) (repoView, tea.Cmd) {
	switch msg := msg.(type) {
	case populateRepoFilesMsg:
		rows := make([]table.Row, len(msg.files))
		for i, f := range msg.files {
			rows[i] = table.Row{f}
		}
		rv.Table = rv.Rows(rows)
		return rv, func() tea.Msg { return statusbar.Msg{} }
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, rv.Open):
			return rv, tea.Batch(
				func() tea.Msg { return statusbar.Msg{Text: "Loading statements file ...", Spinner: true} },
				openStatementsFileCmd(rv.Repo, rv.tagRules, rv.SelectedRow()[0].(string)),
				func() tea.Msg { return setActivePaneMsg{summaryPane} },
			)
		case key.Matches(msg, rv.Reload):
			return rv, tea.Batch(
				func() tea.Msg { return statusbar.Msg{Text: "Loading files ...", Spinner: true} },
				loadRepoFilesCmd(rv.Repo),
			)
		}
	}
	var cmd tea.Cmd
	rv.Table, cmd = rv.Table.Update(msg)
	return rv, cmd
}

func (rv repoView) SetSize(width, height int) repoView {
	rv.Table = rv.Size(width, height)
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
