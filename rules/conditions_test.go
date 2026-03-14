package rules

import (
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule_Evaluate(t *testing.T) {
	tests := []struct {
		name    string
		rules   []byte
		input   map[string]any
		want    bool
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "single rule",
			rules: []byte(`
operator: eq
field: foo
value: bar
`),
			input:   map[string]any{"foo": "bar"},
			want:    true,
			wantErr: require.NoError,
		},
		{
			name: "all",
			rules: []byte(`
all:
- operator: eq
  field: foo
  value: bar
- operator: lte
  field: age
  value: 23
`),
			input:   map[string]any{"foo": "bar", "age": 23},
			want:    true,
			wantErr: require.NoError,
		},
		{
			name: "any",
			rules: []byte(`
any:
- operator: eq
  field: foo
  value: bar
- operator: lte
  field: age
  value: 23
`),
			input:   map[string]any{"foo": "bar", "age": 24},
			want:    true,
			wantErr: require.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r Rule
			require.NoError(t, yaml.Unmarshal(tt.rules, &r))

			ok, err := r.Evaluate(tt.input)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, ok)
		})
	}
}
