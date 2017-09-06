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
	for _, field := range r.Fields.All() {
		if strings.Contains(strings.ToLower(field.RawValue), p.SearchValue) {
			return true
		}
	}
	return false
}

func (p ConsoleProcessor) outputMrk(r Record, filename string) {
	str := ""
	if p.Filters.IncludeLeader() {
		str += fmt.Sprintf("%s\r\n", r.Leader)
	}
	if p.Filters.IncludeRecordInfo() {
		str += fmt.Sprintf("=RIN  pos=%d, length=%d, data offset=%d\r\n", r.Pos, r.Leader.Length, r.Leader.DataOffset)
	}
	if p.Filters.IncludeFileInfo() {
		str += fmt.Sprintf("=FIN  %s\r\n", filename)
	}
	filteredFields := p.Filters.Apply(r.Fields)
	for _, field := range filteredFields.All() {
		str += fmt.Sprintf("%s\r\n", field)
	}
	if str != "" {
		fmt.Printf("%s\r\n", str)
	}
}

func (p ConsoleProcessor) outputJson(r Record, filename string) {
	if p.outputCount > 1 {
		fmt.Printf(", \r\n")
	}

	// TODO: Handle Leader, RecordInfo, and FileInfo fields
	filteredFields := p.Filters.Apply(r.Fields)
	b, err := json.Marshal(filteredFields.All())
	if err != nil {
		fmt.Printf("%s\r\n", err)
	}
	fmt.Printf("{ \"record\": %s}\r\n", b)
}
