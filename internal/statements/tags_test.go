package statements

import (
	"testing"
	"time"

	"github.com/clambin/expensify/csvt"
	"github.com/clambin/expensify/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagRules_Evaluate(t *testing.T) {
	schema := Schemas["bnp-visa"]
	row := schema.ToMap(csvt.Row{
		time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
		1234.5,
		"foo bar",
	})

	tests := []struct {
		name    string
		rules   TagRules
		wantOK  require.BoolAssertionFunc
		wantTag string
	}{
		{"amount - match", TagRules{{"big", rules.Rule{Field: "Amount", Operator: "gte", Value: 1000}}}, require.True, "big"},
		{"amount - mismatch", TagRules{{"big", rules.Rule{Field: "Amount", Operator: "lt", Value: 1000}}}, require.False, ""},
		{"details - contains", TagRules{{"foo", rules.Rule{Field: "Details", Operator: "contains", Value: "foo"}}}, require.True, "foo"},
		{"details - does not contain", TagRules{{"baz", rules.Rule{Field: "Details", Operator: "contains", Value: "baz"}}}, require.False, ""},
		{"details - contains - list", TagRules{{"comb", rules.Rule{Field: "Details", Operator: "contains", Value: []any{"foo", "baz"}}}}, require.True, "comb"},
		{"combined", TagRules{
			{"small", rules.Rule{Field: "Amount", Operator: "lt", Value: 1000}},
			{"foo", rules.Rule{Field: "Details", Operator: "contains", Value: "foo"}},
		}, require.True, "foo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag, ok, err := tt.rules.Evaluate(row)
			require.NoError(t, err)
			tt.wantOK(t, ok)
			assert.Equal(t, tt.wantTag, tag)
		})
	}
}

func TestTag(t *testing.T) {
	date := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	rows := []csvt.Row{
		{date, 1000.0, "Vendor A"},
		{date, 1200.0, "Vendor A"},
		{date, 2000.0, "Vendor B"},
		{date, 3000.0, "Vendor C"},
	}
	schema := Schemas["bnp-visa"]
	tagRules := TagRules{
		{"A", rules.Rule{Field: "Details", Operator: "contains", Value: "Vendor A"}},
		{"C", rules.Rule{Field: "Details", Operator: "eq", Value: "Vendor C"}},
	}

	taggedRows, err := Tag(rows, schema, tagRules)
	require.NoError(t, err)
	tags := make([]string, len(taggedRows))
	for i, taggedRow := range taggedRows {
		tags[i] = taggedRow.Tag
	}
	assert.Equal(t, []string{"A", "A", "", "C"}, tags)
}
