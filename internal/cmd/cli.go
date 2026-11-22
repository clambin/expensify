package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/statements"
	"github.com/goccy/go-yaml"
)

// CLI is the root of the command line grammar.
type CLI struct {
	Config  ConfigFile         `name:"config" short:"c" type:"path" default:"~/.expensify/config.yaml" help:"Path to configuration file"`
	Storage repo.Configuration `embed:"" group:"Storage" prefix:"storage-"`
	Rules   string             `default:"rules.yaml" help:"path to rules file. Relative paths are relative to the configuration file's directory"`

	TUI    TUICmd    `cmd:"" name:"tui" default:"1" hidden:"" help:"start the interactive UI (default)"`
	List   ListCmd   `cmd:"" help:"list statement files"`
	Upload UploadCmd `cmd:"" help:"upload statement files"`
}

// rules loads the tag rules from the rules file. A relative path is relative to the configuration file's directory.
//
// Rules deliberately isn't a type:"path" flag: kong would make relative paths absolute to the current directory.
func (c CLI) rules() (statements.TagRules, error) {
	rulesPath := c.Rules
	switch {
	case strings.HasPrefix(rulesPath, "~/"):
		rulesPath = kong.ExpandPath(rulesPath)
	case !filepath.IsAbs(rulesPath):
		rulesPath = filepath.Join(filepath.Dir(kong.ExpandPath(string(c.Config))), rulesPath)
	}
	body, err := os.ReadFile(rulesPath)
	if err != nil {
		return nil, fmt.Errorf("rules: %w", err)
	}
	var rules statements.TagRules
	if err = yaml.Unmarshal(body, &rules); err != nil {
		return nil, fmt.Errorf("rules: %w", err)
	}
	return rules, nil
}
