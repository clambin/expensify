package analyzer_test

import (
	"bytes"
	"github.com/clambin/expensify/internal/analyzer"
	"github.com/clambin/expensify/internal/payment"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRules_Match(t *testing.T) {
	testcases := []struct {
		name string
		rule analyzer.Rule
		want assert.BoolAssertionFunc
	}{
		{
			name: "empty rule matches",
			rule: analyzer.Rule{},
			want: assert.True,
		},
		{
			name: "match",
			rule: analyzer.Rule{Description: "^payment to foo$"},
			want: assert.True,
		},
		{
			name: "mismatch",
			rule: analyzer.Rule{Description: "^payment to bar$"},
			want: assert.False,
		},
		{
			name: "no regexp",
			rule: analyzer.Rule{Description: " foo"},
			want: assert.True,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			p := payment.CheckingPayment{
				ValueDate:   time.Date(2022, time.December, 19, 0, 0, 0, 0, time.UTC),
				Description: "payment to foo",
			}

			rules := analyzer.Rules{
				"foo": []analyzer.Rule{tt.rule},
			}
			name, got := rules.Match(p)
			tt.want(t, got)
			if got {
				assert.Equal(t, "foo", name)
			}
		})
	}
}

func TestLoadRules(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr assert.ErrorAssertionFunc
		want    analyzer.Rules
	}{
		{
			name:    "empty",
			wantErr: assert.NoError,
		},
		{
			name:    "blank",
			input:   "\n",
			wantErr: assert.NoError,
		},
		{
			name: "valid",
			input: `
foo:
  - description: "bar"
  - description: "snafu"
`,
			wantErr: assert.NoError,
			want: map[string][]analyzer.Rule{
				"foo": {{Description: "bar"}, {Description: "snafu"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := analyzer.LoadRules(bytes.NewBufferString(tt.input))
			tt.wantErr(t, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
