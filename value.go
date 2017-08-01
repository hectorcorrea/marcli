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

func (v Value) String() string {
	ind1 := "\\"
	ind2 := "\\"
	strValue := ""
	if len(v.RawValue) > 3 {
		ind1 = formatIndicator(v.RawValue[0])
		ind2 = formatIndicator(v.RawValue[1])
		strValue = v.RawValue[3:]
	}

	if len(v.SubFieldValues) > 0 {
		// use the subfield values rather than the raw value
		for _, fv := range v.SubFieldValues {
			strValue += fmt.Sprintf("%s%s", fv.SubField, fv.Value)
		}
	}
	return fmt.Sprintf("=%s  %s%s%s", v.Tag, ind1, ind2, strValue)
}

func formatIndicator(value byte) string {
	if value == ' ' {
		return "\\"
	}
	return string(value)
}
