package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var (
	_ help.KeyMap = ApplicationKeyMap{}
	_ help.KeyMap = RepoKeyMap{}
	_ help.KeyMap = SummaryKeyMap{}
	_ help.KeyMap = StatementsListKeyMap{}
	_ help.KeyMap = StatementsDetailsKeyMap{}
)

type KeyMap struct {
	ApplicationKeyMap
	RepoKeyMap
	SummaryKeyMap
	StatementsListKeyMap
	StatementsDetailsKeyMap
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		ApplicationKeyMap:       defaultApplicationKeyMap(),
		RepoKeyMap:              defaultRepoKeyMap(),
		SummaryKeyMap:           defaultSummaryKeyMap(),
		StatementsListKeyMap:    defaultStatementsListKeyMap(),
		StatementsDetailsKeyMap: defaultStatementsDetailsKeyMap(),
	}
}

type ApplicationKeyMap struct {
	Quit             key.Binding
	Next             key.Binding
	Previous         key.Binding
	ClearStatus      key.Binding
	ToggleFullscreen key.Binding
}

func defaultApplicationKeyMap() ApplicationKeyMap {
	return ApplicationKeyMap{
		Quit:             key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "Quit the application")),
		Next:             key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "Select next file")),
		Previous:         key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "Select previous file")),
		ClearStatus:      key.NewBinding(key.WithKeys("alt+c"), key.WithHelp("alt+c", "Clear the status bar")),
		ToggleFullscreen: key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "Enter full screen mode")),
	}
}

func (m ApplicationKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{m.Quit, m.ClearStatus, m.ToggleFullscreen}
}

func (m ApplicationKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}

type RepoKeyMap struct {
	Open   key.Binding
	Reload key.Binding
}

func defaultRepoKeyMap() RepoKeyMap {
	return RepoKeyMap{
		Open:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Open the selected file")),
		Reload: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "Reload the files")),
	}
}

func (m RepoKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{m.Open, m.Reload}
}

func (m RepoKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}

var _ help.KeyMap = SummaryKeyMap{}

type SummaryKeyMap struct {
	Open key.Binding
}

func defaultSummaryKeyMap() SummaryKeyMap {
	return SummaryKeyMap{
		Open: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Open statements for the selected tag")),
	}
}

func (m SummaryKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{m.Open}
}

func (m SummaryKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}

type StatementsListKeyMap struct {
	Details key.Binding
}

func defaultStatementsListKeyMap() StatementsListKeyMap {
	return StatementsListKeyMap{
		Details: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Open the selected statement")),
	}
}

func (m StatementsListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{m.Details}
}

func (m StatementsListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}

type StatementsDetailsKeyMap struct {
	Close key.Binding
}

func defaultStatementsDetailsKeyMap() StatementsDetailsKeyMap {
	return StatementsDetailsKeyMap{
		Close: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Close the details")),
	}
}

func (m StatementsDetailsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{m.Close}
}

func (m StatementsDetailsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}
