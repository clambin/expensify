package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/statements"
	"github.com/clambin/expensify/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLI_StorageConfiguration(t *testing.T) {
	const configFile = `
storage:
  type: s3
  filesystem:
    path: /from/file
  s3:
    url: http://from-file:9000
    bucket: file-bucket
    access-key: file-access-key
    secret-key: file-secret-key
`
	tests := []struct {
		name   string
		config string
		env    map[string]string
		args   []string
		want   repo.Configuration
	}{
		{
			name: "defaults",
			args: []string{"list"},
			want: repo.Configuration{Type: "filesystem", FileSystem: repo.FileSystemConfiguration{Path: "/data"}},
		},
		{
			name:   "config file overrides defaults",
			config: configFile,
			args:   []string{"list"},
			want: repo.Configuration{
				Type:       "s3",
				FileSystem: repo.FileSystemConfiguration{Path: "/from/file"},
				S3: repo.S3Configuration{
					URL:       "http://from-file:9000",
					Bucket:    "file-bucket",
					AccessKey: "file-access-key",
					SecretKey: "file-secret-key",
				},
			},
		},
		{
			name:   "arguments override config file",
			config: configFile,
			args:   []string{"list", "--storage-type", "filesystem", "--storage-filesystem-path", "/from/args", "--storage-s3-access-key", "args-access-key"},
			want: repo.Configuration{
				Type:       "filesystem",
				FileSystem: repo.FileSystemConfiguration{Path: "/from/args"},
				S3: repo.S3Configuration{
					URL:       "http://from-file:9000",
					Bucket:    "file-bucket",
					AccessKey: "args-access-key",
					SecretKey: "file-secret-key",
				},
			},
		},
		{
			name:   "environment overrides config file",
			config: configFile,
			env:    map[string]string{"EXPENSIFY_S3_BUCKET": "env-bucket"},
			args:   []string{"list"},
			want: repo.Configuration{
				Type:       "s3",
				FileSystem: repo.FileSystemConfiguration{Path: "/from/file"},
				S3: repo.S3Configuration{
					URL:       "http://from-file:9000",
					Bucket:    "env-bucket",
					AccessKey: "file-access-key",
					SecretKey: "file-secret-key",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			path := filepath.Join(t.TempDir(), "config.yaml")
			require.NoError(t, os.WriteFile(path, []byte(tt.config), 0600))

			var cli CLI
			parser, err := kong.New(&cli, kong.Configuration(EnvFirst(kongyaml.Loader)))
			require.NoError(t, err)
			_, err = parser.Parse(append([]string{"--config", path}, tt.args...))
			require.NoError(t, err)
			assert.Equal(t, tt.want, cli.Storage)
		})
	}
}

func TestConfigFile(t *testing.T) {
	// kong expands "~" through the OS user database, so tests can't redirect the real default path.
	// Use a grammar whose default configuration path points to a temporary directory instead.
	type grammar struct {
		Config  ConfigFile         `short:"c" type:"path" default:"${config_file}"`
		Storage repo.Configuration `embed:"" prefix:"storage-"`
	}
	parse := func(t *testing.T, defaultPath string, args ...string) (grammar, error) {
		t.Helper()
		var cli grammar
		parser, err := kong.New(&cli, kong.Configuration(EnvFirst(kongyaml.Loader)), kong.Vars{"config_file": defaultPath})
		require.NoError(t, err)
		_, err = parser.Parse(args)
		return cli, err
	}

	t.Run("default file is loaded", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.yaml")
		require.NoError(t, os.WriteFile(path, []byte("storage:\n  type: s3\n"), 0600))
		cli, err := parse(t, path)
		require.NoError(t, err)
		assert.Equal(t, "s3", cli.Storage.Type)
	})

	t.Run("missing default file uses defaults", func(t *testing.T) {
		cli, err := parse(t, filepath.Join(t.TempDir(), "missing.yaml"))
		require.NoError(t, err)
		assert.Equal(t, "filesystem", cli.Storage.Type)
		assert.Equal(t, "/data", cli.Storage.FileSystem.Path)
	})

	t.Run("missing explicit file uses defaults", func(t *testing.T) {
		cli, err := parse(t, filepath.Join(t.TempDir(), "default.yaml"), "--config", filepath.Join(t.TempDir(), "missing.yaml"))
		require.NoError(t, err)
		assert.Equal(t, "filesystem", cli.Storage.Type)
	})

	t.Run("invalid file is an error", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.yaml")
		require.NoError(t, os.WriteFile(path, []byte("storage: [\n"), 0600))
		_, err := parse(t, path)
		assert.Error(t, err)
	})
}

