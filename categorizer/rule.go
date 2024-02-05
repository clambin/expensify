package categorizer

import "regexp"

type Rule struct {
	Amount      float64
	Description string
}

func (r Rule) Match(t Payment) bool {
	if r.Amount != 0 {
		if !t.MatchAmount(r.Amount) {
			return false
		}
	}
	if r.Description != "" {
		if !t.MatchDescription(regexp.MustCompile(r.Description)) {
			return false
		}
	}
	return true
}
