package cmd

import (
	"context"
	"fmt"

	"github.com/clambin/expensify/internal/repo"
)

// ListCmd lists the statement files in the repository.
type ListCmd struct{}

func (c *ListCmd) Run(ctx context.Context, cli *CLI) error {
	r, err := repo.New(ctx, cli.Storage)
	if err != nil {
		return err
	}

	files, err := r.List()
	if err != nil {
		return fmt.Errorf("list: %w", err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
	return nil
}
