package statements

import (
	"bytes"
	"fmt"

	"github.com/clambin/expensify/internal/repo"
	"github.com/clambin/expensify/tcsv"
)

var Schemas = tcsv.Schemas{
	"bnp-debit": {
		Separator: ';',
		Columns: []tcsv.Column{
			{Header: "Volgnummer", ColumnType: tcsv.IgnoreColumn{}},
			{Header: "Uitvoeringsdatum", ColumnType: tcsv.DateColumn{Format: "02/01/2006"}, Label: "Date"},
			{Header: "Valutadatum", ColumnType: tcsv.IgnoreColumn{}},
			{Header: "Bedrag", ColumnType: tcsv.NumberColumn{}, Label: "Amount"},
			{Header: "Valuta rekening", ColumnType: tcsv.IgnoreColumn{}},
			{Header: "Rekeningnummer", ColumnType: tcsv.IgnoreColumn{}},
			{Header: "Type verrichting", ColumnType: tcsv.StringColumn{}, Label: "Type"},
			{Header: "Tegenpartij", ColumnType: tcsv.StringColumn{}, Label: "Counterpart"},
			{Header: "Naam van de tegenpartij", ColumnType: tcsv.StringColumn{}, Label: "CounterpartName"},
			{Header: "Mededeling", ColumnType: tcsv.StringColumn{}, Label: "Message"},
			{Header: "Details", ColumnType: tcsv.StringColumn{}, Label: "Details"},
			{Header: "Status", ColumnType: tcsv.IgnoreColumn{}},
			{Header: "Reden van weigering", ColumnType: tcsv.IgnoreColumn{}},
		},
	},
	"bnp-visa": {
		//separator: ',',
		Columns: []tcsv.Column{
			{Header: "null", ColumnType: tcsv.DateColumn{Format: "02/01/2006"}, Label: "Date"},
			{Header: "null", ColumnType: tcsv.IgnoreColumn{}},
			{Header: "null", ColumnType: tcsv.NumberColumn{}, Label: "Amount"},
			{Header: "null", ColumnType: tcsv.IgnoreColumn{}},
			{Header: "null", ColumnType: tcsv.StringColumn{}, Label: "Details"},
			{Header: "null", ColumnType: tcsv.IgnoreColumn{}},
			{Header: "null", ColumnType: tcsv.IgnoreColumn{}},
		},
	},
}

func Open(r repo.Repo, key string, tagRules TagRules) ([]TaggedRow, tcsv.File, error) {
	body, err := r.Get(key)
	if err != nil {
		return nil, tcsv.File{}, fmt.Errorf("read: %w", err)
	}
	f, err := Schemas.Parse(bytes.NewBuffer(body))
	if err != nil {
		return nil, tcsv.File{}, fmt.Errorf("parse: %w", err)
	}
	taggedStatements, err := Tag(f.Rows, f.Schema, tagRules)
	if err != nil {
		return nil, tcsv.File{}, fmt.Errorf("tag: %w", err)
	}
	return taggedStatements, f, nil
}
