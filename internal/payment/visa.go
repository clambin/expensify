package payment

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type VISAPayment struct {
	ExecutionDate time.Time `payment:"Uitvoeringsdatum"`
	ValueDate     time.Time `payment:"Valutadatum"`
	Amount        float64   `payment:"Bedrag"`
	Currency      string    `payment:"Valuta rekening"`
	Details       string    `payment:"Details"`
	ExchangeRate  string    `payment:"Wisselkoers"`
	Related       string    `payment:"Gerelateerde kost"`
}

func (t VISAPayment) GetAmount() float64 {
	return t.Amount
}

func (t VISAPayment) GetDescription() string {
	return t.Details
}

func (t VISAPayment) Match(re *regexp.Regexp) bool {
	return re.MatchString(t.Details)
}

func ParseVISAPayment(record []string) (Payment, error) {
	if len(record) != 7 {
		return nil, fmt.Errorf("wrong number of fields: %d", len(record))
	}
	transfer := VISAPayment{
		Currency:     record[3],
		Details:      record[4],
		ExchangeRate: record[5],
		Related:      record[6],
	}

	var err error
	if transfer.ExecutionDate, err = time.Parse("02/01/2006", record[0]); err != nil {
		return nil, fmt.Errorf("invalid date format (%s): %w", record[1], err)
	}
	if transfer.ValueDate, err = time.Parse("02/01/2006", record[1]); err != nil {
		return nil, fmt.Errorf("invalid date format (%s): %w", record[2], err)
	}
	record[2] = strings.ReplaceAll(record[2], ",", ".")
	if transfer.Amount, err = strconv.ParseFloat(record[2], 64); err != nil {
		return nil, fmt.Errorf("invalid amount (%s): %w", record[3], err)
	}

	return transfer, nil
}
