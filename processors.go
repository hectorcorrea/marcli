package main

import (
	"fmt"
	"strings"
)

type ConsoleProcessor struct {
	Filters     FieldFilters
	SearchValue string
}

func (p ConsoleProcessor) Process(r Record) {
	match := false
	if p.SearchValue == "" {
		match = true
	} else {
		for _, v := range r.Values {
			if strings.Contains(strings.ToLower(v.RawValue), p.SearchValue) {
				match = true
				break
			}
		}
	}

	if !match {
		return
	}

	if p.Filters.IncludeLeader() {
		fmt.Printf("%s (%d, %d, %d)\r\n", r.Leader, r.Pos, r.Leader.Length, r.Leader.DataOffset)
	}
	for _, value := range p.Filters.Apply(r.Values) {
		fmt.Printf("%s\r\n", value)
	}
	fmt.Printf("\r\n")
}
