package categorizer

import (
	"fmt"
	"github.com/clambin/expensify/payment/checking"
	"github.com/clambin/expensify/payment/visa"
)

type Categorizer struct {
	Rules       Rules
	categorized map[Category]Payments
	unmatched   Payments
}
type Category string

func New(rulesFilename string, paymentFilenames ...string) (*Categorizer, error) {
	rules, err := LoadRulesFromFile(rulesFilename)
	if err != nil {
		return nil, fmt.Errorf("invalid rules file: %w", err)
	}
	c := Categorizer{Rules: rules}
	c.Reset()
	for _, paymentFilename := range paymentFilenames {
		if err = c.Load(paymentFilename); err != nil {
			return nil, fmt.Errorf("%s: %w", paymentFilename, err)
		}

	}
	return &c, nil
}

func (c *Categorizer) Process(payments Payments) {
	for _, transfer := range payments {
		if category, matched := c.matchTransfer(transfer); matched {
			current := c.categorized[category]
			c.categorized[category] = append(current, transfer)
		} else {
			c.unmatched = append(c.unmatched, transfer)
		}
	}
}

func (c *Categorizer) matchTransfer(payment Payment) (Category, bool) {
	for category, rules := range c.Rules {
		for _, rule := range rules {
			if rule.Match(payment) {
				return category, true
			}
		}
	}
	return "", false
}

func (c *Categorizer) Matched() map[Category]Payments {
	return c.categorized
}

func (c *Categorizer) Unmatched() Payments {
	return c.unmatched
}

func (c *Categorizer) Reset() {
	c.categorized = make(map[Category]Payments)
	c.unmatched = nil
}

func (c *Categorizer) Load(filename string) error {
	var payments Payments
	if transfers, err := checking.Parse(filename); err == nil {
		for _, transfer := range transfers {
			payments = append(payments, transfer)
		}
	}
	if transfers, err := visa.Parse(filename); err == nil {
		for _, transfer := range transfers {
			payments = append(payments, transfer)
		}
	}
	if len(payments) == 0 {
		return fmt.Errorf("unrecognized data format")
	}
	c.Process(payments)
	return nil
}
