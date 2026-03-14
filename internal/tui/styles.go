package tui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/lipgloss/v2"
	"codeberg.org/clambin/bubbles/colors"
	"codeberg.org/clambin/bubbles/frame"
	"codeberg.org/clambin/bubbles/table"
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

	selectedFrameStyles = frame.Style{
		Title:  lipgloss.NewStyle().Foreground(colors.Grey74).Bold(true),
		Border: lipgloss.NewStyle().Border(lipgloss.DoubleBorder()),
	}

	frameStyles = frame.Style{
		Title:  lipgloss.NewStyle().Foreground(colors.Grey74).Bold(false),
		Border: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()),
	}

	statusStyles = map[bool]lipgloss.Style{
		false: lipgloss.NewStyle().Foreground(colors.White).Background(colors.Aqua),
		true:  lipgloss.NewStyle().Foreground(colors.White).Background(colors.DarkRed),
	}
)
