package scanner

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"

	"github.com/coxx/es-slowlog/internal/scanner/token"
)

const eof = '\x00'

type Scanner struct {
	reader *bufio.Reader
}

func New(reader io.Reader) *Scanner {
	return &Scanner{reader: bufio.NewReader(reader)}
}

func (s *Scanner) read() rune {
	c, _, err := s.reader.ReadRune()
	if err != nil {
		return eof
	}
	return c
}

func (s *Scanner) unread() {
	_ = s.reader.UnreadRune()
}

func (s *Scanner) Scan() (token.Token, string) {
	char := s.read()
	//fmt.Printf("---> %q\n", char)
	switch {
	case isWhitespace(char):
		s.unread()
		return token.Whitespace, s.scanWhitespace()
	case char == ',':
		return token.Comma, ","
	case unicode.IsLetter(char):
		s.unread()
		return token.FieldName, s.scanFieldName()
	case char == '[':
		s.unread()
		return token.FieldValue, s.scanToClosingBracket()
	case char == '\n':
		return token.EOL, "\n"
	case char == eof:
		return token.EOF, ""
	}
	return token.Error, ""
}

func isFieldNameChar(r rune) bool {
	if unicode.IsDigit(r) || unicode.IsLetter(r) || r == '_' {
		return true
	}
	return false
}

// isWhitespace returns true if input rune r is whitespace.
// In contrast unocode.IsSpace this function does not treat '\n' 'or '\r' as space characters
func isWhitespace(r rune) bool {
	if r == ' ' || r == '\t' {
		return true
	}
	return false
}

// scanLetters consumes and returns all contiguous letters
func (s *Scanner) scanFieldName() string {
	return s.scanContiguous(isFieldNameChar)
}

// scanWhitespace consumes and returns all contiguous whitespace
func (s *Scanner) scanWhitespace() string {
	return s.scanContiguous(isWhitespace)
}

func (s *Scanner) scanContiguous(isClass func(rune) bool) string {
	buff := bytes.Buffer{}
	for {
		char := s.read()
		if !isClass(char) {
			s.unread()
			break
		}
		buff.WriteRune(char)
	}
	return buff.String()
}

func (s *Scanner) scanToClosingBracket() string {
	buff := &bytes.Buffer{}
	openBracketsCount := 0

	for {
		char := s.read()
		if char == eof {
			break
		}
		switch char {
		case '[':
			openBracketsCount += 1
		case ']':
			openBracketsCount -= 1
		}
		buff.WriteRune(char)

		if openBracketsCount == 0 {
			return strings.TrimSuffix(strings.TrimPrefix(buff.String(), "["), "]")
		}
	}
	// TODO: правильная обработка ошибок
	return "SEQUENCE TOO LONG TO FIND CLOSING BRACKET"
}
