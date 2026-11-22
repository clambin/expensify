package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Configuration creates the storage backend. It is populated from command line flags or a configuration file.
type Configuration struct {
	Type       string                  `name:"type" enum:"filesystem,s3" default:"filesystem" env:"EXPENSIFY_STORAGE" help:"storage backend (${enum})"`
	FileSystem FileSystemConfiguration `embed:"" prefix:"filesystem-" envprefix:"EXPENSIFY_FS_"`
	S3         S3Configuration         `embed:"" prefix:"s3-" envprefix:"EXPENSIFY_S3_"`
}

// FileSystemConfiguration creates the filesystem backend.
type FileSystemConfiguration struct {
	Path string `name:"path" type:"path" env:"PATH" default:"/data" help:"root directory of the filesystem backend"`
}

func (b FileSystemConfiguration) backend() Repo {
	return &fsBackend{root: b.Path}
}

// S3Configuration creates the S3 backend.
type S3Configuration struct {
	URL       string `name:"url" env:"URL" help:"URL of the S3 backend"`
	Bucket    string `name:"bucket" env:"BUCKET" help:"bucket of the S3 backend"`
	AccessKey string `name:"access-key" env:"ACCESS_KEY" help:"access key of the S3 backend"`
	SecretKey string `name:"secret-key" env:"SECRET_KEY" help:"secret key of the S3 backend"`
}

func (b S3Configuration) backend(ctx context.Context) (Repo, error) {
	url := b.URL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("ignored"),
		config.WithBaseEndpoint(url),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(b.AccessKey, b.SecretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("s3 config: %w", err)
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.DisableLogOutputChecksumValidationSkipped = true
	})
	return s3Backend{
		client:   client,
		uploader: transfermanager.New(client),
		bucket:   b.Bucket,
	}, nil
}
