package tcsv

import (
	"io"
	"strings"
	"testing"
)

func TestNumberColumn(t *testing.T) {
	tests := []struct {
		name  string
		input string
		pass  bool
		want  float64
	}{
		{"empty", "", false, 0},
		{"no dot", "125", true, 125.0},
		{"dot is decimal", "12.5", true, 12.5},
		{"dot is decimal, with quotes", "\"12.5\"", true, 12.5},
		{"dot is decimal, thousand marker", "1,012.5", true, 1012.5},
		{"comms is decimal", "12,5", true, 12.5},
		{"comms is decimal, with quotes", "\"12,5\"", true, 12.5},
		{"comma is decimal, thousand marker", "1.012,5", true, 1012.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := NumberColumn{}.parse(tt.input)
			if got := err == nil; got != tt.pass {
				t.Fatalf("parse() error = %v, want pss %v", err, tt.pass)
			}
			if val != tt.want {
				t.Fatalf("parse() = %v, want %v", val, tt.want)
			}
		})
	}
}

func TestSkipBOM(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"utf-8", "\xef\xbb\xbfHello", "Hello"},
		{"utf-16-le", "\xff\xfeHello", "Hello"},
		{"utf-16-be", "\xfe\xffHello", "Hello"},
		{"utf-32-le", "\xff\xfe\x00\x00Hello", "Hello"},
		{"utf-32-be", "\x00\x00\xfe\xffHello", "Hello"},
		{"none", "Hello", "Hello"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := io.ReadAll(skipBOM(strings.NewReader(tt.input)))
			if err != nil {
				t.Fatalf("failed to read: %v", err)
			}
			if got := string(out); got != tt.want {
				t.Fatalf("read() = %v, want %v", got, tt.want)
			}
		})
	}
}
