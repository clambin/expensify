package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type statusLine struct {
	msg     statusMsg
	spinner spinner.Model
	width   int
}

func newStatusLine() *statusLine {
	return &statusLine{
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot)),
	}
}

func (s *statusLine) Init() tea.Cmd {
	return func() tea.Msg { return s.spinner.Tick() }
}

func (s *statusLine) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case statusMsg:
		s.msg = msg
	default:
		s.spinner, cmd = s.spinner.Update(msg)
	}
	return cmd
}

func (s *statusLine) View() string {
	msg := s.msg.text
	if s.msg.showSpinner {
		msg += " " + s.spinner.View()
	}
	style := statusStyles.Good
	if s.msg.error {
		style = statusStyles.Error
	}
	return style.Width(s.width).MaxHeight(1).Render(msg)
}

func (s *statusLine) SetSize(width, _ int) {
	s.width = width
}
