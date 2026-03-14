package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/clambin/expensify/internal/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpload(t *testing.T) {
	tests := []struct {
		name    string
		body    []byte
		wantErr require.ErrorAssertionFunc
		wantKey string
	}{
		{
			name: "bnp-debit",
			body: []byte(`Volgnummer;Uitvoeringsdatum;Valutadatum;Bedrag;Valuta rekening;Rekeningnummer;Type verrichting;Tegenpartij;Naam van de tegenpartij;Mededeling;Details;Status;Reden van weigering
2026-;27/01/2026;27/01/2026;-1372.26;EUR;from account;type;to account;to 1;to 2;message;status;
`),
			wantErr: require.NoError,
			wantKey: "bnp-debit-2026-01-27",
		},
		{
			name: "bnp-visa",
			body: []byte(`null,null,null,null,null,null,null
24/01/2026,25/01/2026,"-18,30",EUR,Vendor,Wisselkoers##,Gerelateerde kost##
`),
			wantErr: require.NoError,
			wantKey: "bnp-visa-2026-01-24",
		},
		{
			name:    "invalid format",
			body:    []byte(`foo,bar`),
			wantErr: require.Error,
		},
		{
			name: "empty",
			body: []byte(`null,null,null,null,null,null,null
`),
			wantErr: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filename := filepath.Join(tmpDir, "/test.csv")
			require.NoError(t, os.WriteFile(filename, tt.body, 0644))

			r := repo.New(repo.WithFileSystem(tmpDir))
			key, err := upload(r, filename)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantKey, key)

			if err != nil {
				return
			}

			body, err := r.Get(key)
			require.NoError(t, err)
			assert.Equal(t, tt.body, body)
		})
	}
}

func TestUpload_Duplicate(t *testing.T) {
	body := []byte(`null,null,null,null,null,null,null
24/01/2026,25/01/2026,"-18,30",EUR,Vendor,Wisselkoers##,Gerelateerde kost##
`)
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "/test.csv")
	require.NoError(t, os.WriteFile(filename, body, 0644))

	r := repo.New(repo.WithFileSystem(tmpDir))
	key, err := upload(r, filename)
	require.NoError(t, err)
	assert.Equal(t, "bnp-visa-2026-01-24", key)

	key, err = upload(r, filename)
	require.NoError(t, err)
	assert.Equal(t, "bnp-visa-2026-01-24-1", key)
}
