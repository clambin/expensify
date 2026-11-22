package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/statements"
)

// UploadCmd uploads statement files to the repository.
type UploadCmd struct {
	Files []string `arg:"" type:"existingfile" help:"statement files to upload"`
}

func (c *UploadCmd) Run(ctx context.Context, cli *CLI) error {
	r, err := repo.New(ctx, cli.Storage)
	if err != nil {
		return err
	}
	for _, file := range c.Files {
		key, err := upload(r, filepath.Clean(file))
		if err != nil {
			return fmt.Errorf("failed to upload %q: %w", file, err)
		}
		fmt.Println("Successfully uploaded " + key + ".")
	}
	return nil
}

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
