package main

import (
	"strings"

	"github.com/hectorcorrea/marcli/pkg/marc"
)

type ProcessFileParams struct {
	filename     string
	searchValue  string
	searchRegEx  string
	searchFields []string
	format       string
	filters      marc.FieldFilters
	exclude      marc.FieldFilters
	start        int
	count        int
	hasFields    marc.FieldFilters
	debug        bool
	newLine      string
}

func (p ProcessFileParams) HasFilters() bool {
	return len(p.filters.Fields) > 0 || len(p.exclude.Fields) > 0
}

func (p ProcessFileParams) NewLine() string {
	if strings.ToUpper(newLine) == "CRLF" {
		// Windows style
		return "\r\n"
	} else {
		// Unix style
		return "\n"
	}
}
