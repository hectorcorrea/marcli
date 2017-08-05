package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ConsoleProcessor struct {
	Filters     FieldFilters
	SearchValue string
	Format      string
	outputCount int
}

func (p ConsoleProcessor) Header() {
	if p.Format == "json" {
		fmt.Printf("[\r\n")
	}
}

func (p ConsoleProcessor) Footer() {
	if p.Format == "json" {
		fmt.Printf("\r\n]\r\n")
	}
}

func (p *ConsoleProcessor) Process(r Record) {
	if !p.isMatch(r) {
		return
	}
	p.outputCount += 1
	if p.Format == "json" {
		p.outputJson(r)
	} else {
		p.outputMrk(r)
	}
}

func (p ConsoleProcessor) isMatch(r Record) bool {
	if p.SearchValue == "" {
		return true
	}
	for _, v := range r.Values {
		if strings.Contains(strings.ToLower(v.RawValue), p.SearchValue) {
			return true
		}
	}
	return false
}

func (p ConsoleProcessor) outputMrk(r Record) {
	if p.Filters.IncludeLeader() {
		fmt.Printf("%s (%d, %d, %d)\r\n", r.Leader, r.Pos, r.Leader.Length, r.Leader.DataOffset)
	}
	for _, value := range p.Filters.Apply(r.Values) {
		fmt.Printf("%s\r\n", value)
	}
	fmt.Printf("\r\n")
}

func (p ConsoleProcessor) outputJson(r Record) {
	if p.outputCount > 1 {
		fmt.Printf(", \r\n")
	}

	// Create a copy of the record but only with the
	// values indicated in the filters.
	rr := r
	rr.Values = p.Filters.Apply(r.Values)

	b, err := json.Marshal(rr)
	if err != nil {
		fmt.Printf("%s\r\n", err)
	}
	fmt.Printf("%s\r\n", b)
}
