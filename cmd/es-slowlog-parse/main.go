package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"text/template"

	"github.com/coxx/es-slowlog/internal/parser"
)

func main() {
	stderr := log.New(os.Stderr, "", 0)

	flagUsage := flag.Bool("h", false, "Print usage")
	flagFormat := flag.String("f", "{{.Source}}", "Format")
	flag.Parse()
	if *flagUsage {
		flag.Usage()
		os.Exit(1)
	}

	var formatTemplate *template.Template
	if *flagFormat != "" {
		var errParseTemplate error
		format, errUnquote := strconv.Unquote(`"` + *flagFormat + `"`)
		if errUnquote != nil {
			stderr.Fatalf("Bad format: %v\n", errUnquote)
		}

		formatTemplate, errParseTemplate = template.New("").Parse(format)
		if errParseTemplate != nil {
			stderr.Fatalln("Bad format: %v\n", errParseTemplate)
		}
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
