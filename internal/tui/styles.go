package tui

import (
	"codeberg.org/clambin/bubbles/colors"
	"codeberg.org/clambin/bubbles/frame"
	"codeberg.org/clambin/bubbles/table"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

var (
	tableStyles = table.Styles{
		Header:   lipgloss.NewStyle().Foreground(colors.White).Bold(true),
		Selected: lipgloss.NewStyle().Foreground(colors.Orchid2),
		Cell:     lipgloss.NewStyle().Foreground(colors.Grey74),
	}

	helpStyles = help.Styles{
		ShortKey:       lipgloss.NewStyle().Foreground(colors.Orange3),
		ShortDesc:      lipgloss.NewStyle().Foreground(colors.DarkOrange3),
		ShortSeparator: lipgloss.NewStyle().Foreground(colors.Orange3),
	}

	selectedFrameStyles = frame.Styles{
		Title:  lipgloss.NewStyle().Foreground(colors.Grey74).Bold(true),
		Border: lipgloss.NewStyle().Border(lipgloss.DoubleBorder()),
	}

	frameStyles = frame.Styles{
		Title:  lipgloss.NewStyle().Foreground(colors.Grey74).Bold(false),
		Border: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()),
	}

	statusStyles = struct {
		Good  lipgloss.Style
		Error lipgloss.Style
	}{
		Good:  lipgloss.NewStyle().Foreground(colors.White).Background(colors.Aqua),
		Error: lipgloss.NewStyle().Foreground(colors.White).Background(colors.DarkRed),
	}
)
