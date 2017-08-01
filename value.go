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

func (v SubFieldValue) String() string {
	return fmt.Sprintf("$%s%s", v.SubField, v.Value)
}

// Represents the entire value for a field.
// For example in:
//		=650  \0$aDiabetes$xComplications$zUnited States.
// Value will be:
// 		Value{
//			Tag: "650",
//			Ind1:" ",
//			Ind2: "0",
//			RawValue: "$aDiabetes$xComplications$zUnited States."
//			SubFieldValues (see SubFieldValue definition above)
//	}
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
		// notice that we skip the indicators because they are handled above
		// and valueStr[2] because that's a separator character
		value.RawValue = valueStr[3:]
	}

	if tag > "009" {
		value.SubFieldValues = NewSubFieldValues(valueStr)
	}
	return value
}

func NewSubFieldValues(valueStr string) []SubFieldValue {
	var values []SubFieldValue
	// valueStr comes with the indicators, we skip them:
	//   value[0] indicator 1
	// 	 value[0] indicator 2
	// 	 value[0] separator (ascii 31/0x1f)
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
