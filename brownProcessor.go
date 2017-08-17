package main

import (
	"fmt"
)

type BrownProcessor struct {
}

type BrownRecord struct {
	Bib         string
	Title       string
	Callnumbers []string
}

func NewBrownRecord(r Record) BrownRecord {
	b := BrownRecord{}
	b.Bib = bib(r)
	b.Title = pad(r.SubValueFor("245", "a"))
	b.Callnumbers = callnumbers(r)
	return b
}

func (p BrownProcessor) Header() {
	fmt.Printf("BIB\tTitle\tCallnumber\r\n")
}

func (p BrownProcessor) Footer() {
}

func (p BrownProcessor) Process(f *MarcFile, r Record, count int) {
	b := NewBrownRecord(r)
	if len(b.Callnumbers) == 0 {
		// fmt.Printf("%s\t%s\t%s\r\n", b.Bib, b.Title, "--")
	} else {
		for _, callnumber := range b.Callnumbers {
			fmt.Printf("%s\t%s\t%s\r\n", b.Bib, b.Title, callnumber)
		}
	}
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

	if found, _ := r.ValueFor("090"); !found {
		// TODO: handle other 09X fields
		return numbers
	}

	f_090a := r.SubValueFor("090", "a")
	f_090b := r.SubValueFor("090", "b")
	f_090f := r.SubValueFor("090", "f") // 1-SIZE
	items := r.ValuesFor("945")
	if len(items) == 0 {
		// no items, use the bib call number
		// TODO: do we want this?
		number := concat(f_090f, f_090a, f_090b)
		numbers = append(numbers, number)
		return numbers
	}

	// get the call numbers from the items
	for _, f_945 := range items {
		base := concat(f_090f, f_090a, f_090b)
		f_945a := f_945.SubFieldValue("a")
		if f_945a != "" {
			// Annex Hay items
			base = concat(f_090f, f_945a, "")
		}
		c := f_945.SubFieldValue("c") // volume
		g := f_945.SubFieldValue("g") // copy
		if g == "1" {
			g = ""
		}
		number := concat(base, c, g)
		numbers = append(numbers, number)
	}
	return numbers
}

func pad(str string) string {
	padded := fmt.Sprintf("%-40s", str)
	if len(padded) > 40 {
		return padded[0:40]
	}
	return padded
}

func concat(a, b, c string) string {
	str := ""
	if a != "" {
		str += a
	}

	if b != "" {
		if len(str) > 0 {
			str += " "
		}
		str += b
	}

	if c != "" {
		if len(str) > 0 {
			str += " "
		}
		str += c
	}
	return str
}
