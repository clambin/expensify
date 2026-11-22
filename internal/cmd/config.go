package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/alecthomas/kong"
)

// ConfigFile is a flag holding the path of the configuration file. It works like kong.ConfigFlag,
// except that a missing configuration file is not an error: the flags keep their default values.
// It requires kong to be configured with kong.Configuration(loader).
type ConfigFile string

// BeforeResolve loads the configuration file, if it exists, and adds it as a resolver.
func (c ConfigFile) BeforeResolve(k *kong.Kong, ctx *kong.Context, trace *kong.Path) error {
	path := kong.ExpandPath(string(ctx.FlagValue(trace.Flag).(ConfigFile)))
	resolver, err := k.LoadConfig(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("configuration file: %w", err)
	}
	ctx.AddResolver(resolver)
	return nil
}

// EnvFirst wraps a kong configuration loader, so that environment variables take precedence over the configuration file.
//
// kong applies environment variables together with the defaults, i.e. before it consults resolvers,
// so a value in the configuration file would otherwise override the environment.
func EnvFirst(loader kong.ConfigurationLoader) kong.ConfigurationLoader {
	return func(r io.Reader) (kong.Resolver, error) {
		resolver, err := loader(r)
		if err != nil {
			return nil, err
		}
		return kong.ResolverFunc(func(ctx *kong.Context, parent *kong.Path, flag *kong.Flag) (any, error) {
			for _, env := range flag.Envs {
				if _, ok := os.LookupEnv(env); ok {
					return nil, nil
				}
			}
			return resolver.Resolve(ctx, parent, flag)
		}), nil
	}
}
