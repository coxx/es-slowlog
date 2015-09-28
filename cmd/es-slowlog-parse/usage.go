package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, `
Supported formats:
  'tank' for Yandex Tank
  'vegeta' for Vegeta
  custom template conforming to text/template
    using following structure as input:
    type LogRecord struct {
        Timestamp   string
        Loglevel    string
        Node        string
        Index       string
        Took        string
        TookMillis  int
        Types       string
        Stats       string
        SearchType  string
        TotalShards int
        Source      string
        ExtraSource string
    }`)
}
