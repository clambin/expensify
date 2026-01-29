package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var (
	_ help.KeyMap = ApplicationKeyMap{}
	_ help.KeyMap = BodyKeyMap{}
	_ help.KeyMap = StatementsKeyMap{}
	_ help.KeyMap = DetailsKeyMap{}
)

type KeyMap struct {
	ApplicationKeyMap
	RepoKeyMap
	BodyKeyMap
	DetailsKeyMap
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		ApplicationKeyMap: defaultApplicationKeyMap(),
		RepoKeyMap:        defaultRepoKeyMap(),
		BodyKeyMap:        defaultBodyViewKeyMap(),
		DetailsKeyMap:     defaultDetailsKeyMap(),
	}
}

type ApplicationKeyMap struct {
	Quit key.Binding
	//Reload key.Binding
}

func defaultApplicationKeyMap() ApplicationKeyMap {
	return ApplicationKeyMap{
		Quit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "Quit the application")),
		//Reload: key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "Reload the statements")),
	}
}

func (a ApplicationKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{a.Quit}
}

func (a ApplicationKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{a.ShortHelp()}
}

type RepoKeyMap struct {
	Open key.Binding
}

func defaultRepoKeyMap() RepoKeyMap {
	return RepoKeyMap{
		Open: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Open the selected file")),
	}
}

func (r RepoKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{r.Open}
}

func (r RepoKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{r.ShortHelp()}
}

type BodyKeyMap struct {
	NextPane key.Binding
	PrevPane key.Binding
	Close    key.Binding
	StatementsKeyMap
}

func defaultBodyViewKeyMap() BodyKeyMap {
	return BodyKeyMap{
		NextPane:         key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "Switch to next pane")),
		PrevPane:         key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "Switch to previous pane")),
		Close:            key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Close the statements")),
		StatementsKeyMap: defaultStatementsKeyMap(),
	}
}

func (b BodyKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		b.Close,
		b.NextPane,
		b.PrevPane,
	}
}

func (b BodyKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{b.ShortHelp()}
}

type StatementsKeyMap struct {
	Open key.Binding
}

func defaultStatementsKeyMap() StatementsKeyMap {
	return StatementsKeyMap{
		Open: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Open the selected statement")),
	}
}

func (s StatementsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{s.Open}
}

func (s StatementsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{s.ShortHelp()}
}

type DetailsKeyMap struct {
	Close key.Binding
}

func defaultDetailsKeyMap() DetailsKeyMap {
	return DetailsKeyMap{
		Close: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Close the details")),
	}
}

func (d DetailsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{d.Close}
}

func (d DetailsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{d.ShortHelp()}
}
