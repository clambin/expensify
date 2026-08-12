package cmd

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"codeberg.org/clambin/bubbles/msglogger"
	"github.com/clambin/expensify/internal/configuration"
	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/tui"
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use: "expensify",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, repo, err := loadConfigurationWithRepo(cmd)
			if err != nil {
				return err
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
			model := tui.New(repo, cfg.Rules, tui.DefaultKeyMap(), logger)
			p := tea.NewProgram(model) //tea.WithoutCatchPanics(),

			_, err = p.Run()
			return err
		},
	}
)

func init() {
	RootCmd.PersistentFlags().StringP("config", "c", "/etc/config.yaml", "configuration file")
}

func loadConfiguration(cmd *cobra.Command) (configuration.Configuration, error) {
	filename, err := cmd.Flags().GetString("config")
	if err != nil {
		return configuration.Configuration{}, err
	}
	return configuration.LoadFileConfiguration(filename)
}

func loadConfigurationWithRepo(cmd *cobra.Command) (cfg configuration.Configuration, r repo.Repo, err error) {
	if cfg, err = loadConfiguration(cmd); err != nil {
		return configuration.Configuration{}, nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	if r, err = cfg.Backend.Backend(cmd.Context()); err != nil {
		return configuration.Configuration{}, nil, fmt.Errorf("failed to load repository backend: %w", err)
	}
	return cfg, r, err
}
