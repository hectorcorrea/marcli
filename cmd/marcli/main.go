package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/hectorcorrea/marcli/pkg/marc"
)

var fileName, search, fields, exclude, format, hasFields string
var start, count int
var debug bool

func init() {
	flag.StringVar(&fileName, "file", "", "MARC file to process. Required.")
	flag.StringVar(&search, "match", "", "String that must be present in the content of the record, case insensitive.")
	flag.StringVar(&fields, "fields", "", "Comma delimited list of fields to output.")
	flag.StringVar(&exclude, "exclude", "", "Comma delimited list of fields to exclude from the output.")
	flag.StringVar(&format, "format", "mrk", "Output format. Accepted values: mrk, mrc, json, or solr.")
	flag.IntVar(&start, "start", 1, "Number of first record to load")
	flag.IntVar(&count, "count", -1, "Total number of records to load (-1 no limit)")
	flag.StringVar(&hasFields, "hasFields", "", "Comma delimited list of fields that must be present in the record.")
	flag.BoolVar(&debug, "debug", false, "When true it does not stop on errors")
	flag.Parse()
}

func main() {
	if fileName == "" {
		showSyntax()
		return
	}

	params := ProcessFileParams{
		filename:    fileName,
		searchValue: strings.ToLower(search),
		filters:     marc.NewFieldFilters(fields),
		exclude:     marc.NewFieldFilters(exclude),
		start:       start,
		count:       count,
		hasFields:   marc.NewFieldFilters(hasFields),
		debug:       debug,
	}

	if len(params.filters.Fields) > 0 && len(params.exclude.Fields) > 0 {
		panic("Cannot specify fields and exclude at the same time.")
	}

	var err error
	if format == "mrc" {
		err = toMrc(params)
	} else if format == "mrk" {
		err = toMrk(params)
	} else if format == "json" {
		err = toJson(params)
	} else if format == "solr" {
		err = toSolr(params)
	} else {
		err = errors.New("Invalid format")
	}
	if err != nil {
		panic(err)
	}
}

func showSyntax() {
	fmt.Printf("marcli parameters:\r\n")
	fmt.Printf("\r\n")
	flag.PrintDefaults()
	fmt.Printf("\r\n")
	fmt.Printf(`
NOTES:
	The match parameter is used to filter records based on the content of the
values in the record. The hasFields parameter is used to filter records based
on the presence of certain fields on the record (regardless of their value).

	You can only use the fields or exclude parameter, but not both.
`)
	fmt.Printf("\r\n")
	fmt.Printf("\r\n")
}
