package main

import (
	"fmt"
)

// Record is a struct representing a MARC record. It has a Fields slice
// which contains both ControlFields and DataFields.
type Record struct {
	Data   []byte
	Fields []Field
	Leader Leader
}

// Contains returns true if Record contains the value passed.
func (r Record) Contains(searchValue string) bool {
	if searchValue == "" {
		return true
	}
	for _, field := range r.Fields {
		if field.Contains(searchValue) {
			return true
		}
	}
	return false
}

// ControlNum returns the control number (tag 001) for the record.
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
