package main

import (
	"fmt"
	"github.com/clambin/expensify/categorizer"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	cmd            *cobra.Command
	rulesFilename  string
	summaryDetails bool
	showIgnored    bool
)

func init() {
	cmd = &cobra.Command{
		Use:   "expensify",
		Short: "summarize monthly expenses",
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	rulesPath := filepath.Join(homeDir, ".expensify", "rules.yaml")
	cmd.PersistentFlags().StringVarP(&rulesFilename, "rules", "r", rulesPath, "rules file")

	summaryCmd := &cobra.Command{
		Use:   "summary",
		Run:   showSummary,
		Short: "report summary of expenses",
	}
	summaryCmd.Flags().BoolVarP(&summaryDetails, "detail", "d", false, "print summary details")
	summaryCmd.Flags().BoolVarP(&showIgnored, "ignored", "i", false, "show ignored category")
	cmd.AddCommand(summaryCmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "unmatched",
		Run:   showUnmatched,
		Short: "show all payments that did not match any rules",
	})
}

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}

func showSummary(_ *cobra.Command, args []string) {
	c, err := categorizer.New(rulesFilename, args...)
	if err != nil {
		fmt.Println(err)
		return
	}

	for category, entries := range c.Matched() {
		if category == "ignore" && !showIgnored {
			continue
		}
		fmt.Printf("%s: %.2f\n", category, entries.Total())
		if summaryDetails {
			for _, entry := range entries {
				fmt.Printf("\t%.2f\t%s\n", entry.GetAmount(), entry.GetDescription())
			}
		}
	}
}

func showUnmatched(_ *cobra.Command, args []string) {
	c, err := categorizer.New(rulesFilename, args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, unmatched := range c.Unmatched() {
		fmt.Printf("\t%.2f\t%s\n", unmatched.GetAmount(), unmatched.GetDescription())
	}
}
