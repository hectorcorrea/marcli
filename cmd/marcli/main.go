package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/hectorcorrea/marcli/pkg/marc"
)

var fileName, search, searchRegEx, searchFields, fields, exclude, format, hasFields, newLine string
var start, count int
var debug bool

func init() {
	flag.StringVar(&fileName, "file", "", "MARC file to process. Required.")
	flag.StringVar(&search, "match", "", "String that must be present in the content of the record, case insensitive.")
	flag.StringVar(&searchRegEx, "matchRegEx", "", "A regular expression to match the record.")
	flag.StringVar(&searchFields, "matchFields", "", "Comma delimited list of fields to search, used when match parameter is indicated, defaults to all fields.")
	flag.StringVar(&fields, "fields", "", "Comma delimited list of fields to output.")
	flag.StringVar(&exclude, "exclude", "", "Comma delimited list of fields to exclude from the output.")
	flag.StringVar(&format, "format", "mrk", "Output format. Accepted values: mrk, mrc, xml, json, solr, or count-only.")
	flag.IntVar(&start, "start", 1, "Number of first record to load.")
	flag.IntVar(&count, "count", -1, "Total number of records to load (-1 no limit).")
	flag.StringVar(&hasFields, "hasFields", "", "Comma delimited list of fields that must be present in the record.")
	flag.BoolVar(&debug, "debug", false, "When true it does not stop on errors.")
	flag.StringVar(&newLine, "newLine", "LF", "Character(s) to use to indicate new lines. Valid values LF or CRLF.")
	flag.Parse()
}

func main() {
	if fileName == "" {
		showSyntax()
		return
	}

	params := ProcessFileParams{
		filename:     fileName,
		format:       format,
		searchValue:  strings.ToLower(search),
		searchRegEx:  searchRegEx,
		searchFields: searchFieldsFromString(searchFields),
		filters:      marc.NewFieldFilters(fields),
		exclude:      marc.NewFieldFilters(exclude),
		start:        start,
		count:        count,
		hasFields:    marc.NewFieldFilters(hasFields),
		debug:        debug,
		newLine:      newLine,
	}

	if len(params.filters.Fields) > 0 && len(params.exclude.Fields) > 0 {
		panic("Cannot specify fields and exclude at the same time.")
	}

	if params.searchValue != "" && params.searchRegEx != "" {
		panic("Cannot specify match and matchRegEx at the same time.")
	}

	var err error
	if format == "mrk" || format == "count-only" {
		err = toMrk(params)
	} else if format == "mrc" {
		err = toMrc(params)
	} else if format == "json" {
		err = toJson(params)
	} else if format == "solr" {
		err = toSolr(params)
	} else if format == "xml" {
		err = toXML(params)
	} else {
		err = errors.New("invalid format")
	}
	if err != nil {
		panic(err)
	}
}

func showSyntax() {
	fmt.Println("marcli parameters:")
	fmt.Println()
	flag.PrintDefaults()
	fmt.Println()
	fmt.Printf(`
NOTES:
	The match parameter is used to filter records based on their content.
By default marcli searches in all the fields for each record, you can use
the matchFields parameter to limit the search to only certain fields (subfields
are not supported in matchFields, i.e. 245 is OK, 245a is not)

    The matchRegEx parameter can be used to filter records based on a regular expression
(e.g. '.*03-\d\d-06.*' to get records with dates from March 2006)

    The hasFields parameter is used to filter records based on the presence
of certain fields on the record (regardless of their value).

	You can only use the fields or exclude parameter, but not both.
`)
	fmt.Println()
	fmt.Println()
}

func searchFieldsFromString(searchFieldsString string) []string {
	values := []string{}
	for _, value := range strings.Split(searchFieldsString, ",") {
		if strings.TrimSpace(searchFieldsString) != "" {
			values = append(values, value)
		}
	}
	return values
}
