package tui

import (
	"fmt"

	"codeberg.org/clambin/bubbles/table"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/statements"
)

type repoView struct {
	*table.Table
	repo.Repo
	RepoKeyMap
	tagRules []statements.TagRule
}

func newRepoView(r repo.Repo, tagRules []statements.TagRule, keyMap RepoKeyMap) *repoView {
	return &repoView{
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

func (rv *repoView) Init() tea.Cmd {
	return tea.Batch(
		rv.Table.Init(),
		func() tea.Msg { return statusMsg{text: "Loading files ...", showSpinner: true} },
		rv.loadRepoFilesCmd(),
	)
}

func (rv *repoView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case populateRepoFilesMsg:
		rows := make([]table.Row, len(msg.files))
		for i, f := range msg.files {
			rows[i] = table.Row{f}
		}
		rv.SetRows(rows)
		return func() tea.Msg { return statusMsg{} }
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, rv.Open):
			return tea.Batch(
				func() tea.Msg { return statusMsg{text: "Loading statements file ...", showSpinner: true} },
				rv.openStatementsFileCmd(rv.SelectedRow[0].(string)),
				func() tea.Msg { return setActivePaneMsg{summaryPane} },
			)
		case key.Matches(msg, rv.Reload):
			return tea.Batch(
				func() tea.Msg { return statusMsg{text: "Loading files ...", showSpinner: true} },
				rv.loadRepoFilesCmd(),
			)
		default:
			return rv.Table.Update(msg)
		}
	default:
		return rv.Table.Update(msg)
	}
}

func (rv *repoView) loadRepoFilesCmd() tea.Cmd {
	return func() tea.Msg {
		files, err := rv.List()
		if err != nil {
			return errorMsg(fmt.Errorf("error loading files: %w", err))
		}
		return populateRepoFilesMsg{files: files}
	}
}

func (rv *repoView) openStatementsFileCmd(key string) tea.Cmd {
	return func() tea.Msg {
		taggedStatements, file, err := statements.Open(rv.Repo, key, rv.tagRules)
		if err != nil {
			return errorMsg(fmt.Errorf("open: %w", err))
		}
		return populateStatementsMsg{taggedStatements: taggedStatements, file: file}
	}
}
