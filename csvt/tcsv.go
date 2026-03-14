// Package csvt provides a typed CSV parser
package csvt

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Row []any

type File struct {
	SchemaName string
	Rows       []Row
	Schema
}

type Schemas map[string]Schema

type Schema struct {
	// Columns define the columns in the CSV file
	Columns []Column
	// Separator is the separator used in the CSV file for this Column
	Separator rune
}

// Column represents a column in the CSV file
type Column struct {
	// the Header is the name of the column, as it exists in the CSV file
	Header string
	// ColumnType defines the type of the values in the column.
	// During parsing, the value is converted from a string to the appropriate type.
	ColumnType ColumnType
	// Label is the internal name for the column. It's used to represent a row as a map.
	Label string
}

func (s Schema) GetColumns() []string {
	cols := make([]string, 0, len(s.Columns))
	for _, c := range s.Columns {
		if _, ok := c.ColumnType.(IgnoreColumn); !ok {
			cols = append(cols, c.Label)
		}
	}
	return cols
}

func (s Schema) ToMap(row Row) map[string]any {
	value := make(map[string]any, len(row))
	var colIdx int
	for _, col := range s.Columns {
		if _, ok := col.ColumnType.(IgnoreColumn); !ok {
			value[col.Label] = row[colIdx]
			colIdx++
		}
	}
	return value
}

func (s Schemas) Parse(r io.Reader) (File, error) {
	content, err := io.ReadAll(skipBOM(r))
	if err != nil {
		return File{}, fmt.Errorf("failed to read body: %w", err)
	}

	for schemaName := range s {
		if f, err := s.tryParse(content, schemaName); err == nil {
			return f, nil
		}
	}
	return File{}, errors.New("not a valid CSV file: no schema found")
}

func (s Schemas) tryParse(content []byte, schemaName string) (File, error) {
	schema := s[schemaName]

	r := csv.NewReader(bytes.NewReader(content))
	if schema.Separator != 0 {
		r.Comma = schema.Separator
	}
	header, err := r.Read()
	if err != nil {
		return File{}, fmt.Errorf("failed to read header: %w", err)
	}

	if !slices.Equal(header, schema.csvHeader()) {
		return File{}, errors.New("not a valid CSV file: header does not match schema")
	}

	rows, err := r.ReadAll()
	if err != nil {
		return File{}, fmt.Errorf("failed to read all rows: %w", err)
	}

	file := File{Schema: schema, SchemaName: schemaName}
	file.Rows, err = schema.parse(rows)
	if err != nil {
		return File{}, fmt.Errorf("failed to parse rows: %w", err)
	}
	return file, nil
}

func (s Schema) csvHeader() []string {
	headers := make([]string, len(s.Columns))
	for i, c := range s.Columns {
		headers[i] = c.Header
	}
	return headers
}

func (s Schema) parse(records [][]string) ([]Row, error) {
	rows := make([]Row, len(records))

	for i, record := range records {
		if len(record) != len(s.Columns) {
			return nil, fmt.Errorf("invalid row %d: expected %d columns, got %d", i+1, len(s.Columns), len(record))
		}

		rows[i] = make(Row, 0, len(s.Columns))
		for c, cell := range record {
			col := s.Columns[c]
			if _, ok := col.ColumnType.(IgnoreColumn); ok {
				continue
			}

			v, err := col.ColumnType.parse(cell)
			if err != nil {
				return nil, fmt.Errorf("failed to parse row %d, column %d: %w", i+1, c+1, err)
			}
			rows[i] = append(rows[i], v)
		}
	}
	return rows, nil
}

type ColumnType interface {
	parse(string) (any, error)
}

type IgnoreColumn struct{}

func (IgnoreColumn) parse(string) (any, error) {
	return nil, errors.New("column ignored")
}

type StringColumn struct{}

func (StringColumn) parse(s string) (any, error) {
	return s, nil
}

type NumberColumn struct{}

func (f NumberColumn) parse(s string) (any, error) {
	// Trim spaces and optional quotes
	s = strings.TrimSpace(s)
	s = strings.Trim(s, `"'`)

	if s == "" {
		return 0.0, fmt.Errorf("empty amount")
	}

	// Find the last occurrence of each separator
	lastDot := strings.LastIndex(s, ".")
	lastComma := strings.LastIndex(s, ",")

	switch {
	case lastDot >= 0 && lastComma >= 0:
		// Both present → last one is the decimal separator
		if lastDot > lastComma {
			// "." decimal, "," thousands
			s = strings.ReplaceAll(s, ",", "")
		} else {
			// "," decimal, "." thousands
			s = strings.ReplaceAll(s, ".", "")
			s = strings.ReplaceAll(s, ",", ".")
		}

	case lastComma >= 0:
		// TODO: Only comma → assume decimal. else: only dot or no thousands separator → already OK
		s = strings.ReplaceAll(s, ",", ".")
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, fmt.Errorf("parse amount %q: %w", s, err)
	}
	return v, nil
}

type DateColumn struct {
	Format string
}

func (d DateColumn) parse(s string) (any, error) {
	return time.Parse(d.Format, s)
}

// skipBOM skips the UTF-8 BOM (byte order mark) from the beginning of the file
func skipBOM(r io.Reader) io.Reader {
	var prefix [4]byte
	n, err := io.ReadFull(r, prefix[:])
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return r
	}

	var skip int
	for _, bomSequence := range [][]byte{
		{0xEF, 0xBB, 0xBF},
		{0x00, 0x00, 0xFE, 0xFF},
		{0xFF, 0xFE, 0x00, 0x00},
		{0xFE, 0xFF},
		{0xFF, 0xFE},
	} {
		bomSequenceLen := len(bomSequence)
		if n >= bomSequenceLen && bytes.Equal(prefix[:bomSequenceLen], bomSequence) {
			skip = bomSequenceLen
			break
		}
	}

	return io.MultiReader(
		bytes.NewReader(prefix[skip:n]),
		r,
	)
}
