package statements

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatementSchemas(t *testing.T) {
	tests := []struct {
		name           string
		filename       string
		wantSchemaName string
		wantRowCount   int
	}{
		{"debit", "bnp-debit.csv", "bnp-debit", 2},
		{"visa", "bnp-visa.csv", "bnp-visa", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join("testdata", tt.filename))
			require.NoError(t, err)
			file, err := Schemas.Parse(bytes.NewReader(content))
			require.NoError(t, err)
			assert.Equal(t, tt.wantSchemaName, file.SchemaName)
			assert.Len(t, file.Rows, tt.wantRowCount)
			wantColumns := len(file.GetColumns())
			for _, row := range file.Rows {
				assert.Len(t, row, wantColumns)
			}
		})
	}
}
