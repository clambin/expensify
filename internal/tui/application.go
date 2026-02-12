package tui

import (
	"codeberg.org/clambin/bubbles/frame"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/statements"
)

// paneID identifies the different panes in the application
type paneID int

const (
	repoPane paneID = iota
	summaryPane
	statementsPane
)

// pane is the expected interface for panes in the application
type pane interface {
	Init() tea.Cmd
	Update(msg tea.Msg) tea.Cmd
	View() string
	SetSize(width, height int)
	help.KeyMap
}

var (
	_ tea.Model   = Application{}
	_ help.KeyMap = Application{}
)

type Application struct {
	help       help.Model
	panes      map[paneID]pane
	statusLine *statusLine
	keyMap     ApplicationKeyMap
	activePane paneID
	width      int
	height     int
	fullscreen bool
}

func New(repo repo.Repo, tagRules []statements.TagRule, keyMap KeyMap) tea.Model {
	h := help.New()
	h.Styles = helpStyles

	return Application{
		keyMap: keyMap.ApplicationKeyMap,
		help:   h,
		panes: map[paneID]pane{
			repoPane:       newRepoView(repo, tagRules, keyMap.RepoKeyMap),
			summaryPane:    newSummaryView(keyMap.SummaryKeyMap),
			statementsPane: newStatementsView(keyMap.StatementsListKeyMap, keyMap.StatementsDetailsKeyMap),
		},
		statusLine: newStatusLine(),
		activePane: repoPane,
	}
}

func (a Application) Init() tea.Cmd {
	initCmds := make([]tea.Cmd, 0, len(a.panes))
	for _, c := range a.panes {
		initCmds = append(initCmds, c.Init())
	}
	initCmds = append(initCmds, a.statusLine.Init())
	return tea.Batch(initCmds...)
}

func (a Application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width, a.height = msg.Width, msg.Height
		return a.sizePanes(), nil
	case setActivePaneMsg:
		a.activePane = msg.paneID
		return a.sizePanes(), nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.keyMap.Quit):
			return a, tea.Quit
		case key.Matches(msg, a.keyMap.Next):
			return a, func() tea.Msg { return setActivePaneMsg{paneID: (a.activePane + 1) % 3} }
		case key.Matches(msg, a.keyMap.Previous):
			return a, func() tea.Msg { return setActivePaneMsg{paneID: (a.activePane - 1 + 3) % 3} }
		case key.Matches(msg, a.keyMap.ClearStatus):
			return a, func() tea.Msg { return statusMsg{} }
		case key.Matches(msg, a.keyMap.FullScreenOn):
			a.fullscreen = true
			return a.sizePanes(), nil
		case key.Matches(msg, a.keyMap.FullScreenOff):
			a.fullscreen = false
			return a.sizePanes(), nil
		default:
			return a, a.panes[a.activePane].Update(msg)
		}
	default:
		cmds := make([]tea.Cmd, 0, len(a.panes))
		for _, c := range a.panes {
			cmds = append(cmds, c.Update(msg))
		}
		return a, tea.Batch(append(cmds, a.statusLine.Update(msg))...)
	}
}

func (a Application) View() string {
	var top string
	if a.fullscreen {
		top = a.panes[a.activePane].View()

	} else {
		top = a.viewPaned()
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		top,
		a.statusLine.View(),
		a.help.View(a),
	)
}

func (a Application) ShortHelp() []key.Binding {
	return append(a.keyMap.ShortHelp(), a.panes[a.activePane].ShortHelp()...)
}

func (a Application) FullHelp() [][]key.Binding {
	return [][]key.Binding{a.ShortHelp()}
}

func (a Application) sizePanes() Application {
	borderWidth := frameStyles.Border.GetHorizontalBorderSize()
	borderHeight := frameStyles.Border.GetVerticalBorderSize()

	workingHeight := a.height - 2 // one for status line, one for help

	if a.fullscreen {
		a.panes[a.activePane].SetSize(a.width, workingHeight)
		return a
	}

	headerWidth := a.width / 2
	headerHeight := workingHeight / 3

	a.panes[repoPane].SetSize(headerWidth-borderWidth, headerHeight-borderHeight)
	a.panes[summaryPane].SetSize(a.width-headerWidth-borderWidth, headerHeight-borderHeight)
	a.panes[statementsPane].SetSize(a.width-borderWidth, workingHeight-headerHeight-borderHeight)
	a.statusLine.SetSize(a.width, 1)
	a.help.Width = a.width
	return a
}

func (a Application) viewPaned() string {
	s := map[paneID]frame.Styles{
		repoPane:       frameStyles,
		summaryPane:    frameStyles,
		statementsPane: frameStyles,
	}
	s[a.activePane] = selectedFrameStyles

	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Left,
			frame.Draw("files", lipgloss.Center, a.panes[repoPane].View(), s[repoPane]),
			frame.Draw("summary", lipgloss.Center, a.panes[summaryPane].View(), s[summaryPane]),
		),
		frame.Draw("statements", lipgloss.Center, a.panes[statementsPane].View(), s[statementsPane]),
	)
}
