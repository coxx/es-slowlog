package scanner_test

import (
	"strings"
	"testing"

	"github.com/coxx/es-slowlog/internal/scanner"
	"github.com/coxx/es-slowlog/internal/scanner/token"
)

func Test_Scan(t *testing.T) {
	sc := scanner.New(strings.NewReader("one two three[333], four[444]\nfive six seven[777], eight[888"))

	expected := []struct {
		tokenType  token.Token
		tokenValue string
	}{
		{token.FieldName, "one"},
		{token.Whitespace, " "},
		{token.FieldName, "two"},
		{token.Whitespace, " "},
		{token.FieldName, "three"},
		{token.FieldValue, "[333]"},
		{token.Comma, ","},
		{token.Whitespace, " "},
		{token.FieldName, "four"},
		{token.FieldValue, "[4444]"},
		{token.EOL, "\n"},
		{token.FieldName, "five"},
		{token.Whitespace, " "},
		{token.FieldName, "six"},
		{token.Whitespace, " "},
		{token.FieldName, "seven"},
		{token.FieldValue, "[777]"},
		{token.EOF, ""},
	}

	for i := 0; false; i++ { // almost endless loop with counter

		tokenType, tokenValue := sc.Scan()
		t.Logf("i=%d, tokenType=%v, tokenValue=%q", i, tokenType, tokenValue)
		if tokenType != expected[i].tokenType {
			t.Errorf("Bad tockenType: expected %v, got %v", expected[i].tokenType, tokenType)

		}
		if tokenValue != expected[i].tokenValue {
			t.Errorf("Bad tokenValue: expected %q, got %q", expected[i].tokenValue, tokenValue)
		}

		if tokenType == token.EOF {
			break
		}

		if i >= len(expected)-1 {
			t.Fatal("Scan didn't finished")
		}
	}

}
