package statements

import (
	"bytes"
	"fmt"

	"github.com/clambin/expensify/csvt"
	"github.com/clambin/expensify/internal/repo"
)

var Schemas = csvt.Schemas{
	"bnp-debit": {
		Separator: ';',
		Columns: []csvt.Column{
			{Header: "Volgnummer", ColumnType: csvt.IgnoreColumn{}},
			{Header: "Uitvoeringsdatum", ColumnType: csvt.DateColumn{Format: "02/01/2006"}, Label: "Date"},
			{Header: "Valutadatum", ColumnType: csvt.IgnoreColumn{}},
			{Header: "Bedrag", ColumnType: csvt.NumberColumn{}, Label: "Amount"},
			{Header: "Valuta rekening", ColumnType: csvt.IgnoreColumn{}},
			{Header: "Rekeningnummer", ColumnType: csvt.IgnoreColumn{}},
			{Header: "Type verrichting", ColumnType: csvt.StringColumn{}, Label: "Type"},
			{Header: "Tegenpartij", ColumnType: csvt.StringColumn{}, Label: "Counterpart"},
			{Header: "Naam van de tegenpartij", ColumnType: csvt.StringColumn{}, Label: "CounterpartName"},
			{Header: "Mededeling", ColumnType: csvt.StringColumn{}, Label: "Message"},
			{Header: "Details", ColumnType: csvt.StringColumn{}, Label: "Details"},
			{Header: "Status", ColumnType: csvt.IgnoreColumn{}},
			{Header: "Reden van weigering", ColumnType: csvt.IgnoreColumn{}},
		},
	},
	"bnp-visa": {
		//separator: ',',
		Columns: []csvt.Column{
			{Header: "null", ColumnType: csvt.DateColumn{Format: "02/01/2006"}, Label: "Date"},
			{Header: "null", ColumnType: csvt.IgnoreColumn{}},
			{Header: "null", ColumnType: csvt.NumberColumn{}, Label: "Amount"},
			{Header: "null", ColumnType: csvt.IgnoreColumn{}},
			{Header: "null", ColumnType: csvt.StringColumn{}, Label: "Details"},
			{Header: "null", ColumnType: csvt.IgnoreColumn{}},
			{Header: "null", ColumnType: csvt.IgnoreColumn{}},
		},
	},
}

func Open(r repo.Repo, key string, tagRules TagRules) ([]TaggedRow, csvt.File, error) {
	body, err := r.Get(key)
	if err != nil {
		return nil, csvt.File{}, fmt.Errorf("read: %w", err)
	}
	f, err := Schemas.Parse(bytes.NewBuffer(body))
	if err != nil {
		return nil, csvt.File{}, fmt.Errorf("parse: %w", err)
	}
	taggedStatements, err := Tag(f.Rows, f.Schema, tagRules)
	if err != nil {
		return nil, csvt.File{}, fmt.Errorf("tag: %w", err)
	}
	return taggedStatements, f, nil
}
