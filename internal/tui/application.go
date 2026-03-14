package tui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"codeberg.org/clambin/bubbles/frame"
	"github.com/clambin/expensify/bubbles/statusbar"
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

var (
	_ tea.Model   = Application{}
	_ help.KeyMap = Application{}
)

type Application struct {
	help           help.Model
	keyMap         ApplicationKeyMap
	repoPane       repoView
	summaryPane    summaryView
	statusLine     statusbar.Model
	statementsPane statementsView
	activePane     paneID
	width          int
	height         int
	fullscreen     bool
}

func New(repo repo.Repo, tagRules []statements.TagRule, keyMap KeyMap) tea.Model {
	h := help.New()
	h.Styles = helpStyles

	return Application{
		keyMap:         keyMap.ApplicationKeyMap,
		help:           h,
		repoPane:       newRepoView(repo, tagRules, keyMap.RepoKeyMap),
		summaryPane:    newSummaryView(keyMap.SummaryKeyMap),
		statementsPane: newStatementsView(keyMap.StatementsListKeyMap, keyMap.StatementsDetailsKeyMap),
		statusLine:     statusbar.New(statusStyles, spinner.WithSpinner(spinner.Dot)),
		activePane:     repoPane,
	}
}

func (a Application) Init() tea.Cmd {
	return tea.Batch(
		a.repoPane.Init(),
		a.summaryPane.Init(),
		a.statementsPane.Init(),
		a.statusLine.Init(),
	)
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
			switch a.activePane {
			case repoPane:
				a.repoPane, cmd = a.repoPane.Update(msg)
			case summaryPane:
				a.summaryPane, cmd = a.summaryPane.Update(msg)
			case statementsPane:
				a.statementsPane, cmd = a.statementsPane.Update(msg)
			}
			return a, cmd
		}
	default:
		cmds := make([]tea.Cmd, 3)
		a.repoPane, cmds[0] = a.repoPane.Update(msg)
		a.summaryPane, cmds[1] = a.summaryPane.Update(msg)
		a.statementsPane, cmds[2] = a.statementsPane.Update(msg)
		return a, tea.Batch(cmds...)
	}
}

func (a Application) View() tea.View {
	var top string
	switch a.fullscreen {
	case false:
		top = a.viewPaned()
	case true:
		switch a.activePane {
		case repoPane:
			top = a.repoPane.View()
		case summaryPane:
			top = a.summaryPane.View()
		case statementsPane:
			top = a.statementsPane.View()
		}
	}
	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Top,
		top,
		a.statusLine.Width(a.width).View(),
		a.help.View(a),
	))
	v.AltScreen = true
	return v
}

func (a Application) ShortHelp() []key.Binding {
	var b []key.Binding
	switch a.activePane {
	case repoPane:
		b = a.repoPane.ShortHelp()
	case summaryPane:
		b = a.summaryPane.ShortHelp()
	case statementsPane:
		b = a.statementsPane.ShortHelp()
	}
	return append(a.keyMap.ShortHelp(), b...)
}

func (a Application) FullHelp() [][]key.Binding {
	return [][]key.Binding{a.ShortHelp()}
}

func (a Application) sizePanes() Application {
	workingHeight := a.height - 2 // one for status line, one for help

	if a.fullscreen {
		switch a.activePane {
		case repoPane:
			a.repoPane = a.repoPane.SetSize(a.width, workingHeight)
		case summaryPane:
			a.summaryPane = a.summaryPane.SetSize(a.width, workingHeight)
		case statementsPane:
			a.statementsPane = a.statementsPane.SetSize(a.width, workingHeight)
		}
		return a
	}

	borderWidth := frameStyles.Border.GetHorizontalBorderSize()
	borderHeight := frameStyles.Border.GetVerticalBorderSize()

	headerWidth := a.width / 2
	headerHeight := workingHeight / 3

	a.repoPane = a.repoPane.SetSize(headerWidth-borderWidth, headerHeight-borderHeight)
	a.summaryPane = a.summaryPane.SetSize(a.width-headerWidth-borderWidth, headerHeight-borderHeight)
	a.statementsPane = a.statementsPane.SetSize(a.width-borderWidth, workingHeight-headerHeight-borderHeight)
	a.help.SetWidth(a.width)
	return a
}

func (a Application) viewPaned() string {
	styles := map[paneID]frame.Style{
		repoPane:       frameStyles,
		summaryPane:    frameStyles,
		statementsPane: frameStyles,
	}
	styles[a.activePane] = selectedFrameStyles

	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Left,
			frame.Render("files", lipgloss.Center, styles[repoPane], a.repoPane.View()),
			frame.Render("summary", lipgloss.Center, styles[summaryPane], a.summaryPane.View()),
		),
		frame.Render("statements", lipgloss.Center, styles[statementsPane], a.statementsPane.View()),
	)
}
