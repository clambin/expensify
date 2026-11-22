package cmd

import (
	"context"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"codeberg.org/clambin/bubbles/msglogger"
	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/tui"
)

// TUICmd starts the interactive UI.
type TUICmd struct{}

func (c *TUICmd) Run(ctx context.Context, cli *CLI) error {
	r, err := repo.New(ctx, cli.Storage)
	if err != nil {
		return err
	}

	rules, err := cli.rules()
	if err != nil {
		return fmt.Errorf("failed to read rules file: %w", err)
	}

	var logger *msglogger.MsgLogger
	/*
		f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		logger = msglogger.New(slog.NewTextHandler(f, nil),
			spinner.TickMsg{},
			// tea.KeyPressMsg{},
		)
	*/
	model := tui.New(r, rules, tui.DefaultKeyMap(), logger)
	p := tea.NewProgram(model, tea.WithContext(ctx)) //tea.WithoutCatchPanics(),

	_, err = p.Run()
	return err
}
