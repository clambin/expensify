package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/clambin/expensify/internal/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var cli cmd.CLI
	kctx := kong.Parse(&cli,
		kong.Name("expensify"),
		kong.Configuration(cmd.EnvFirst(kongyaml.Loader)),
		kong.BindTo(ctx, (*context.Context)(nil)),
	)
	kctx.FatalIfErrorf(kctx.Run(ctx, &cli))
}
