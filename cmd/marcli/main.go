package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/hectorcorrea/marcli/pkg/marc"
)

var fileName, search, fields, format string
var start, count int

func init() {
	flag.StringVar(&fileName, "file", "", "MARC file to process. Required.")
	flag.StringVar(&search, "match", "", "Only records that match the string passed, case insensitive.")
	flag.StringVar(&fields, "fields", "", "Comma delimited list of fields to output.")
	flag.StringVar(&format, "format", "mrk", "Output format. Accepted values: mrk, mrc, json, or solr.")
	flag.IntVar(&start, "start", 1, "Number of first record to load")
	flag.IntVar(&count, "count", -1, "Total number of records to load (-1 no limit)")

	flag.Parse()
}

func main() {
	if fileName == "" {
		fmt.Printf("marcli parameters:\r\n")
		flag.PrintDefaults()
		return
	}
	var err error
	searchValue := strings.ToLower(search)
	filters := marc.NewFieldFilters(fields)
	if format == "mrc" {
		err = toMrc(fileName, searchValue, filters, start, count)
	} else if format == "mrk" {
		err = toMrk(fileName, searchValue, filters, start, count)
	} else if format == "json" {
		err = toJson(fileName, searchValue, filters, start, count)
	} else if format == "solr" {
		err = toSolr(fileName, searchValue, filters, start, count)
	} else {
		err = errors.New("Invalid format")
	}
	if err != nil {
		panic(err)
	}
}
