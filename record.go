package main

import (
	"fmt"
)

// Record is a struct representing a MARC record. It has a Fields slice
// which contains both ControlFields and DataFields.
type Record struct {
	Data   []byte
	Fields []DataField
	Leader Leader
}

func (r Record) IsMatch(searchValue string) bool {
	return true
	// TODO
	// if searchValue == "" {
	// 	return true
	// }
	// for _, field := range r.Fields.All() {
	// 	if strings.Contains(strings.ToLower(field.RawValue), searchValue) {
	// 		return true
	// 	}
	// }
	// return false
}

func (r Record) ControlNum() string {
	for _, f := range r.Fields {
		if f.Tag == "001" {
			return f.Value
		}
	}
	return ""
}

func (r Record) String() string {
	return fmt.Sprintf("Leader: %s", r.Leader)
}