const testRules = `
- tag: groceries
  rule:
    operator: eq
    field: Vendor
    value: shop
`

var wantTestRules = statements.TagRules{
	{Tag: "groceries", Rule: rules.Rule{Operator: "eq", Field: "Vendor", Value: "shop"}},
}

func TestCLI_rules(t *testing.T) {
	tests := []struct {
		name      string
		rulesFile string // where to write the rules, relative to the configuration directory
		rules     func(configDir string) string
		body      string
		wantErr   require.ErrorAssertionFunc
		want      statements.TagRules
	}{
		{
			name:      "relative to configuration directory",
			rulesFile: "rules.yaml",
			rules:     func(string) string { return "rules.yaml" },
			body:      testRules,
			wantErr:   require.NoError,
			want:      wantTestRules,
		},
		{
			name:      "relative subdirectory",
			rulesFile: "sub/rules.yaml",
			rules:     func(string) string { return "sub/rules.yaml" },
			body:      testRules,
			wantErr:   require.NoError,
			want:      wantTestRules,
		},
		{
			name:      "absolute",
			rulesFile: "elsewhere/rules.yaml",
			rules:     func(dir string) string { return filepath.Join(dir, "elsewhere", "rules.yaml") },
			body:      testRules,
			wantErr:   require.NoError,
			want:      wantTestRules,
		},
		{
			name:    "missing",
			rules:   func(string) string { return "missing.yaml" },
			wantErr: require.Error,
		},
		{
			name:      "invalid",
			rulesFile: "rules.yaml",
			rules:     func(string) string { return "rules.yaml" },
			body:      "- tag: [",
			wantErr:   require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configDir := t.TempDir()
			if tt.rulesFile != "" {
				path := filepath.Join(configDir, tt.rulesFile)
				require.NoError(t, os.MkdirAll(filepath.Dir(path), 0700))
				require.NoError(t, os.WriteFile(path, []byte(tt.body), 0600))
			}

			cli := CLI{
				Config: ConfigFile(filepath.Join(configDir, "config.yaml")),
				Rules:  tt.rules(configDir),
			}
			got, err := cli.rules()
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestCLI_rules_Parsed checks rules() with values set by kong: the default rules path, and a rules path set in the configuration file.
func TestCLI_rules_Parsed(t *testing.T) {
	tests := []struct {
		name      string
		config    string
		rulesFile string
	}{
		{name: "default", rulesFile: "rules.yaml"},
		{name: "configuration file", config: "rules: custom/rules.yaml\n", rulesFile: "custom/rules.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// run from a different directory, so relative paths can't accidentally resolve against the working directory
			t.Chdir(t.TempDir())

			configDir := t.TempDir()
			configPath := filepath.Join(configDir, "config.yaml")
			require.NoError(t, os.WriteFile(configPath, []byte(tt.config), 0600))
			rulesPath := filepath.Join(configDir, tt.rulesFile)
			require.NoError(t, os.MkdirAll(filepath.Dir(rulesPath), 0700))
			require.NoError(t, os.WriteFile(rulesPath, []byte(testRules), 0600))

			var cli CLI
			parser, err := kong.New(&cli, kong.Configuration(EnvFirst(kongyaml.Loader)))
			require.NoError(t, err)
			_, err = parser.Parse([]string{"--config", configPath, "list"})
			require.NoError(t, err)

			got, err := cli.rules()
			require.NoError(t, err)
			assert.Equal(t, wantTestRules, got)
		})
	}
}
