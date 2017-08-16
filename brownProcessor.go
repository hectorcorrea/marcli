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
	b.Callnumbers = callnumbers(r)
	return b
}

func bib(r Record) string {
	bib := r.SubValueFor("907", "a")
	if bib != "" {
		bib = bib[1:(len(bib) - 1)]
	}
	return bib
}

func callnumbers(r Record) []string {
	var numbers []string
	f_090a := r.SubValueFor("090", "a")
	f_090b := r.SubValueFor("090", "b")
	f_090f := r.SubValueFor("090", "f")
	base := fmt.Sprintf("%s %s %s", f_090f, f_090a, f_090b)
	items := r.ValuesFor("945")
	if len(items) == 0 {
		numbers = append(numbers, base)
	} else {
		for _, f_945 := range items {
			// f_945a := f_945.SubFieldValue("a")
			// if f_945a != "" {
			// 	base = fmt.Sprintf("%s %s", f_090f, f_945a)
			// }
			c := f_945.SubFieldValue("c")
			g := f_945.SubFieldValue("g")
			number := fmt.Sprintf("%s %s %s", base, c, g)
			numbers = append(numbers, number)
		}
	}
	return numbers
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
