package main

import (
	"errors"
	"fmt"
	"strings"
)

type FieldFilters struct {
	Fields []FieldFilter
}

type FieldFilter struct {
	Tag       string
	Subfields string
}

// str is a comma delimited string in the format NNNabc,NNNabc
// where NNN represents the MARC field to output and abc...z represents
// a set of subfields to include. If no subfields are indicated all
// subfields for the field are assummed.
// Example:
//		"700a" represents MARC field 700, subfield a.
// 		"700ag" represents MARC field 700, subfields a and g.
//		"700" represents field 700 and all its subfields.
func NewFieldFilters(fieldsStr string) FieldFilters {
	if fieldsStr == "" {
		return FieldFilters{}
	}
	filters := FieldFilters{}
	for _, value := range strings.Split(fieldsStr, ",") {
		filters.addFilter(value)
	}
	return filters
}

func (filters FieldFilters) String() string {
	s := "Filters {\r\n"
	for _, field := range filters.Fields {
		if field.Subfields == "" {
			s += fmt.Sprintf("\tTag: %s\r\n", field.Tag)
		} else {
			s += fmt.Sprintf("\tTag: %s subfields: %s\r\n", field.Tag, field.Subfields)
		}
	}
	s += "}\r\n"
	return s
}

// fieldStr is a string in the format NNNabc
func (filters *FieldFilters) addFilter(fieldStr string) error {
	if len(fieldStr) < 3 {
		return errors.New("Invalid field string (too short)")
	}
	tag := fieldStr[0:3]
	subfields := ""
	if len(fieldStr) > 3 {
		subfields = fieldStr[3:]
	}
	filter := FieldFilter{Tag: tag, Subfields: subfields}
	filters.Fields = append(filters.Fields, filter)
	return nil
}

func (filters FieldFilters) Apply(values []Value) []Value {
	if len(filters.Fields) == 0 {
		return values
	}
	var filtered []Value
	for _, field := range filters.Fields {
		for _, value := range values {
			if value.Tag == field.Tag {
				filtered = append(filtered, value)
			}
		}
	}
	return filtered
}

func (filters FieldFilters) IncludeLeader() bool {
	if len(filters.Fields) == 0 {
		return true
	}
	for _, field := range filters.Fields {
		if field.Tag == "LDR" {
			return true
		}
	}
	return false
}