package main

import (
	"fmt"
	"strings"
)

// Represents a single subfield value.
// For example in:
//		=650  \0$aDiabetes$xComplications$zUnited States.
// an example of SubFieldValue will be:
// 		SubFieldValue{
//			SubField: "a",
//			Value: "Diabetes"
//		}
type SubFieldValue struct {
	SubField string
	Value    string
}

// Represents the entire value for a field.
// For example in:
//		=650  \0$aDiabetes$xComplications$zUnited States.
// Field would be:
// 		Field{
//			Tag: "650",
//			Ind1:" ",
//			Ind2: "0",
//			RawValue: "$aDiabetes$xComplications$zUnited States."
//			SubFields (see SubFieldValue definition above)
//	}
type Field struct {
	Tag       string
	Ind1      string
	Ind2      string
	RawValue  string // includes indicators and separator character
	SubFields []SubFieldValue
}

type Fields struct {
	fields []Field
}

func (v SubFieldValue) String() string {
	return fmt.Sprintf("$%s%s", v.SubField, v.Value)
}

func (f Fields) All() []Field {
	return f.fields
}

func (f *Fields) Add(field Field) {
	f.fields = append(f.fields, field)
}

func NewField(tag, valueStr string) Field {
	value := Field{Tag: tag}
	if tag <= "008" {
		// Control fields (001-008) don't have indicators or subfields
		// so we just get the value as-is.
		value.RawValue = valueStr
		return value
	}

	// Process the indicators and subfields
	if len(valueStr) >= 2 {
		value.Ind1 = string(valueStr[0])
		value.Ind2 = string(valueStr[1])
	}
	if len(valueStr) > 2 {
		// notice that we skip the indicators [0] and [1] because they are handled
		// above and valueStr[2] because that's a separator character
		value.RawValue = valueStr[3:]
	}
	value.SubFields = NewSubFieldValues(valueStr)
	return value
}

func NewSubFieldValues(valueStr string) []SubFieldValue {
	var values []SubFieldValue
	// valueStr comes with the indicators, we skip them:
	//   valueStr[0] indicator 1
	// 	 valueStr[1] indicator 2
	// 	 valueStr[2] separator (ascii 31/0x1f)
	separator := 0x1f
	tokens := strings.Split(valueStr[3:], string(separator))
	for _, token := range tokens {
		value := SubFieldValue{
			SubField: string(token[0]),
			Value:    token[1:],
		}
		values = append(values, value)
	}
	return values
}

func (f Field) String() string {
	ind1 := formatIndicator(f.Ind1)
	ind2 := formatIndicator(f.Ind2)
	strValue := ""
	if len(f.SubFields) > 0 {
		// use the subfield values
		for _, s := range f.SubFields {
			strValue += fmt.Sprintf("$%s%s", s.SubField, s.Value)
		}
	} else {
		// use the raw value
		strValue = f.RawValue
	}
	return fmt.Sprintf("=%s  %s%s%s", f.Tag, ind1, ind2, strValue)
}

func (f Field) SubFieldValue(subfield string) string {
	for _, s := range f.SubFields {
		if s.SubField == subfield {
			return s.Value
		}
	}
	return ""
}

// For a given value, extract the subfield values in the string
// indicated. "subfields" is a plain string, like "abu", to
// indicate subfields a, b, and u.
func (f Field) SubFieldValues(subfields string) []SubFieldValue {
	var values []SubFieldValue
	for _, sub := range f.SubFields {
		if strings.Contains(subfields, sub.SubField) {
			value := SubFieldValue{
				SubField: sub.SubField,
				Value:    sub.Value,
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

func (f Fields) Get(tag string) []Field {
	var fields []Field
	for _, field := range f.fields {
		if field.Tag == tag {
			fields = append(fields, field)
		}
	}
	return fields
}

func (f Fields) GetOne(tag string) (bool, Field) {
	fields := f.Get(tag)
	if len(fields) == 0 {
		return false, Field{}
	}
	return true, fields[0]
}

func (f Fields) GetValue(tag string, subfield string) string {
	value := ""
	found, field := f.GetOne(tag)
	if found {
		if subfield == "" {
			value = field.RawValue
		} else {
			value = field.SubFieldValue(subfield)
		}
	}
	return value
}

func (f Fields) GetValues(tag string, subfield string) []string {
	var values []string
	for _, field := range f.Get(tag) {
		var value string
		if subfield == "" {
			value = field.RawValue
		} else {
			value = field.SubFieldValue(subfield)
		}
		if value != "" {
			values = append(values, value)
		}
	}
	return values
}
