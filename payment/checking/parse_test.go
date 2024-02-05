package checking_test

import (
	"github.com/clambin/expensify/payment/checking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		pass  bool
		count int
	}{
		{name: "input", pass: true, count: 64},
		{name: "empty", pass: true, count: 0},
		{name: "blank", pass: true, count: 0},
		{name: "bad_record", pass: false},
		{name: "missing_field", pass: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join("testdata", tt.name+".csv")
			transfers, err := checking.Parse(filename)
			if tt.pass {
				require.NoError(t, err)
				assert.Len(t, transfers, tt.count)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
