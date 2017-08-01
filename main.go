package main

import (
	"flag"
	"fmt"
	"strings"
)

var fileName, search, output string

func init() {
	flag.StringVar(&fileName, "f", "", "MARC file to process. Required.")
	flag.StringVar(&search, "m", "", "Only records that match the string passed (case insensitive).")
	flag.StringVar(&output, "o", "", "Comma delimited list of fields to output.")
	flag.Parse()
}

func main() {
	if fileName == "" {
		fmt.Printf("marcli parameters:\r\n")
		flag.PrintDefaults()
		return
	}

	file, err := NewMarcFile(fileName)
	if err != nil {
		panic(err)
	}

	processor := ConsoleProcessor{
		Filters:     NewFieldFilters(output),
		SearchValue: strings.ToLower(search),
	}

	err = file.ReadAll(processor)
	if err != nil {
		panic(err)
	}
}
