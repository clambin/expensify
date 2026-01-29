package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/statements"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(uploadCmd)
}

var (
	uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "upload statements",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, repo, err := loadConfigurationWithRepo(cmd)
			if err != nil {
				return err
			}
			for _, arg := range args {
				key, err := upload(repo, filepath.Clean(arg))
				if err != nil {
					return fmt.Errorf("failed to upload %q: %w", arg, err)
				}
				fmt.Println("Successfully uploaded " + key + ".")
			}
			return nil
		},
	}
)

func upload(repo repo.Repo, filename string) (string, error) {
	r, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %q: %w", filename, err)
	}
	defer func() { _ = r.Close() }()
	key, body, err := analyze(r)
	if err != nil {
		return "", fmt.Errorf("failed to analyze file %q: %w", filename, err)
	}
	// add suffix if key exists
	_, err = repo.Get(key)
	suffix := 1
	for err == nil {
		key = fmt.Sprintf("%s-%d", key, suffix)
		_, err = repo.Get(key)
		suffix++
	}
	err = repo.Add(key, body)
	return key, err
}

func analyze(r io.Reader) (string, []byte, error) {
	var content bytes.Buffer
	f, err := statements.Schemas.Parse(io.TeeReader(r, &content))
	if err != nil {
		return "", nil, fmt.Errorf("parse: %w", err)
	}
	if len(f.Rows) == 0 {
		return "", nil, fmt.Errorf("no statements found")
	}
	minDate := f.ToMap(f.Rows[0])["Date"].(time.Time)
	for _, stmt := range f.Rows[1:] {
		if stmtDate := f.ToMap(stmt)["Date"].(time.Time); stmtDate.Before(minDate) {
			minDate = stmtDate
		}
	}
	return f.SchemaName + "-" + minDate.Format("2006-01-02"), content.Bytes(), err
}
