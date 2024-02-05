package visa

import (
	"regexp"
	"time"
)

type Transfer struct {
	ExecutionDate time.Time `payment:"Uitvoeringsdatum"`
	ValueDate     time.Time `payment:"Valutadatum"`
	Amount        float64   `payment:"Bedrag"`
	Currency      string    `payment:"Valuta rekening"`
	Details       string    `payment:"Details"`
	ExchangeRate  string    `payment:"Wisselkoers"`
	Related       string    `payment:"Gerelateerde kost"`
}

func (t Transfer) GetAmount() float64 {
	return t.Amount
}

func (t Transfer) MatchAmount(amount float64) bool {
	return t.Amount == amount
}

func (t Transfer) GetDescription() string {
	return t.Details
}

func (t Transfer) MatchDescription(re *regexp.Regexp) bool {
	return re.MatchString(t.Details)
}

type Transfers []Transfer
