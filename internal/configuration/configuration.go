package configuration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/internal/statements"
	"github.com/goccy/go-yaml"
)

var (
	defaultConfiguration = Configuration{
		Backend: BackendConfiguration{
			Type: "filesystem",
			FileSystem: BackendFileSystemConfiguration{
				RootDir: "/data",
			},
		},
	}
)

type Configuration struct {
	Backend BackendConfiguration `yaml:"backend"`
	Rules   []statements.TagRule `yaml:"rules"`
}

type BackendConfiguration struct {
	Type       string                         `yaml:"type"`
	FileSystem BackendFileSystemConfiguration `yaml:"filesystem"`
	S3         BackendS3Configuration         `yaml:"s3"`
}

func (b BackendConfiguration) Backend(ctx context.Context) (repo.Repo, error) {
	switch strings.ToLower(b.Type) {
	case "filesystem":
		return b.FileSystem.backend()
	case "s3":
		return b.S3.backend(ctx)
	default:
		return nil, fmt.Errorf("unknown backend type: %s", b.Type)
	}
}

type BackendFileSystemConfiguration struct {
	RootDir string `yaml:"rootDir"`
}

func (b BackendFileSystemConfiguration) backend() (repo.Repo, error) {
	return repo.New(repo.WithFileSystem(b.RootDir)), nil
}

type BackendS3Configuration struct {
	URL       string `yaml:"url"`
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
}

func (b BackendS3Configuration) backend(ctx context.Context) (repo.Repo, error) {
	url := b.URL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("ignored"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(b.AccessKey, b.SecretKey, "")),
		config.WithBaseEndpoint(url),
	)
	if err != nil {
		return nil, fmt.Errorf("s3 config: %w", err)
	}
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.DisableLogOutputChecksumValidationSkipped = true
	})
	return repo.New(repo.WithS3Backend(s3Client, b.Bucket)), nil
}

func LoadFileConfiguration(path string) (Configuration, error) {
	if path == "" {
		var err error
		if path, err = homeConfigPath(); err != nil {
			return Configuration{}, fmt.Errorf("unable to determine user configuration directory: %w", err)
		}
	}
	configuration := defaultConfiguration
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return configuration, nil
		}
		return Configuration{}, fmt.Errorf("open: %w", err)
	}

	defer func() { _ = f.Close() }()
	err = yaml.NewDecoder(f).Decode(&configuration)
	return configuration, err
}

func homeConfigPath() (string, error) {
	path, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("userHomeDir: %w", err)
	}
	path = filepath.Join(path, "expensify")
	if err = os.MkdirAll(path, 0700); err != nil {
		return "", fmt.Errorf("mkdirAll: %w", err)
	}
	return filepath.Join(path, "config.yaml"), nil
}
