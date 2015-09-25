package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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

func vegetaFormat(targetAddress string) string {
	const format = "/{{.Index}}/{{.Types}}/_search?search_type={{.SearchType}}\n@<<BODY\n{{.Source}}\nBODY\n"
	return cleanupAddress(targetAddress) + format
}

func main() {
	const defaultFormat = "{{.Source}}"

	stderr := log.New(os.Stderr, "", 0)

	flagUsage := flag.Bool("h", false, "Print usage")
	flagFormat := flag.String("f", defaultFormat, "Format template")
	flagAddress := flag.String("a", "", "Target address")
	flag.Parse()
	if *flagUsage {
		flag.Usage()
		os.Exit(1)
	}

	if *flagFormat == "vegeta" && *flagAddress == "" {
		stderr.Fatalln("Target address should be specified")
	}

	var format string
	switch *flagFormat {
	case "vegeta":
		format = vegetaFormat(*flagAddress)

	case "":
		format = defaultFormat
	default:
		var errUnquote error
		format, errUnquote = strconv.Unquote(`"` + *flagFormat + `"`)
		if errUnquote != nil {
			stderr.Fatalf("Bad format: %v\n", errUnquote)
		}
	}

	log.Printf("format = %q", format)

	var formatTemplate *template.Template
	var errParseTemplate error
	formatTemplate, errParseTemplate = template.New("").Parse(format)
	if errParseTemplate != nil {
		stderr.Fatalln("Bad format: %v\n", errParseTemplate)
	}

	parser := parser.New(os.Stdin)

	for {
		logRecord, err := parser.Parse()
		if err == io.EOF {
			break
		}
		if err != nil {
			stderr.Fatalf("Can't parse input: %v\n", err)
		}

		if formatTemplate != nil {
			b := &bytes.Buffer{}
			if err := formatTemplate.Execute(b, logRecord); err != nil {
				stderr.Fatalln("Can't execute template: %v\n", err)
			}
			fmt.Println(b.String())

		} else {
			fmt.Println(logRecord.Source)
		}

	}

}
