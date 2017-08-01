package main

import (
	"fmt"
	"strings"
)

type ConsoleProcessor struct {
	Filters FieldFilters
}

func (p ConsoleProcessor) Process(r Record) {
	if p.Filters.IncludeLeader() {
		fmt.Printf("%s (%d, %d, %d)\r\n", r.Leader, r.Pos, r.Leader.Length, r.Leader.DataOffset)
	}
	for _, value := range p.Filters.Apply(r.Values) {
		fmt.Printf("%s\r\n", value)
	}
	fmt.Printf("\r\n\r\n")
}

type ExtractProcessor struct {
	Filters FieldFilters
	Value   string // value to search
}

func (p ExtractProcessor) Process(r Record) {
	match := false
	for _, v := range r.Values {
		if strings.Contains(strings.ToLower(v.RawValue), p.Value) {
			match = true
			break
		}
	}

	if match {
		if p.Filters.IncludeLeader() {
			fmt.Printf("%s (%d, %d, %d)\r\n", r.Leader, r.Pos, r.Leader.Length, r.Leader.DataOffset)
		}
		for _, value := range p.Filters.Apply(r.Values) {
			fmt.Printf("%s\r\n", value)
		}
		fmt.Printf("\r\n\r\n")
	}
}
