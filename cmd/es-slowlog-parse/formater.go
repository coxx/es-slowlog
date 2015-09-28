package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"text/template"

	"github.com/coxx/es-slowlog/internal/parser"
)

func cleanupAddress(addr string) string {
	// add default schema
	if !(strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://")) {
		addr = "http://" + addr
	}
	// trim last slash
	if strings.HasSuffix(addr, "/") {
		addr = strings.TrimSuffix(addr, "/")
	}
	return addr
}

func applyTemplate(t *template.Template, data interface{}) (string, error) {
	b := bytes.Buffer{}
	err := t.Execute(&b, data)
	if err != nil {
		return "", err
	}
	return b.String(), err
}

type formaterFunc func(parser.LogRecord) (string, error)

func defaultFormater() formaterFunc {
	return func(data parser.LogRecord) (string, error) {
		return data.Source, nil
	}
}

func vegetaFormater(targetAddress string) formaterFunc {
	const format = "/{{.Index}}/{{.Types}}/_search?search_type={{.SearchType}}\n@<<BODY\n{{.Source}}\nBODY\n"
	tmpl, _ := template.New("").Parse("GET " + cleanupAddress(targetAddress) + format)
	return func(rec parser.LogRecord) (string, error) {
		s, err := applyTemplate(tmpl, rec)
		return s, err
	}
}

func tankFormater(targetAddress string) (formaterFunc, error) {
	const format = `GET /{{.Index}}/{{.Types}}/_search?search_type={{.SearchType}} HTTP/1.1
Host: %s

{{.Source}}`
	// set Host header
	url, errUrl := url.Parse(cleanupAddress(targetAddress))
	if errUrl != nil {
		return nil, errUrl
	}
	host := url.Host
	tmpl, errTmpl := template.New("").Parse(fmt.Sprintf(format, host))
	if errTmpl != nil {
		return nil, errTmpl
	}

	return func(rec parser.LogRecord) (string, error) {
		s, err := applyTemplate(tmpl, rec)
		if err != nil {
			return "", err
		}
		s = fmt.Sprint(len(s)+2, "\n", s, "\n")
		return s, nil
	}, nil
}

func customFormater(format string) (formaterFunc, error) {
	//var errUnquote error
	format, errUnquote := strconv.Unquote(`"` + format + `"`)
	if errUnquote != nil {
		return nil, errUnquote
	}
	tmpl, errParse := template.New("").Parse(format)
	if errParse != nil {
		return nil, errParse
	}
	return func(rec parser.LogRecord) (string, error) {
		s, err := applyTemplate(tmpl, rec)
		return s, err
	}, nil
}

func newFormater(format, address string) (formaterFunc, error) {
	switch format {
	case "vegeta":
		return vegetaFormater(address), nil
	case "tank":
		return tankFormater(address)
	case "":
		return defaultFormater(), nil
	default:
		return customFormater(format)
	}
}
