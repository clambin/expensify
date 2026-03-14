package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFileConfiguration(t *testing.T) {
	configuration, err := LoadFileConfiguration("")
	require.NoError(t, err)
	assert.Equal(t, defaultConfiguration.Backend.Type, configuration.Backend.Type)
	assert.Equal(t, defaultConfiguration.Backend.FileSystem.RootDir, configuration.Backend.FileSystem.RootDir)
}

func TestBackendConfiguration_Backend(t *testing.T) {
	tests := []struct {
		name    string
		cfg     BackendConfiguration
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "filesystem",
			cfg: BackendConfiguration{
				Type: "filesystem",
				FileSystem: BackendFileSystemConfiguration{
					RootDir: "/data",
				},
			},
			wantErr: require.NoError,
		},
		{
			name: "s3",
			cfg: BackendConfiguration{
				Type: "s3",
				S3: BackendS3Configuration{
					URL:       "http://localhost:9000",
					Bucket:    "foo",
					AccessKey: "username",
					SecretKey: "password",
				},
			},
			wantErr: require.NoError,
		},
		{
			name:    "invalid",
			cfg:     BackendConfiguration{Type: "invalid"},
			wantErr: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.cfg.Backend(t.Context())
			tt.wantErr(t, err)
		})
	}
}
