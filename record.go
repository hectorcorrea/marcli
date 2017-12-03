package main

import (
	"strings"
)

type Record struct {
	Leader    Leader
	Directory []DirEntry
	Fields    Fields
	Pos       int
}

func (r Record) IsMatch(searchValue string) bool {
	if searchValue == "" {
		return true
	}
	for _, field := range r.Fields.All() {
		if strings.Contains(strings.ToLower(field.RawValue), searchValue) {
			return true
		}
	}
	return false
}
