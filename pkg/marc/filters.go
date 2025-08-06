package marc

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

var ErrInvalidFieldString = errors.New("invalid field string (too short)")

// fieldsStr is a comma delimited string in the format NNNabc,NNNabc
// where NNN represents the MARC field to output and abc...z represents
// a set of subfields to include. If no subfields are indicated all
// subfields for the field are assummed.
// Example:
//
//	"700a" represents MARC field 700, subfield a.
//	"700ag" represents MARC field 700, subfields a and g.
//	"700" represents field 700 and all its subfields.
func NewFieldFilters(fieldsStr string) FieldFilters {
	if fieldsStr == "" {
		return FieldFilters{}
	}
	filters := FieldFilters{}
	for _, value := range strings.Split(fieldsStr, ",") {
		filter, err := NewFieldFilter(value)
		if err != nil {
			// TODO: handle error
			return FieldFilters{}
		}
		filters.Fields = append(filters.Fields, filter)
	}
	return filters
}

// fieldStr is a string in the format NNNabc
func NewFieldFilter(fieldStr string) (FieldFilter, error) {
	if len(fieldStr) < 3 {
		return FieldFilter{}, ErrInvalidFieldString
	}
	tag := fieldStr[:3]
	subfields := ""
	if len(fieldStr) > 3 {
		subfields = fieldStr[3:]
	}
	filter := FieldFilter{Tag: tag, Subfields: subfields}
	return filter, nil
}

func (filters FieldFilters) String() string {
	s := "Filters {\n"
	for _, field := range filters.Fields {
		if field.Subfields == "" {
			s += fmt.Sprintf("\tTag: %s\n", field.Tag)
		} else {
			s += fmt.Sprintf("\tTag: %s subfields: %s\n", field.Tag, field.Subfields)
		}
	}
	s += "}\n"
	return s
}

func (filters FieldFilters) IncludeField(name string) bool {
	for _, field := range filters.Fields {
		if field.Tag == name {
			return true
		}
	}
	return false
}

func (filters FieldFilters) IncludeLeader() bool {
	// return true if no fields specified: leader is part of MARC data
	return len(filters.Fields) == 0 || filters.IncludeField("LDR")
}
