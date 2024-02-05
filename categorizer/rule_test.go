package categorizer_test

import (
	"github.com/clambin/expensify/categorizer"
	"github.com/clambin/expensify/payment/checking"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRule_Match(t *testing.T) {
	transfer := checking.Transfer{
		ValueDate:   time.Date(2022, time.December, 19, 0, 0, 0, 0, time.UTC),
		Amount:      10.0,
		Description: "payment to foo",
	}

	testcases := []struct {
		name  string
		rule  categorizer.Rule
		match bool
	}{
		{
			name:  "empty rule matches",
			rule:  categorizer.Rule{},
			match: true,
		},
		{
			name:  "match amount",
			rule:  categorizer.Rule{Amount: 10.0},
			match: true,
		},
		{
			name:  "amount mismatch",
			rule:  categorizer.Rule{Amount: 15.0},
			match: false,
		},
		{
			name:  "match description",
			rule:  categorizer.Rule{Description: "^payment to foo$"},
			match: true,
		},
		{
			name:  "description mismatch",
			rule:  categorizer.Rule{Description: "^payment to bar$"},
			match: false,
		},
		{
			name:  "match both",
			rule:  categorizer.Rule{Amount: 10.0, Description: "^payment to foo$"},
			match: true,
		},
		{
			name:  "must match both",
			rule:  categorizer.Rule{Amount: 10.0, Description: "^payment to bar$"},
			match: false,
		},
		{
			name:  "no regexp",
			rule:  categorizer.Rule{Description: " foo"},
			match: true,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.match, tt.rule.Match(transfer))
		})
	}

	rule := categorizer.Rule{
		Description: "foo",
	}

	matched := rule.Match(checking.Transfer{Description: "foo"})
	assert.True(t, matched)
}
