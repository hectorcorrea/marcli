package marc

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrInvalidIndicators  = errors.New("invalid Indicators detected")
	ErrBadSubfieldsLength = errors.New("bad SubFields length")
)

// Field represents a field inside a MARC record. Notice that the
// field could be a "control" field (tag 001-009) or a "data" field
// (any other tag)
//
// For example in:
//
//	=650  \0$aDiabetes$xComplications$zUnited States.
//
// Field would be:
//
//		Field{
//			Tag: "650",
//			Value: ""
//			Indicator1: " ",
//			Indicator2: "0",
//			SubFields (see SubField definition below)
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
//
//	=650  \0$aDiabetes$xComplications$zUnited States.
//
// an example of SubFieldValue will be:
//
//	SubField{
//		Code: "a",
//		Value: "Diabetes"
//	}
type SubField struct {
	Code  string
	Value string
}

// MakeField creates a field object with the data received.
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
		return f, ErrInvalidIndicators
	}

	if len(data) < 4 { // Each data field contains at least one subfield code.
		return f, ErrBadSubfieldsLength
	}

	for _, sf := range bytes.Split(data[3:], []byte{st}) {
		if len(sf) > 1 {
			f.SubFields = append(f.SubFields, SubField{string(sf[0]), string(sf[1:])})
		}
	}
	return f, nil
}

// IsControlField returns true if the field is a control field (tag 001-009)
func (f Field) IsControlField() bool {
	return strings.HasPrefix(f.Tag, "00")
}

// Contains returns true if the field contains the passed string or matches the regex.
func (f Field) Contains(str string, regEx string) bool {
	if str != "" {
		return f.containsValue(str)
	} else {
		return f.containsRegEx(regEx)
	}
}

func (f Field) containsValue(str string) bool {
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

func (f Field) containsRegEx(regEx string) bool {
	re := regexp.MustCompile(regEx)

	if f.IsControlField() {
		matches := re.FindStringSubmatch(f.Value)
		// if matches != nil {
		// 	fmt.Printf("Control field match %s: %#v\n", f.Tag, matches)
		// }
		return matches != nil
	}

	for _, sub := range f.SubFields {
		matches := re.FindStringSubmatch(sub.Value)
		if matches != nil {
			// fmt.Printf("Field match %s: %#v\n", f.Tag, matches)
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
