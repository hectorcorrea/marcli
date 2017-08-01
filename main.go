package main

import (
	"flag"
	"fmt"
	"strings"
)

var fileName, extract, output string

func init() {
	flag.StringVar(&fileName, "f", "", "MARC file to process. Required.")
	flag.StringVar(&extract, "x", "", "Extract record where string is found (case insensitive).")
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

	filters := NewFieldFilters(output)
	var processor RecordProcessor
	if extract != "" {
		processor = ExtractProcessor{
			Filters: filters,
			Value:   strings.ToLower(extract),
		}
	} else {
		processor = ConsoleProcessor{
			Filters: filters,
		}
	}

	err = file.ReadAll(processor)
	if err != nil {
		panic(err)
	}
}
