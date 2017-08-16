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

func (p ConsoleProcessor) Process(f *MarcFile, r Record, count int) {
	if !p.isMatch(r) {
		return
	}
	p.outputCount = count
	if p.Format == "json" {
		p.outputJson(r, f.Name)
	} else {
		p.outputMrk(r, f.Name)
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

func (p ConsoleProcessor) outputMrk(r Record, filename string) {
	if p.Filters.IncludeLeader() {
		fmt.Printf("%s\r\n", r.Leader)
	}
	if p.Filters.IncludeRecordInfo() {
		fmt.Printf("=RIN  pos=%d, length=%d, data offset=%d\r\n", r.Pos, r.Leader.Length, r.Leader.DataOffset)
	}
	if p.Filters.IncludeFileInfo() {
		fmt.Printf("=FIN  %s\r\n", filename)
	}
	for _, value := range p.Filters.Apply(r.Values) {
		fmt.Printf("%s\r\n", value)
	}
	fmt.Printf("\r\n")
}

func (p ConsoleProcessor) outputJson(r Record, filename string) {
	if p.outputCount > 1 {
		fmt.Printf(", \r\n")
	}

	// TODO: Handle Leader, RecordInfo, and FileInfo fields
	output := p.Filters.Apply(r.Values)
	b, err := json.Marshal(output)
	if err != nil {
		fmt.Printf("%s\r\n", err)
	}
	fmt.Printf("{ \"record\": %s}\r\n", b)
}
