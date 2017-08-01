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
	RawValue       string // includes indicators and separator characters
	SubFieldValues []SubFieldValue
}

func (v Value) String() string {
	ind1 := formatIndicator(v.RawValue[0])
	ind2 := formatIndicator(v.RawValue[1])
	return fmt.Sprintf("=%s  %s%s%s", v.Tag, ind1, ind2, v.RawValue[3:])
}

func formatIndicator(value byte) string {
	if value == ' ' {
		return "\\"
	}
	return string(value)
}
