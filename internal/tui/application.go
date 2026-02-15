package tui

import (
	"codeberg.org/clambin/bubbles/frame"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/clambin/expensify/bubbles/statusbar"
	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/statements"
)

type pane interface {
	tea.Model
	help.KeyMap
	SetSize(int, int) tea.Model
}

// paneID identifies the different panes in the application
type paneID int

const (
	repoPane paneID = iota
	summaryPane
	statementsPane
)

var (
	_ tea.Model   = Application{}
	_ help.KeyMap = Application{}
)

type Application struct {
	help       help.Model
	panes      map[paneID]tea.Model
	statusLine tea.Model
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
		panes: map[paneID]tea.Model{
			repoPane:       newRepoView(repo, tagRules, keyMap.RepoKeyMap),
			summaryPane:    newSummaryView(keyMap.SummaryKeyMap),
			statementsPane: newStatementsView(keyMap.StatementsListKeyMap, keyMap.StatementsDetailsKeyMap),
		},
		statusLine: statusbar.New(statusStyles, spinner.WithSpinner(spinner.Dot)),
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
	case statusbar.Msg:
		var cmd tea.Cmd
		a.statusLine, cmd = a.statusLine.Update(msg)
		return a, cmd
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
			return a, func() tea.Msg { return statusbar.Msg{} }
		case key.Matches(msg, a.keyMap.ToggleFullscreen):
			a.fullscreen = !a.fullscreen
			return a.sizePanes(), nil
		default:
			var cmd tea.Cmd
			a.panes[a.activePane], cmd = a.panes[a.activePane].Update(msg)
			return a, cmd
		}
	default:
		cmds := make([]tea.Cmd, 0, len(a.panes))
		for pid := range a.panes {
			var cmd tea.Cmd
			a.panes[pid], cmd = a.panes[pid].Update(msg)
			cmds = append(cmds, cmd)
		}
		return a, tea.Batch(cmds...)
	}
}

func (a Application) View() string {
	var top string
	switch a.fullscreen {
	case true:
		top = a.panes[a.activePane].View()
	case false:
		top = a.viewPaned()
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		top,
		a.statusLine.(statusbar.Model).Width(a.width).View(),
		a.help.View(a),
	)
}

func (a Application) ShortHelp() []key.Binding {
	return append(a.keyMap.ShortHelp(), a.panes[a.activePane].(help.KeyMap).ShortHelp()...)
}

func (a Application) FullHelp() [][]key.Binding {
	return [][]key.Binding{a.ShortHelp()}
}

func (a Application) sizePanes() Application {
	workingHeight := a.height - 2 // one for status line, one for help

	if a.fullscreen {
		a.panes[a.activePane].(pane).SetSize(a.width, workingHeight)
		return a
	}

	borderWidth := frameStyles.Border.GetHorizontalBorderSize()
	borderHeight := frameStyles.Border.GetVerticalBorderSize()

	headerWidth := a.width / 2
	headerHeight := workingHeight / 3

	a.panes[repoPane] = a.panes[repoPane].(pane).SetSize(headerWidth-borderWidth, headerHeight-borderHeight)
	a.panes[summaryPane] = a.panes[summaryPane].(pane).SetSize(a.width-headerWidth-borderWidth, headerHeight-borderHeight)
	a.panes[statementsPane] = a.panes[statementsPane].(pane).SetSize(a.width-borderWidth, workingHeight-headerHeight-borderHeight)
	a.statusLine.(statusbar.Model).Width(a.width).View()
	a.help.Width = a.width
	return a
}

func (a Application) viewPaned() string {
	s := map[paneID]frame.Style{
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
