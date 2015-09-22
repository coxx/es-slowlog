package scanner

import (
	"strings"
	"testing"
)

// Test for Scanner.read()
func Test__Scanner__read(t *testing.T) {
	scanner := New(strings.NewReader("123ABCD"))

	for _, expected := range "123ABCD" {
		got := scanner.read()

		if got != expected {
			t.Errorf("Expected %q, got %q", expected, got)
		}
	}
}

// Test for Scanner.unread()
func Test__Scanner__unread(t *testing.T) {
	scanner := New(strings.NewReader("123ABCD"))

	for _, expected := range "123ABCD" {
		_ = scanner.read()

		scanner.unread()

		got := scanner.read()
		if got != expected {
			t.Errorf("Expected %q, got %q", expected, got)
		}
	}
}

func Test__Scanner__scanToClosingBracket(t *testing.T) {
	testData := []struct {
		input    string
		expected string
	}{
		{"[12345]", "12345"},
		{"[[12345]]", "[12345]"},
		{"[[123[]45]]", "[123[]45]"},
		{"[[123[]45]]xxx", "[123[]45]"},
		{"[123][456]", "123"},
	}

	for _, td := range testData {
		scanner := New(strings.NewReader(td.input))
		got := scanner.scanToClosingBracket()

		if got != td.expected {
			t.Errorf("Expected %q, got %q", td.expected, got)
		}
	}
}
