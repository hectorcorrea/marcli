package main

import (
	"errors"
	"flag"
	"fmt"
	"marcli/export"
	"marcli/marc"
	"strings"
)

var fileName, search, fields, format string

func init() {
	flag.StringVar(&fileName, "file", "", "MARC file to process. Required.")
	flag.StringVar(&search, "match", "", "Only records that match the string passed, case insensitive.")
	flag.StringVar(&fields, "fields", "", "Comma delimited list of fields to output.")
	flag.StringVar(&format, "format", "mrk", "Output format. Accepted values: mrk, mrc, json, or solr.")
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
		// TODO: support filters in ToMrc exporter
		err = export.ToMrc(fileName, searchValue, filters)
	} else if format == "mrk" {
		err = export.ToMrk(fileName, searchValue, filters)
	} else if format == "json" {
		err = export.ToJson(fileName, searchValue, filters)
	} else {
		err = errors.New("Invalid format")
	}
	if err != nil {
		panic(err)
	}
}
