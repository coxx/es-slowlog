package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/coxx/es-slowlog/internal/parser"
)

func main() {
	const defaultFormat = "{{.Source}}"

	stderr := log.New(os.Stderr, "", 0)

	flag.Usage = usage
	flagUsage := flag.Bool("h", false, "Print usage")
	flagFormat := flag.String("f", defaultFormat, "Output format")
	flagAddress := flag.String("a", "", "Target address")
	flag.Parse()
	if *flagUsage {
		flag.Usage()
		os.Exit(1)
	}

	if (*flagFormat == "vegeta" || *flagFormat == "tank") && *flagAddress == "" {
		stderr.Fatalln("Target address should be specified")
	}

	formater, err := newFormater(*flagFormat, *flagAddress)
	if err != nil {
		stderr.Fatalf("Bad format: %v\n", err)
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
		s, err := formater(logRecord)
		if err != nil {
			stderr.Fatalln(err)
		}
		fmt.Println(s)
	}
}
