package payment

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
)

type Payment interface {
	GetAmount() float64
	GetDescription() string
	Match(regexp *regexp.Regexp) bool
}

var _ Payment = CheckingPayment{}
var _ Payment = VISAPayment{}

type Payments []Payment

func (p Payments) Total() float64 {
	var total float64
	for i := range p {
		total += p[i].GetAmount()
	}
	return total
}

func Read(r io.Reader) (Payments, error) {
	var input bytes.Buffer
	if payments, err := ReadPayments(io.TeeReader(r, &input), ParseVISAPayment); err == nil {
		return payments, nil
	}
	return ReadPayments(&input, ParseCheckingPayment)
}

func ReadPayments(r io.Reader, parse func([]string) (Payment, error)) (Payments, error) {
	parser := csv.NewReader(r)
	parser.Comma = ';'
	records, err := parser.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("csv: %w", err)
	}
	if len(records) < 2 {
		return nil, nil
	}

	transfers := make(Payments, len(records)-1)
	for i := range records[1:] {
		if transfers[i], err = parse(records[i+1]); err != nil {
			return nil, fmt.Errorf("parse: %w", err)
		}
	}
	return transfers, nil
}
