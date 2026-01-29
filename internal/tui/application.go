package tui

import (
	"fmt"
	"os"

	"codeberg.org/clambin/bubbles/table"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/statements"
)

type activePane int

const (
	noPane activePane = iota
	statementsPane
	detailsPane
)

var _ tea.Model = Application{}

type Application struct {
	keyMap      ApplicationKeyMap
	help        help.Model
	repo        repo.Repo
	repoView    *repoView
	bodyView    *bodyView
	detailsView *detailsView
	tagRules    []statements.TagRule
	activePane  activePane
	width       int
	height      int
}

func New(repo repo.Repo, tagRules []statements.TagRule, keyMap KeyMap) (tea.Model, error) {
	files, err := repo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	rows := make([]table.Row, len(files))
	for i, file := range files {
		rows[i] = table.Row{file}
	}

	h := help.New()
	h.Styles = helpStyles

	return Application{
		keyMap:      keyMap.ApplicationKeyMap,
		help:        h,
		repo:        repo,
		repoView:    newRepoView(rows, repo, keyMap.RepoKeyMap),
		bodyView:    newBodyView(keyMap.BodyKeyMap),
		detailsView: newDetailsView(keyMap.DetailsKeyMap),
		tagRules:    tagRules,
		activePane:  noPane,
	}, nil
}

func (a Application) Init() tea.Cmd {
	return tea.Batch(
		a.repoView.Init(),
		a.bodyView.Init(),
	)
}

func (a Application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width, a.height = msg.Width, msg.Height
		a.resize()
	case errorMsg:
		_, _ = fmt.Fprintln(os.Stderr, msg.err)
	case openStatementsMsg:
		cmd = loadStatementsCmd(a.repo, a.tagRules, msg.name)
	case loadStatementsMsg:
		a.bodyView.open(msg.file, msg.taggedStatements)
		a.activePane = statementsPane
		a.resize()
	case closeStatementsMsg:
		a.activePane = noPane
		a.resize()
	case openDetailsMsg:
		a.detailsView.load(msg.taggedRow, msg.schema)
		a.activePane = detailsPane
	case closeDetailsMsg:
		a.activePane = statementsPane
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.keyMap.Quit):
			return a, tea.Quit
		default:
			switch a.activePane {
			case noPane:
				cmd = a.repoView.Update(msg)
			case statementsPane:
				cmd = a.bodyView.Update(msg)
			case detailsPane:
				cmd = a.detailsView.Update(msg)
			}
		}
	default:
		switch a.activePane {
		case noPane:
			cmd = a.repoView.Update(msg)
		case statementsPane:
			cmd = a.bodyView.Update(msg)
		case detailsPane:
			cmd = a.detailsView.Update(msg)
		}
	}
	return a, cmd
}

func (a Application) View() string {
	components := make([]string, 0, 3)
	components = append(components, a.repoView.View())
	switch a.activePane {
	case statementsPane:
		components = append(components, a.bodyView.View())
	case detailsPane:
		components = append(components, a.detailsView.View())
	default:
	}
	components = append(components, a.viewShortHelp())
	return lipgloss.JoinVertical(lipgloss.Top, components...)
}

func (a Application) resize() {
	switch a.activePane {
	case noPane:
		a.repoView.SetSize(a.width, a.height-1)
	default:
		a.repoView.SetSize(a.width, 5)
		a.bodyView.setSize(a.width, a.height-5-1)
		a.detailsView.SetSize(a.width, a.height-5-1)
	}
}

func (a Application) viewShortHelp() string {
	keys := a.keyMap.ShortHelp()
	switch a.activePane {
	case noPane:
		keys = append(keys, a.repoView.ShortHelp()...)
	case statementsPane:
		keys = append(keys, a.bodyView.ShortHelp()...)
	case detailsPane:
		keys = append(keys, a.detailsView.ShortHelp()...)
	}
	return a.help.ShortHelpView(keys)
}
