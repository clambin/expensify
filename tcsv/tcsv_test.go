package tcsv_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/clambin/expensify/tcsv"
)

func TestSchema(t *testing.T) {
	s := tcsv.Schemas{
		"test": tcsv.Schema{
			Separator: ',',
			Columns: []tcsv.Column{
				{"name", tcsv.StringColumn{}, "Name"},
				{"ignore", tcsv.IgnoreColumn{}, ""},
				{"age", tcsv.NumberColumn{}, "Age"},
			},
		},
	}

	input := strings.NewReader(`name,ignore,age
John Doe,ignored,35
Jane Doe,ignored,25
`)

	f, err := s.Parse(input)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	want := []string{"Name", "Age"}
	if columns := f.GetColumns(); !slices.Equal(columns, want) {
		t.Fatalf("expected columns to be %v, got %v", want, columns)
	}

	if got := len(f.Rows); got != 2 {
		t.Fatalf("expected 2 rows, got %d", got)
	}
	for i, want := range []map[string]any{
		{"Name": "John Doe", "Age": 35.0},
		{"Name": "Jane Doe", "Age": 25.0},
	} {
		if actual := f.ToMap(f.Rows[i]); !reflect.DeepEqual(actual, want) {
			t.Errorf("row %d: got %v, want %v", i, actual, want)
		}
	}
}
