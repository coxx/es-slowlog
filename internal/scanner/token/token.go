package token

// Token is the set of lexical tokens
type Token int

//go:generate stringer -type Token

const (
	Error Token = iota
	EOF
	EOL
	Whitespace
	Comma
	FieldName
	FieldValue
)
