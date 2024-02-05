package visa_test

import (
	"github.com/clambin/expensify/payment/visa"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestTransfer_MatchAmount(t *testing.T) {
	transfer := visa.Transfer{
		Amount:   10.0,
		Currency: "EUR",
		Details:  "payment for bar to foo",
	}

	testcases := []struct {
		name   string
		amount float64
		match  bool
	}{
		{
			name:   "match amount",
			amount: 10.0,
			match:  true,
		},
		{
			name:   "amount mismatch ",
			amount: 15.0,
			match:  false,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.match, transfer.MatchAmount(tt.amount))
		})
	}
}

func TestTransfer_MatchDescription(t *testing.T) {
	transfer := visa.Transfer{
		Amount:   10.0,
		Currency: "EUR",
		Details:  "payment for bar to foo",
	}

	testcases := []struct {
		name        string
		description string
		match       bool
	}{
		{
			name:        "match details",
			description: "^payment for bar to foo",
			match:       true,
		},
		{
			name:        "mismatch",
			description: "^payment to bar",
			match:       false,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.match, transfer.MatchDescription(regexp.MustCompile(tt.description)))
		})
	}
}
