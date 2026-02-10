package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list statement files",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, repo, err := loadConfigurationWithRepo(cmd)
		if err != nil {
			return err
		}
		files, err := repo.List()
		if err != nil {
			return fmt.Errorf("list: %w", err)
		}
		for _, file := range files {
			fmt.Println(file)
		}
		return nil
	},
}
