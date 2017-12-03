package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ProcessorSolr struct {
	Filters     FieldFilters
	SearchValue string
}

type SolrDocument struct {
	Id              string   `json:"id"`
	Author          string   `json:"author,omitempty"`
	AuthorDate      string   `json:"authorDate,omitempty"`
	AuthorFuller    string   `json:"authorFuller,omitempty"`
	Title           string   `json:"title,omitempty"`
	Responsibility  string   `json:"responsibility,omitempty"`
	Publisher       string   `json:"publisher,omitempty"`
	Subjects        []string `json:"subjects,omitempty"`
	SubjectsForm    []string `json:"subjectsForm,omitempty"`
	SubjectsGeneral []string `json:"subjectsGeneral,omitempty"`
	SubjectsChrono  []string `json:"subjectsChrono,omitempty"`
	SubjectsGeo     []string `json:"subjectsGeo,omitempty"`
}

func NewSolrDocument(r Record) SolrDocument {
	doc := SolrDocument{}
	id := r.Fields.GetValue("001", "")
	if id == "" {
		id = "INVALID"
	}
	doc.Id = strings.TrimSpace(id)
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

	titleA := r.Fields.GetValue("245", "a")
	titleB := r.Fields.GetValue("245", "b")
	titleC := r.Fields.GetValue("245", "c")
	doc.Title = concat(titleA, titleB)
	doc.Responsibility = titleC

	doc.Publisher = r.Fields.GetValue("260", "a")
	doc.Subjects = subjects(r, "a")
	doc.SubjectsForm = subjects(r, "v")
	doc.SubjectsGeneral = subjects(r, "x")
	doc.SubjectsChrono = subjects(r, "y")
	doc.SubjectsGeo = subjects(r, "z")
	return doc
}

func (p ProcessorSolr) Header() {
	fmt.Printf("[\r\n")
}

func (p ProcessorSolr) Footer() {
	fmt.Printf("\r\n]\r\n")
}

func (p ProcessorSolr) ProcessRecord(f *MarcFile, r Record) {
	doc := NewSolrDocument(r)
	str, err := json.Marshal(doc)
	if err != nil {
		fmt.Printf("%s\r\n", err)
	}
	fmt.Printf("%s", str)
}

func (p ProcessorSolr) Separator() {
	fmt.Printf(", \r\n")
}

func subjects(r Record, subfield string) []string {
	var values []string
	for _, f_650 := range r.Fields.Get("650") {
		value := f_650.SubFieldValue(subfield)
		if value != "" {
			values = append(values, trimPeriod(value))
		}
	}
	return values
}
