package payment_test

import (
	"bytes"
	"github.com/clambin/expensify/internal/payment"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPayments_Total(t *testing.T) {
	var payments payment.Payments
	assert.Zero(t, payments.Total())

	payments = append(payments, payment.CheckingPayment{Amount: 10}, payment.VISAPayment{Amount: 100})
	assert.Equal(t, 110.0, payments.Total())
}

func TestRead(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr assert.ErrorAssertionFunc
		want    payment.Payments
	}{
		{
			name: "visa",
			input: `Uitvoeringsdatum;Valutadatum;Bedrag;Valuta rekening;Details;Wisselkoers;Gerelateerde kost
31/10/2023;01/11/2023;-93,43;EUR;foo;Wisselkoers##;Gerelateerde kost##
`,
			wantErr: assert.NoError,
			want: payment.Payments{payment.VISAPayment{
				ExecutionDate: time.Date(2023, time.October, 31, 0, 0, 0, 0, time.UTC),
				ValueDate:     time.Date(2023, time.November, 1, 0, 0, 0, 0, time.UTC),
				Amount:        -93.43,
				Currency:      "EUR",
				Details:       "foo",
				ExchangeRate:  "Wisselkoers##",
				Related:       "Gerelateerde kost##",
			}},
		},
		{
			name: "checking",
			input: `Volgnummer;Uitvoeringsdatum;Valutadatum;Bedrag;Valuta rekening;Rekeningnummer;Type verrichting;Tegenpartij;Naam van de tegenpartij;Mededeling;Details;Status;Reden van weigering
id;31/10/2023;31/10/2023;-11.00;EUR;foo;Kaartbetaling;;;;bar;Geaccepteerd;
`,
			wantErr: assert.NoError,
			want: payment.Payments{payment.CheckingPayment{
				ID:            "id",
				ExecutionDate: time.Date(2023, time.October, 31, 0, 0, 0, 0, time.UTC),
				ValueDate:     time.Date(2023, time.October, 31, 0, 0, 0, 0, time.UTC),
				Amount:        -11, Currency: "EUR",
				Account: "foo",
				Type:    "Kaartbetaling",
				Details: "bar",
				Status:  "Geaccepteerd",
			}},
		},
		{
			name: "invalid",
			input: `Volgnummer;Uitvoeringsdatum;Valutadatum;Bedrag;Valuta rekening;Rekeningnummer;Type verrichting;Tegenpartij;Naam van de tegenpartij;Mededeling;Details;Status;Reden van weigering
id;31/10/2023;31/10/2023;-11.00;EUR;foo;Kaartbetaling;;;;bar;
`,
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := payment.Read(bytes.NewBufferString(tt.input))
			tt.wantErr(t, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestReadPayments(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		f       func([]string) (payment.Payment, error)
		wantErr assert.ErrorAssertionFunc
		want    payment.Payments
	}{
		{
			name:    "visa - empty",
			f:       payment.ParseVISAPayment,
			wantErr: assert.NoError,
		},
		{
			name:    "checking - empty",
			f:       payment.ParseCheckingPayment,
			wantErr: assert.NoError,
		},
		{
			name: "visa - invalid",
			input: `Uitvoeringsdatum;Valutadatum;Bedrag;Valuta rekening;Details;Wisselkoers;Gerelateerde kost
mm/dd/yyyy;01/11/2023;-93,43;EUR;foo;Wisselkoers##;Gerelateerde kost##
`,
			f:       payment.ParseVISAPayment,
			wantErr: assert.Error,
		},
		{
			name: "visa - valid",
			input: `Uitvoeringsdatum;Valutadatum;Bedrag;Valuta rekening;Details;Wisselkoers;Gerelateerde kost
31/10/2023;01/11/2023;-93,43;EUR;foo;Wisselkoers##;Gerelateerde kost##
`,
			f:       payment.ParseVISAPayment,
			wantErr: assert.NoError,
			want: payment.Payments{
				payment.VISAPayment{
					ExecutionDate: time.Date(2023, time.October, 31, 0, 0, 0, 0, time.UTC),
					ValueDate:     time.Date(2023, time.November, 1, 0, 0, 0, 0, time.UTC),
					Amount:        -93.43,
					Currency:      "EUR",
					Details:       "foo",
					ExchangeRate:  "Wisselkoers##",
					Related:       "Gerelateerde kost##",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewBufferString(tt.input)
			got, err := payment.ReadPayments(r, tt.f)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
