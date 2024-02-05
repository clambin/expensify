package analyzer

import (
	"fmt"
	"github.com/clambin/expensify/internal/payment"
	"io"
)

func Analyze(r io.Reader, rules Rules) (map[string]payment.Payments, payment.Payments, error) {
	payments, err := payment.Read(r)
	if err != nil {
		return nil, nil, fmt.Errorf("read: %w", err)
	}

	matched := make(map[string]payment.Payments)
	var unmatched []payment.Payment
	for _, p := range payments {
		if name, ok := rules.Match(p); ok {
			matched[name] = append(matched[name], p)
		} else {
			unmatched = append(unmatched, p)
		}
	}
	return matched, unmatched, nil
}
