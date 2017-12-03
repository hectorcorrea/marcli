package main

import (
	"flag"
	"fmt"
	"strings"
)

var fileName, search, fields, format string

func init() {
	flag.StringVar(&fileName, "file", "", "MARC file to process. Required.")
	flag.StringVar(&search, "match", "", "Only records that match the string passed, case insensitive.")
	flag.StringVar(&fields, "fields", "", "Comma delimited list of fields to output.")
	flag.StringVar(&format, "format", "mrk", "Output format. Accepted values: mrk, json, or solr.")
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

	searchValue := strings.ToLower(search)
	var processor Processor
	if format == "brown" {
		processor = ProcessorBrown{
			Filters: NewFieldFilters(fields),
		}
	} else if format == "solr" {
		processor = ProcessorSolr{
			Filters: NewFieldFilters(fields),
		}
	} else {
		processor = ConsoleProcessor{
			Filters: NewFieldFilters(fields),
			Format:  format,
		}
	}
	err = file.ReadAll(processor, searchValue)

	if err != nil {
		panic(err)
	}
}
