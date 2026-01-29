package repo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Repo interface {
	Add(key string, body []byte) error
	List() ([]string, error)
	Get(key string) ([]byte, error)
}

type cfg struct {
	backend Repo
}

type Option func(*cfg)

func New(opts ...Option) Repo {
	var cfg cfg
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg.backend
}

func WithFileSystem(path string) Option {
	return func(c *cfg) {
		c.backend = &fsBackend{root: path}
	}
}

func WithS3Backend(client S3Client, bucket string) Option {
	return func(c *cfg) {
		c.backend = &s3Backend{
			bucket:   bucket,
			client:   client,
			uploader: manager.NewUploader(client),
		}
	}
}

var (
	_, _ Repo = (*fsBackend)(nil), (*s3Backend)(nil)
)

type fsBackend struct {
	root string
}

func (f fsBackend) List() ([]string, error) {
	var files []string
	err := filepath.WalkDir(f.root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, strings.TrimPrefix(path, f.root+string(filepath.Separator)))
		}
		return nil
	})
	return files, err
}

func (f fsBackend) Get(key string) ([]byte, error) {
	path := filepath.Join(f.root, key)
	return os.ReadFile(path)
}

func (f fsBackend) Add(key string, body []byte) error {
	path := filepath.Join(f.root, key)
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	return os.WriteFile(path, body, 0644)
}

type S3Client interface {
	manager.UploadAPIClient
	ListObjects(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type s3Backend struct {
	client   S3Client
	uploader *manager.Uploader
	bucket   string
}

func (s s3Backend) List() ([]string, error) {
	o, err := s.client.ListObjects(context.Background(), &s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("s3: listObjects: %w", err)
	}
	objects := make([]string, 0, len(o.Contents))
	for _, c := range o.Contents {
		objects = append(objects, aws.ToString(c.Key))
	}
	return objects, nil
}

func (s s3Backend) Get(key string) ([]byte, error) {
	o, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("s3: getObject: %w", err)
	}
	defer func() { _ = o.Body.Close() }()
	return io.ReadAll(o.Body)
}

func (s s3Backend) Add(key string, body []byte) error {
	_, err := s.uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewBuffer(body),
	})
	return err
}
