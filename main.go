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
	searchValue := strings.ToLower(search)
	filters := NewFieldFilters(fields)
	// TODO: handle multiple formats (brown, solr, json)
	err := mrkProcessor(fileName, searchValue, filters)
	if err != nil {
		panic(err)
	}
}
