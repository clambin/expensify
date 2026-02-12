package statusbar

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Msg struct {
	Text    string
	Warn    bool
	Spinner bool
}

type Style map[bool]lipgloss.Style

var _ tea.Model = Model{}

type Model struct {
	style   Style
	msg     Msg
	spinner spinner.Model
	width   int
}

func New(style Style, opts ...spinner.Option) Model {
	return Model{
		style:   style,
		spinner: spinner.New(opts...),
	}
}

func (s Model) Init() tea.Cmd {
	return s.spinner.Tick
}

func (s Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Msg:
		s.msg = msg
		return s, nil
	default:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd
	}
}

func (s Model) View() string {
	msg := s.msg.Text
	if s.msg.Spinner {
		msg += " " + s.spinner.View()
	}
	return s.style[s.msg.Warn].Width(s.width).MaxHeight(1).Render(msg)
}

func (s Model) Width(width int) Model {
	s.width = width
	return s
}
