package main

import (
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
}

func (p ProcessFileParams) HasFilters() bool {
	return len(p.filters.Fields) > 0 || len(p.exclude.Fields) > 0
}
