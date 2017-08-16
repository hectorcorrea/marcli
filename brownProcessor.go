package main

import (
	"fmt"
)

type BrownProcessor struct {
	Filters     FieldFilters
	SearchValue string // not used
	Format      string // not used
	outputCount int
}

type BrownRecord struct {
	Bib         string
	Title       string
	Callnumbers []string
}

func NewBrownRecord(r Record) BrownRecord {
	b := BrownRecord{}
	b.Bib = bib(r)
	b.Title = r.SubValueFor("245", "a")
	// b.Callnumbers = []
	// b.Callnumbers = append(b.Callnumbers, "bbb")
	return b
}

func bib(r Record) string {
	bib := r.SubValueFor("907", "a")
	if bib != "" {
		bib = bib[1:(len(bib) - 1)]
	}
	return bib
}

func (p BrownProcessor) Header() {
	fmt.Printf("BIB\tTitle\tCallnumber\r\n")
}

func (p BrownProcessor) Footer() {
}

func (p BrownProcessor) Process(f *MarcFile, r Record, count int) {
	p.outputTSV(r)
}

func (p BrownProcessor) outputTSV(r Record) {
	b := NewBrownRecord(r)
	if len(b.Callnumbers) == 0 {
		fmt.Printf("%s\t%s\t%s\r\n", b.Bib, b.Title, "--")
	} else {
		for _, callnumber := range b.Callnumbers {
			fmt.Printf("%s\t%s\t%s\r\n", b.Bib, b.Title, callnumber)
		}
	}
}
