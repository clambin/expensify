package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
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
			model, err := tui.New(repo, cfg.Rules, tui.DefaultKeyMap())
			if err != nil {
				return fmt.Errorf("failed to create application: %w", err)
			}
			_, err = tea.NewProgram(model, tea.WithAltScreen(), tea.WithoutCatchPanics()).Run()
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
