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

func (p BrownProcessor) Header() {
	fmt.Printf("BIB\tTitle\tCallnumber\r\n")
}

func (p BrownProcessor) Footer() {
}

func (p BrownProcessor) Process(f *MarcFile, r Record, count int) {
	p.outputCount = count
	p.outputTSV(r, f.Name)
}

func (p BrownProcessor) outputTSV(r Record, filename string) {
	for _, value := range p.Filters.Apply(r.Values) {
		fmt.Printf("BROWN: %s\r\n", value)
	}
	fmt.Printf("\r\n")
}
