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

// For a given list of values, it returns only those values
// match the filters. The filter is done by Tag and if
// available by Sub Field.
func (filters FieldFilters) Apply(values []Value) []Value {
	if len(filters.Fields) == 0 {
		return values
	}

	var filtered []Value
	for _, field := range filters.Fields {
		// Process all the values that match the tag
		// (there could be more than one)
		for _, value := range valuesForTag(values, field.Tag) {
			if len(field.Subfields) == 0 {
				// add the value as-is, no need to filter by subfield
				filtered = append(filtered, value)
			} else {
				//... filter the value by subfield
				newValue := value
				newValue.RawValue = ""
				newValue.SubFieldValues = subFieldValuesFromValue(value, field.Subfields)
				filtered = append(filtered, newValue)
			}
		}
	}
	return filtered
}

// For a given value, extract the subfield values in the string
// indicated. "subfields" is a plain string, like "abu", to
// indicate subfields a, b, and u.
func subFieldValuesFromValue(value Value, subfields string) []SubFieldValue {
	var newValues []SubFieldValue
	for _, sv := range value.SubFieldValues {
		if strings.Contains(subfields, sv.SubField) {
			yy := SubFieldValue{
				SubField: sv.SubField,
				Value:    sv.Value,
			}
			newValues = append(newValues, yy)
		}
	}
	return newValues
}

func valuesForTag(values []Value, tag string) []Value {
	var vv []Value
	for _, value := range values {
		if value.Tag == tag {
			vv = append(vv, value)
		}
	}
	return vv
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
	if len(filters.Fields) == 0 {
		// included by default because it is part of the MARC data
		return true
	}
	return filters.IncludeField("LDR")
}

func (filters FieldFilters) IncludeFileInfo() bool {
	return filters.IncludeField("FIN")
}

func (filters FieldFilters) IncludeRecordInfo() bool {
	return filters.IncludeField("RIN")
}
