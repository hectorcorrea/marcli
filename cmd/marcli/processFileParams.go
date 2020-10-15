package main

import (
	"github.com/hectorcorrea/marcli/pkg/marc"
)

type ProcessFileParams struct {
	filename    string
	searchValue string
	filters     marc.FieldFilters
	start       int
	count       int
	hasFields   marc.FieldFilters
}
