package main

import (
	"fmt"
)

type SubFieldValue struct {
	SubField string
	Value    string
}

func (v SubFieldValue) String() string {
	return fmt.Sprintf("$%s%s", v.SubField, v.Value)
}

type Value struct {
	Tag            string
	Ind1           string
	Ind2           string
	RawValue       string // includes indicators and separator character
	SubFieldValues []SubFieldValue
}

func NewValue(tag, valueStr string) Value {
	value := Value{Tag: tag}

	if len(valueStr) >= 2 {
		value.Ind1 = string(valueStr[0])
		value.Ind2 = string(valueStr[1])
	}

	if len(valueStr) > 2 {
		value.RawValue = valueStr[3:]
	}

	if tag > "009" {
		value.SubFieldValues = NewFieldsFromString(valueStr)
	}
	return value
}

func (v Value) String() string {
	ind1 := formatIndicator(v.Ind1)
	ind2 := formatIndicator(v.Ind2)
	strValue := ""
	if len(v.SubFieldValues) > 0 {
		// use the subfield values
		for _, fv := range v.SubFieldValues {
			strValue += fmt.Sprintf("$%s%s", fv.SubField, fv.Value)
		}
	} else {
		// use the raw value
		strValue = v.RawValue
	}
	return fmt.Sprintf("=%s  %s%s%s", v.Tag, ind1, ind2, strValue)
}

func formatIndicator(value string) string {
	if value == " " {
		return "\\"
	}
	return value
}
