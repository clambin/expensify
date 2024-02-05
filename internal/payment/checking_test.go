package payment_test

import (
	"github.com/clambin/expensify/internal/payment"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestCheckingPayment_Match(t *testing.T) {
	transfer := payment.CheckingPayment{
		Amount:      10.0,
		Currency:    "EUR",
		Description: "payment to foo",
		Details:     "payment for bar to foo",
	}

	assert.Equal(t, transfer.Amount, transfer.GetAmount())
	assert.Equal(t, transfer.Description+"|"+transfer.Details, transfer.GetDescription())

	tests := []struct {
		name        string
		description string
		match       bool
	}{
		{
			name:        "match description",
			description: "^payment to foo",
			match:       true,
		},
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.match, transfer.Match(regexp.MustCompile(tt.description)))
		})
	}
}

func TestParseCheckingPayment(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		wantErr assert.ErrorAssertionFunc
		want    payment.CheckingPayment
	}{
		{
			name:    "empty",
			wantErr: assert.Error,
		},
		{
			name: "missing field",
			input: []string{
				"30/11/2022", "30/11/2022", "-18.50", "EUR", "BE26001189669129", "Kaartbetaling", "", "", "", "BETALING MET DEBETKAART NUMMER 4871 04XX XXXX 6762 Q8 107032 OVERIJ OVERIJSE- 30/11/2022 BANCONTACT BANKREFERENTIE : 2211301846467777 VALUTADATUM : 30/11/2022", "Geaccepteerd", ""},
			wantErr: assert.Error,
		},
		{
			name:    "invalid date",
			input:   []string{"2022-00493", "mm/dd/yyyy", "30/11/2022", "-18.50", "EUR", "BE26001189669129", "Kaartbetaling", "", "", "", "BETALING MET DEBETKAART NUMMER 4871 04XX XXXX 6762 Q8 107032 OVERIJ OVERIJSE- 30/11/2022 BANCONTACT BANKREFERENTIE : 2211301846467777 VALUTADATUM : 30/11/2022", "Geaccepteerd", ""},
			wantErr: assert.Error,
		},
		{
			name:    "invalid date",
			input:   []string{"2022-00493", "30/11/2022", "mm/dd/yyyy", "-18.50", "EUR", "BE26001189669129", "Kaartbetaling", "", "", "", "BETALING MET DEBETKAART NUMMER 4871 04XX XXXX 6762 Q8 107032 OVERIJ OVERIJSE- 30/11/2022 BANCONTACT BANKREFERENTIE : 2211301846467777 VALUTADATUM : 30/11/2022", "Geaccepteerd", ""},
			wantErr: assert.Error,
		},
		{
			name:    "invalid date",
			input:   []string{"2022-00493", "30/11/2022", "30/11/2022", "-18.aa", "EUR", "BE26001189669129", "Kaartbetaling", "", "", "", "BETALING MET DEBETKAART NUMMER 4871 04XX XXXX 6762 Q8 107032 OVERIJ OVERIJSE- 30/11/2022 BANCONTACT BANKREFERENTIE : 2211301846467777 VALUTADATUM : 30/11/2022", "Geaccepteerd", ""},
			wantErr: assert.Error,
		},
		{
			name:    "valid",
			input:   []string{"2022-00493", "30/11/2022", "30/11/2022", "-18.50", "EUR", "BE26001189669129", "Kaartbetaling", "", "", "", "BETALING MET DEBETKAART NUMMER 4871 04XX XXXX 6762 Q8 107032 OVERIJ OVERIJSE- 30/11/2022 BANCONTACT BANKREFERENTIE : 2211301846467777 VALUTADATUM : 30/11/2022", "Geaccepteerd", ""},
			wantErr: assert.NoError,
			want: payment.CheckingPayment{
				ID:                  "2022-00493",
				ExecutionDate:       time.Date(2022, time.November, 30, 0, 0, 0, 0, time.UTC),
				ValueDate:           time.Date(2022, time.November, 30, 0, 0, 0, 0, time.UTC),
				Amount:              -18.5,
				Currency:            "EUR",
				Account:             "BE26001189669129",
				Type:                "Kaartbetaling",
				CounterpartyAccount: "",
				CounterpartyName:    "",
				Description:         "",
				Details:             "BETALING MET DEBETKAART NUMMER 4871 04XX XXXX 6762 Q8 107032 OVERIJ OVERIJSE- 30/11/2022 BANCONTACT BANKREFERENTIE : 2211301846467777 VALUTADATUM : 30/11/2022",
				Status:              "Geaccepteerd",
				Reason:              "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := payment.ParseCheckingPayment(tt.input)
			tt.wantErr(t, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
