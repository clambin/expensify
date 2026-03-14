package statements

import (
	"github.com/clambin/expensify/csvt"
	"github.com/clambin/expensify/rules"
)

type TaggedRow struct {
	Tag string
	csvt.Row
}

type TagRule struct {
	Tag  string
	Rule rules.Rule
}

type TagRules []TagRule

func (r TagRules) Evaluate(row map[string]any) (string, bool, error) {
	for _, rule := range r {
		ok, err := rule.Rule.Evaluate(row)
		if err != nil {
			return "", false, err
		}
		if ok {
			return rule.Tag, true, nil
		}
	}
	return "", false, nil
}

func Tag(rows []csvt.Row, s csvt.Schema, rules TagRules) ([]TaggedRow, error) {
	taggedRows := make([]TaggedRow, len(rows))
	for i, row := range rows {
		tp := TaggedRow{Row: row}
		tag, ok, err := rules.Evaluate(s.ToMap(row))
		if err != nil {
			return nil, err
		}
		if ok {
			tp.Tag = tag
		}
		taggedRows[i] = tp
	}
	return taggedRows, nil
}
