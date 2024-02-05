package visa

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func Parse(filename string) (Transfers, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	parser := csv.NewReader(f)
	parser.Comma = ';'
	records, err := parser.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parse payment: %w", err)
	}
	var transfers Transfers
	if len(records) == 0 {
		return transfers, nil
	}
	for _, record := range records[1:] {
		var transfer Transfer
		if transfer, err = parseRecord(record); err != nil {
			return transfers, err
		}
		transfers = append(transfers, transfer)
	}
	return transfers, nil
}

func parseRecord(record []string) (transfer Transfer, err error) {
	if len(record) != 7 {
		return Transfer{}, fmt.Errorf("wrong number of fields: %d", len(record))
	}
	transfer = Transfer{
		Currency:     record[3],
		Details:      record[4],
		ExchangeRate: record[5],
		Related:      record[6],
	}

	if transfer.ExecutionDate, err = time.Parse("02/01/2006", record[0]); err != nil {
		return Transfer{}, fmt.Errorf("invalid date format (%s): %w", record[1], err)
	}
	if transfer.ValueDate, err = time.Parse("02/01/2006", record[1]); err != nil {
		return Transfer{}, fmt.Errorf("invalid date format (%s): %w", record[2], err)
	}
	record[2] = strings.ReplaceAll(record[2], ",", ".")
	if transfer.Amount, err = strconv.ParseFloat(record[2], 32); err != nil {
		return Transfer{}, fmt.Errorf("invalid amount (%s): %w", record[3], err)
	}

	return transfer, nil
}
