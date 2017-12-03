package main

import (
	"encoding/json"
	"fmt"
)

type ConsoleProcessor struct {
	Filters     FieldFilters
	SearchValue string
	Format      string
}

func (p ConsoleProcessor) Header() {
	if p.Format == "json" {
		fmt.Printf("[\r\n")
	}
}

func (p ConsoleProcessor) Footer() {
	if p.Format == "json" {
		fmt.Printf("]\r\n")
	}
}

func (p ConsoleProcessor) ProcessRecord(f *MarcFile, r Record) {
	if p.Format == "json" {
		p.outputJson(r, f.Name)
	} else {
		p.outputMrk(r, f.Name)
	}
}

func (p ConsoleProcessor) Separator() {
	if p.Format == "json" {
		fmt.Printf(", \r\n")
	} else {
		fmt.Printf("\r\n")
	}
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
	filteredFields := p.Filters.Apply(r.Fields)
	for _, field := range filteredFields.All() {
		fmt.Printf("%s\r\n", field)
	}
}

func (p ConsoleProcessor) outputJson(r Record, filename string) {
	// TODO: Handle Leader, RecordInfo, and FileInfo fields
	filteredFields := p.Filters.Apply(r.Fields)
	b, err := json.Marshal(filteredFields.All())
	if err != nil {
		fmt.Printf("%s\r\n", err)
	}
	// fmt.Printf("{ \"record\": %s}\r\n", b)
	fmt.Printf("%s\r\n", b)
}
