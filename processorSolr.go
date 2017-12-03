package main

import (
	"encoding/json"
	"fmt"
)

type ProcessorSolr struct {
	Filters     FieldFilters
	SearchValue string
}

type SolrDocument struct {
	Id           string   `json:"id"`
	Author       string   `json:"author,omitempty"`
	AuthorDate   string   `json:"authorDate,omitempty"`
	AuthorFuller string   `json:"authorFuller,omitempty"`
	Title        string   `json:"title,omitempty"`
	Publisher    string   `json:"publisher,omitempty"`
	Subjects     []string `json:"subjects,omitempty"`
}

func NewSolrDocument(r Record) SolrDocument {
	doc := SolrDocument{}
	id := r.Fields.GetValue("001", "")
	if id == "" {
		id = "INVALID"
	}
	doc.Id = id
	author := r.Fields.GetValue("100", "a")
	if author != "" {
		doc.Author = author
		doc.AuthorDate = r.Fields.GetValue("100", "d")
		doc.AuthorFuller = r.Fields.GetValue("100", "a")
	} else {
		doc.Author = r.Fields.GetValue("110", "a")
		doc.AuthorDate = ""
		doc.AuthorFuller = ""
	}

	doc.Title = r.Fields.GetValue("245", "a")
	doc.Publisher = r.Fields.GetValue("260", "a")
	doc.Subjects = subjects(r)
	return doc
}

func (p ProcessorSolr) Header() {
	fmt.Printf("[\r\n")
}

func (p ProcessorSolr) Footer() {
	fmt.Printf("]\r\n")
}

func (p ProcessorSolr) ProcessRecord(f *MarcFile, r Record) {
	doc := NewSolrDocument(r)
	str, err := json.Marshal(doc)
	if err != nil {
		fmt.Printf("%s\r\n", err)
	}
	fmt.Printf("%s\r\n", str)
}

func (p ProcessorSolr) Separator() {
	fmt.Printf(", \r\n")
}

func subjects(r Record) []string {
	var s []string
	for _, f_650 := range r.Fields.Get("650") {
		f_650a := f_650.SubFieldValue("a")
		s = append(s, f_650a)
	}
	return s
}
