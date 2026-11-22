package repo

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/minio"
)

func TestRepo(t *testing.T) {
	t.Run("fs", func(t *testing.T) {
		testRepo(t, FileSystemConfiguration{Path: t.TempDir()}.backend())
	})

	t.Run("s3", func(t *testing.T) {
		ctx := t.Context()
		s3Container, err := minio.Run(ctx, "pgsty/silo:latest")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, testcontainers.TerminateContainer(s3Container))
		}()
		ep, err := s3Container.ConnectionString(ctx)
		require.NoError(t, err)

		cfg, err := config.LoadDefaultConfig(
			ctx,
			config.WithRegion("ignored"),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(s3Container.Username, s3Container.Password, ""),
			),
		)

		require.NoError(t, err)
		s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String("http://" + ep)
			o.UsePathStyle = true
		})

		_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String("expensify"),
		})
		require.NoError(t, err)

		b, err := S3Configuration{
			URL:       "http://" + ep,
			Bucket:    "expensify",
			AccessKey: s3Container.Username,
			SecretKey: s3Container.Password,
		}.backend(t.Context())
		require.NoError(t, err)

		testRepo(t, b)
	})
}

func testRepo(t testing.TB, b Repo) {
	t.Helper()

	// list empty repo
	files, err := b.List()
	require.NoError(t, err)
	assert.Len(t, files, 0)

	// get a non-existent file
	_, err = b.Get("foo")
	require.Error(t, err)

	files, _ = b.List()
	assert.Empty(t, files)

	// add a new file
	err = b.Add("foo", []byte("bar"))
	require.NoError(t, err)

	files, _ = b.List()
	assert.Equal(t, []string{"foo"}, files)

	// get the file
	body, err := b.Get("foo")
	require.NoError(t, err)
	assert.Equal(t, "bar", string(body))

	// keys can be nested
	err = b.Add("bar/foo", []byte("baz"))
	require.NoError(t, err)
	body, err = b.Get("bar/foo")
	require.NoError(t, err)
	assert.Equal(t, "baz", string(body))

	// list files
	files, err = b.List()
	require.NoError(t, err)
	assert.Equal(t, []string{"bar/foo", "foo"}, files)
}
