package main

import (
	"fmt"
)

type BrownProcessor struct {
}

type BrownRecord struct {
	Bib   string
	Title string
	Items []BrownItem
}

type BrownItem struct {
	Callnumber string
	Barcode    string
}

func NewBrownRecord(r Record) BrownRecord {
	b := BrownRecord{}
	b.Bib = bib(r)
	b.Title = pad(r.SubValueFor("245", "a"))
	b.Items = items(r)
	return b
}

func (p BrownProcessor) Header() {
	fmt.Printf("BIB\tTitle\tCallnumber\tBarcode\r\n")
}

func (p BrownProcessor) Footer() {
}

func (p BrownProcessor) Process(f *MarcFile, r Record, count int) {
	b := NewBrownRecord(r)
	if len(b.Items) == 0 {
		// fmt.Printf("%s\t%s\t%s\r\n", b.Bib, b.Title, "--")
	} else {
		for _, item := range b.Items {
			fmt.Printf("%s\t%s\t%s\t%s\r\n", b.Bib, b.Title, item.Callnumber, item.Barcode)
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

func baseCallNumber(r Record) (bool, Value) {
	// 090 ab            LC CALL NO(c)
	if found, value := r.ValueFor("090"); found {
		return true, value
	}

	// 091 ab            HARRIS CALL NO(e)
	if found, value := r.ValueFor("091"); found {
		return true, value
	}

	// 092 ab            JCB CALL NO(f)
	if found, value := r.ValueFor("092"); found {
		return true, value
	}

	// 096 ab           SUDOCS CALL NO(v)
	if found, value := r.ValueFor("096"); found {
		return true, value
	}

	// 099 ab            OTHER BROWN CALL (l)
	if found, value := r.ValueFor("099"); found {
		return true, value
	}

	return false, Value{}
}

func items(r Record) []BrownItem {
	var items []BrownItem

	marcItems := r.ValuesFor("945")
	if len(marcItems) == 0 {
		return items
	}

	// Base call number from the 09X field
	found, f_090 := baseCallNumber(r)
	if !found {
		return items
	}

	f_090a := f_090.SubFieldValue("a")
	f_090b := f_090.SubFieldValue("b")
	f_090f := f_090.SubFieldValue("f") // 1-SIZE

	// get the call numbers from the items
	for _, f_945 := range marcItems {
		barcode := f_945.SubFieldValue("i")
		base := concat(f_090f, f_090a, f_090b)
		f_945a := f_945.SubFieldValue("a")
		f_945b := f_945.SubFieldValue("b")
		if f_945a != "" {
			// use the values in the item record
			base = concat(f_090f, f_945a, f_945b)
		}
		c := f_945.SubFieldValue("c") // volume
		g := f_945.SubFieldValue("g") // copy
		if g == "1" {
			g = ""
		} else if g > "1" {
			g = "c. " + g
		}
		number := concat(base, c, g)
		item := BrownItem{Callnumber: number, Barcode: barcode}
		items = append(items, item)
	}
	return items
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
