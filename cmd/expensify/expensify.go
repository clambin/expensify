package main

import (
	"fmt"
	"github.com/clambin/expensify/internal/analyzer"
	"github.com/clambin/expensify/internal/payment"
	"github.com/clambin/expensify/pkg/maps"
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
	rules, err := loadRules()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, filename := range args {
		matched, _, err := analyze(filename, rules)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, category := range maps.SortedKeys(matched) {
			if category == "ignored" && showIgnored == false {
				continue
			}
			fmt.Printf("%s: %.2f\n", category, matched[category].Total())
			if summaryDetails {
				for i := range matched[category] {
					fmt.Printf("\t%.2f\t%s\n", matched[category][i].GetAmount(), matched[category][i].GetDescription())
				}
			}
		}
	}
}

func showUnmatched(_ *cobra.Command, args []string) {
	rules, err := loadRules()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, filename := range args {
		_, unmatched, err := analyze(filename, rules)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for i := range unmatched {
			fmt.Printf("\t%.2f\t%s\n", unmatched[i].GetAmount(), unmatched[i].GetDescription())
		}
	}
}

func loadRules() (analyzer.Rules, error) {
	f, err := os.Open(rulesFilename)
	if err != nil {
		return nil, fmt.Errorf("rules: %w", err)
	}
	defer func() { _ = f.Close() }()
	return analyzer.LoadRules(f)
}

func analyze(filename string, rules analyzer.Rules) (map[string]payment.Payments, payment.Payments, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = f.Close() }()
	return analyzer.Analyze(f, rules)
}
