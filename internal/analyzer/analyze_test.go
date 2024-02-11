package analyzer_test

import (
	"bytes"
	"github.com/clambin/expensify/internal/analyzer"
	"github.com/clambin/expensify/internal/payment"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAnalyze(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantErr       assert.ErrorAssertionFunc
		wantMatched   map[string]payment.Payments
		wantUnmatched payment.Payments
	}{
		{
			name: "valid",
			input: `Uitvoeringsdatum;Valutadatum;Bedrag;Valuta rekening;Details;Wisselkoers;Gerelateerde kost
31/10/2023;01/11/2023;-100;EUR;foo;;
31/10/2023;01/11/2023;-100;EUR;bar;;
31/10/2023;01/11/2023;-100;EUR;snafu;;
`,

			wantErr: assert.NoError,
			wantMatched: map[string]payment.Payments{
				"foo": {payment.VISAPayment{
					ExecutionDate: time.Date(2023, time.October, 31, 0, 0, 0, 0, time.UTC),
					ValueDate:     time.Date(2023, time.November, 1, 0, 0, 0, 0, time.UTC),
					Amount:        -100,
					Currency:      "EUR",
					Details:       "foo",
				}},
				"bar": {payment.VISAPayment{
					ExecutionDate: time.Date(2023, time.October, 31, 0, 0, 0, 0, time.UTC),
					ValueDate:     time.Date(2023, time.November, 1, 0, 0, 0, 0, time.UTC),
					Amount:        -100,
					Currency:      "EUR",
					Details:       "bar",
				}},
			},
			wantUnmatched: payment.Payments{
				payment.VISAPayment{
					ExecutionDate: time.Date(2023, time.October, 31, 0, 0, 0, 0, time.UTC),
					ValueDate:     time.Date(2023, time.November, 1, 0, 0, 0, 0, time.UTC),
					Amount:        -100,
					Currency:      "EUR",
					Details:       "snafu",
				},
			},
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
			t.Parallel()

			rules := analyzer.Rules{
				"foo": []analyzer.Rule{{Description: "^foo"}},
				"bar": []analyzer.Rule{{Description: "^bar"}},
			}

			matched, unmatched, err := analyzer.Analyze(bytes.NewBufferString(tt.input), rules)
			tt.wantErr(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantMatched, matched)
			assert.Equal(t, tt.wantUnmatched, unmatched)
		})
	}
}
