package main

import (
	"encoding/json"
	"fmt"
)

type ProcessorSolr struct {
	Filters     FieldFilters
	SearchValue string
	outputCount int
}

type SolrDocument struct {
	Id           string   // MARC 001
	Author       string   // MARC 100a
	AuthorDate   string   // MARC 100d
	AuthorFuller string   // MARC 100q
	Title        string   // MARC 245
	Publisher    string   // MARC 260
	Subjects     []string // MARC 650
}

func NewSolrDocument(r Record) (SolrDocument, bool) {
	doc := SolrDocument{}
	doc.Id = r.Fields.GetValue("001", "")
	doc.Author = r.Fields.GetValue("100", "a")
	doc.AuthorDate = r.Fields.GetValue("100", "d")
	doc.AuthorFuller = r.Fields.GetValue("100", "a")
	doc.Title = r.Fields.GetValue("245", "a")
	doc.Publisher = r.Fields.GetValue("260", "a")
	doc.Subjects = subjects(r)
	empty := (doc.Id == "" || doc.Author == "" || doc.Title == "")
	return doc, empty
}

func (p ProcessorSolr) Header() {
	fmt.Printf("[\r\n")
}

func (p ProcessorSolr) Footer() {
	fmt.Printf("]\r\n")
}

func (p ProcessorSolr) Process(f *MarcFile, r Record, count int) {
	doc, empty := NewSolrDocument(r)
	if empty {
		return
	}
	if count > 1 {
		fmt.Printf(", \r\n")
	}
	str, err := json.Marshal(doc)
	if err != nil {
		fmt.Printf("%s\r\n", err)
	}
	fmt.Printf("%s\r\n", str)
}

func subjects(r Record) []string {
	var s []string
	for _, f_650 := range r.Fields.Get("650") {
		f_650a := f_650.SubFieldValue("a")
		s = append(s, f_650a)
	}
	return s
}
