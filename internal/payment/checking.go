package payment

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CheckingPayment struct {
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

func (t CheckingPayment) GetAmount() float64 {
	return t.Amount
}

func (t CheckingPayment) GetDescription() string {
	return strings.Join([]string{t.Description, t.Details}, "|")
}

func (t CheckingPayment) Match(re *regexp.Regexp) bool {
	if re.MatchString(t.Description) {
		return true
	}
	return re.MatchString(t.Details)
}

func ParseCheckingPayment(record []string) (Payment, error) {
	if len(record) != 13 {
		return nil, fmt.Errorf("wrong number of fields: %d", len(record))
	}
	transfer := CheckingPayment{
		ID:                  record[0],
		Currency:            record[4],
		Account:             record[5],
		Type:                record[6],
		CounterpartyAccount: record[7],
		CounterpartyName:    record[8],
		Description:         record[9],
		Details:             record[10],
		Status:              record[11],
		Reason:              record[12],
	}

	var err error
	if transfer.ExecutionDate, err = time.Parse("02/01/2006", record[1]); err != nil {
		return nil, fmt.Errorf("invalid date format (%s): %w", record[1], err)
	}
	if transfer.ValueDate, err = time.Parse("02/01/2006", record[2]); err != nil {
		return nil, fmt.Errorf("invalid date format (%s): %w", record[2], err)
	}
	record[3] = strings.Replace(record[3], ",", ".", -1)
	if transfer.Amount, err = strconv.ParseFloat(record[3], 64); err != nil {
		return nil, fmt.Errorf("invalid amount (%s): %w", record[3], err)
	}

	return transfer, nil
}
