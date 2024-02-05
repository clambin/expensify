package payment_test

import (
	"github.com/clambin/expensify/internal/payment"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestVISAPayment_Match(t *testing.T) {
	transfer := payment.VISAPayment{
		Amount:   10.0,
		Currency: "EUR",
		Details:  "payment for bar to foo",
	}

	assert.Equal(t, transfer.Amount, transfer.GetAmount())
	assert.Equal(t, transfer.Details, transfer.GetDescription())

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
			assert.Equal(t, tt.match, transfer.Match(regexp.MustCompile(tt.description)))
		})
	}
}

func TestParseVISAPayment(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		wantErr assert.ErrorAssertionFunc
		want    payment.VISAPayment
	}{
		{
			name:    "empty",
			wantErr: assert.Error,
		},
		{
			name:    "missing field",
			input:   []string{"30/11/2022", "01/12/2022", "-3,aa", "EUR", "SODEXO 4726 SWIFT ADE     NL - 01310  LA HULPE"},
			wantErr: assert.Error,
		},
		{
			name:    "invalid amount",
			input:   []string{"30/11/2022", "01/12/2022", "-3,aa", "EUR", "SODEXO 4726 SWIFT ADE     NL - 01310  LA HULPE", "", ""},
			wantErr: assert.Error,
		},
		{
			name:    "invalid date",
			input:   []string{"dd/mm/yyyy", "01/12/2022", "-3,25", "EUR", "SODEXO 4726 SWIFT ADE     NL - 01310  LA HULPE", "", ""},
			wantErr: assert.Error,
		},
		{
			name:    "invalid date",
			input:   []string{"30/11/2022", "dd/mm/yyyy", "-3,25", "EUR", "SODEXO 4726 SWIFT ADE     NL - 01310  LA HULPE", "", ""},
			wantErr: assert.Error,
		},
		{
			name:    "valid",
			input:   []string{"30/11/2022", "01/12/2022", "-3,25", "EUR", "SODEXO 4726 SWIFT ADE     NL - 01310  LA HULPE", "", ""},
			wantErr: assert.NoError,
			want: payment.VISAPayment{
				ExecutionDate: time.Date(2022, time.November, 30, 0, 0, 0, 0, time.UTC),
				ValueDate:     time.Date(2022, time.December, 1, 0, 0, 0, 0, time.UTC),
				Amount:        -3.25,
				Currency:      "EUR",
				Details:       "SODEXO 4726 SWIFT ADE     NL - 01310  LA HULPE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := payment.ParseVISAPayment(tt.input)
			tt.wantErr(t, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
