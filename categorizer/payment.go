package categorizer

import "regexp"

type Payment interface {
	GetAmount() float64
	MatchAmount(float64) bool
	GetDescription() string
	MatchDescription(regexp *regexp.Regexp) bool
}

type Payments []Payment

func (c Payments) Total() float64 {
	var sum float64
	for _, entry := range c {
		sum += entry.GetAmount()
	}
	return sum
}
