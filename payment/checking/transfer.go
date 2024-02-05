package checking

import (
	"regexp"
	"strings"
	"time"
)

type Transfer struct {
	ID                  string    `payment:"Volgnummer"`
	ExecutionDate       time.Time `payment:"Uitvoeringsdatum"`
	ValueDate           time.Time `payment:"Valutadatum"`
	Amount              float64   `payment:"Bedrag"`
	Currency            string    `payment:"Valuta rekening"`
	Account             string    `payment:"Rekeningnummer"`
	Type                string    `payment:"Type verrichting"`
	CounterpartyAccount string    `payment:"Tegenpartij"`
	CounterpartyName    string    `payment:"Naam van de tegenpartij"`
	Description         string    `payment:"Mededeling"`
	Details             string    `payment:"Details"`
	Status              string    `payment:"Status"`
	Reason              string    `payment:"Reden van weigering"`
}

func (t Transfer) GetAmount() float64 {
	return t.Amount
}

func (t Transfer) MatchAmount(amount float64) bool {
	return t.Amount == amount
}

func (t Transfer) GetDescription() string {
	return strings.Join([]string{t.Description, t.Details}, "|")
}

func (t Transfer) MatchDescription(re *regexp.Regexp) bool {
	if re.MatchString(t.Description) {
		return true
	}
	return re.MatchString(t.Details)
}

type Transfers []Transfer
