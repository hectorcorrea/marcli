package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// Field represents a field inside a MARC record. Notice that the
// field could be a "control" field (tag 001-009) or a "data" field
// (any other tag)
//
// For example in:
//		=650  \0$aDiabetes$xComplications$zUnited States.
// Field would be:
// 		Field{
//			Tag: "650",
//			Value: ""
//			Indicator1: " ",
//			Indicator2: "0",
//			SubFields (see SubField definition above)
//	}
type Field struct {
	Tag        string     // for both Control and Data fields
	Value      string     // for Control fields
	Indicator1 string     // for Data fields
	Indicator2 string     // for Data fields
	SubFields  []SubField // for Data fields
}

// SubField contains a Code and a Value.
// For example in:
//		=650  \0$aDiabetes$xComplications$zUnited States.
// an example of SubFieldValue will be:
// 		SubField{
//			Code: "a",
//			Value: "Diabetes"
//		}
type SubField struct {
	Code  string
	Value string
}

// MakeField creates a field objet with the data received.
func MakeField(tag string, data []byte) (Field, error) {
	f := Field{}
	f.Tag = tag

	// It's a control field
	if strings.HasPrefix(tag, "00") {
		f.Value = string(data)
		return f, nil
	}

	if len(data) > 2 {
		f.Indicator1 = string(data[0])
		f.Indicator2 = string(data[1])
	} else {
		return f, errors.New("Invalid Indicators detected")
	}

	for _, sf := range bytes.Split(data[3:], []byte{st}) {
		if len(sf) > 0 {
			f.SubFields = append(f.SubFields, SubField{string(sf[0]), string(sf[1:])})
		} else {
			return f, errors.New("Extraneous field terminator")
		}
	}
	return f, nil
}

// IsControlField returns true if the field is a control field (tag 001-009)
func (f Field) IsControlField() bool {
	return strings.HasPrefix(f.Tag, "00")
}

// Contains returns true if the field contains the passed string.
func (f Field) Contains(str string) bool {
	str = strings.ToLower(str)
	if f.IsControlField() {
		return strings.Contains(strings.ToLower(f.Value), str)
	}

	for _, sub := range f.SubFields {
		if strings.Contains(strings.ToLower(sub.Value), str) {
			return true
		}
	}
	return false
}

func (f Field) String() string {
	if f.IsControlField() {
		return fmt.Sprintf("=%s  %s", f.Tag, f.Value)
	}
	str := fmt.Sprintf("=%s  %s%s", f.Tag, formatIndicator(f.Indicator1), formatIndicator(f.Indicator2))
	for _, sub := range f.SubFields {
		str += fmt.Sprintf("$%s%s", sub.Code, sub.Value)
	}
	return str
}

// GetSubFields returns an array of subfields that match the set of subfields
// indicated in the filter string. "filter" is a plain string, like "abu", to
// indicate what subfields are to be returned.
func (f Field) GetSubFields(filter string) []SubField {
	values := []SubField{}
	for _, sub := range f.SubFields {
		if strings.Contains(filter, sub.Code) {
			value := SubField{
				Code:  sub.Code,
				Value: sub.Value,
			}
			values = append(values, value)
		}
	}
	return values
}

func formatIndicator(value string) string {
	if value == " " {
		return "\\"
	}
	return value
}
