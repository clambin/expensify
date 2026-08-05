package tui

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"codeberg.org/clambin/bubbles/statusbar"
	"codeberg.org/clambin/bubbles/table"
	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/statements"
)

type repoView struct {
	repo.Repo
	RepoKeyMap
	tagRules []statements.TagRule
	table.Table
	statusBar *statusbar.StatusBar
}

func newRepoView(r repo.Repo, tagRules []statements.TagRule, keyMap RepoKeyMap, statusBar *statusbar.StatusBar) repoView {
	return repoView{
		Table: table.New().
			Columns([]table.Column{{Name: "Name"}}).
			Styles(tableStyles),
		Repo:       r,
		tagRules:   tagRules,
		RepoKeyMap: keyMap,
		statusBar:  statusBar,
	}
}

func (rv repoView) Init() tea.Cmd {
	return tea.Batch(
		rv.Table.Init(),
		rv.statusBar.SetStatusCmd("Loading files ...", statusbar.WithShowSpinner(true)),
		rv.loadRepoFilesCmd(),
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
		return rv, rv.statusBar.SetStatusCmd(fmt.Sprintf("Loaded %d files", len(msg.files)), statusbar.WithTTL(10*time.Second))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, rv.Open):
			return rv, tea.Batch(
				rv.statusBar.SetStatusCmd("Loading statements file ...", statusbar.WithShowSpinner(true)),
				rv.openStatementsFileCmd(rv.SelectedRow()[0].(string)),
				func() tea.Msg { return setActivePaneMsg{summaryPane} },
			)
		case key.Matches(msg, rv.Reload):
			return rv, tea.Batch(
				rv.statusBar.SetStatusCmd("Loading files ...", statusbar.WithShowSpinner(true)),
				rv.loadRepoFilesCmd(),
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

func (rv repoView) loadRepoFilesCmd() tea.Cmd {
	return func() tea.Msg {
		files, err := rv.Repo.List()
		if err != nil {
			return rv.statusBar.SetStatusCmd("Error loading files: "+err.Error(), statusbar.WithLevel("ERROR"))()
		}
		return tea.Batch(
			rv.statusBar.ClearStatusCmd(),
			func() tea.Msg { return populateRepoFilesMsg{files: files} },
		)()
	}
}

func (rv repoView) openStatementsFileCmd(key string) tea.Cmd {
	return func() tea.Msg {
		taggedStatements, file, err := statements.Open(rv.Repo, key, rv.tagRules)
		if err != nil {
			return rv.statusBar.SetStatusCmd("Error loading file "+key+": "+err.Error(), statusbar.WithLevel("ERROR"))()
		}
		return tea.Batch(
			rv.statusBar.ClearStatusCmd(),
			func() tea.Msg { return populateStatementsMsg{taggedStatements: taggedStatements, file: file} },
		)()
	}
}
