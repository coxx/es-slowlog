package parser

import (
	"io"
	"strconv"
	"strings"

	"github.com/coxx/es-slowlog/internal/scanner"
	"github.com/coxx/es-slowlog/internal/scanner/token"
)

// LogRecord represens single record from the ElasticSearch slowlog file
type LogRecord struct {
	Timestamp   string // 0
	Loglevel    string // 1
	Node        string // 2
	Index       string // 3
	Took        string
	TookMillis  int
	Types       string
	Stats       string
	SearchType  string
	TotalShards int
	Source      string
	ExtraSource string
}

func (record *LogRecord) setFieldByName(name string, v string) {
	switch name {
	case "took":
		record.Took = v
	case "took_millis":
		record.TookMillis, _ = strconv.Atoi(v)
	case "types":
		record.Types = v
	case "stats":
		record.Stats = v
	case "search_type":
		record.SearchType = strings.ToLower(v)
	case "total_shards":
		record.TotalShards, _ = strconv.Atoi(v)
	case "source":
		record.Source = v
	case "extra_source":
		record.ExtraSource = v
	}
}

func (record *LogRecord) setFieldByNo(no int, v string) {
	switch no {
	case 0:
		record.Timestamp = v
	case 1:
		record.Loglevel = v
	case 3:
		record.Node = v
	case 4:
		record.Index = v
	}
}

type Parser struct {
	scanner *scanner.Scanner
}

func New(reader io.Reader) Parser {
	return Parser{scanner: scanner.New(reader)}
}

func (p *Parser) Parse() (LogRecord, error) {
	//type fieldType struct{ name, value string }
	//var fields []fieldType

	record := LogRecord{}

	// scann log line up to EOL

	// var field fieldType
	fieldName := ""
	// fieldValue := ""
	fieldNo := 0

scann:
	for {
		tok, lit := p.scanner.Scan()
		switch tok {
		case token.EOF:
			return record, io.EOF
		case token.EOL:
			break scann
		case token.Comma, token.Whitespace: // ignore whitespace and comma
			continue scann
		case token.FieldName:
			fieldName = lit
			continue scann

		case token.FieldValue:
			if fieldName != "" {
				(&record).setFieldByName(fieldName, lit)
			} else {
				(&record).setFieldByNo(fieldNo, lit)
			}
			fieldNo++
			fieldName = ""
		}
	}
	return record, nil
}
