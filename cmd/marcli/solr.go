package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hectorcorrea/marcli/pkg/marc"
)

type SolrDocument struct {
	Id              string   `json:"id"`
	Author          string   `json:"author_txt_en,omitempty"`
	AuthorDate      string   `json:"author_date_s,omitempty"`
	AuthorFuller    string   `json:"author_fuller_txt_en,omitempty"`
	AuthorsOther    []string `json:"authors_other_txts_en,omitempty"`
	Title           string   `json:"title_txt_en,omitempty"`
	Responsibility  string   `json:"responsibility_txt_en,omitempty"`
	PublisherPlace  string   `json:"publisher_place_s,omitempty"`
	PublisherName   string   `json:"publisher_name_s,omitempty"`
	PublisherDate   string   `json:"publisher_date_s,omitempty"`
	Urls            []string `json:"urls_ss,omitempty"`
	Subjects        []string `json:"subjects_ss,omitempty"`
	SubjectsForm    []string `json:"subjects_form_ss,omitempty"`
	SubjectsGeneral []string `json:"subjects_general_ss,omitempty"`
	SubjectsChrono  []string `json:"subjects_chrono_ss,omitempty"`
	SubjectsGeo     []string `json:"subjects_geo_ss,omitempty"`
}

func NewSolrDocument(r marc.Record) SolrDocument {
	doc := SolrDocument{}
	id := r.GetValue("001", "")
	if id == "" {
		id = "INVALID"
	}
	doc.Id = strings.TrimSpace(id)
	author := r.GetValue("100", "a")
	if author != "" {
		doc.Author = author
		doc.AuthorDate = r.GetValue("100", "d")
		doc.AuthorFuller = r.GetValue("100", "q")
	} else {
		doc.Author = r.GetValue("110", "a")
		doc.AuthorDate = ""
		doc.AuthorFuller = ""
	}
	doc.AuthorsOther = r.GetValues("700", "a")

	titleA := r.GetValue("245", "a")
	titleB := r.GetValue("245", "b")
	titleC := r.GetValue("245", "c")
	doc.Title = concat(titleA, titleB)
	doc.Responsibility = titleC

	doc.PublisherPlace = r.GetValue("260", "a")
	doc.PublisherName = r.GetValue("260", "b")
	doc.PublisherDate = r.GetValue("260", "c")
	doc.Urls = r.GetValues("856", "u")
	doc.Subjects = subjects(r, "a")
	doc.SubjectsForm = subjects(r, "v")
	doc.SubjectsGeneral = subjects(r, "x")
	doc.SubjectsChrono = subjects(r, "y")
	doc.SubjectsGeo = subjects(r, "z")
	return doc
}

func toSolr(params ProcessFileParams) error {
	if params.HasFilters() {
		return errors.New("filters not supported for this format")
	}

	if count == 0 {
		return nil
	}

	file, err := os.Open(params.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var i, out int
	marc := marc.NewMarcFile(file)

	fmt.Printf("[")
	for marc.Scan() {
		r, err := marc.Record()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if i++; i < start {
			continue
		}
		if r.Contains(params.searchValue, params.searchRegEx, params.searchFields) && r.HasFields(params.hasFields) {
			if out > 0 {
				fmt.Printf(",\r\n")
			} else {
				fmt.Printf("\r\n")
			}
			doc := NewSolrDocument(r)
			b, err := json.Marshal(doc)
			if err != nil {
				fmt.Printf("%s\r\n", err)
			}
			fmt.Printf("%s", b)
			if out++; out == count {
				break
			}
		}
	}
	fmt.Printf("\r\n]\r\n")

	return marc.Err()
}

func subjects(r marc.Record, subfield string) []string {
	var values []string
	for _, fieldValue := range r.GetValues("650", subfield) {
		values = append(values, trimPeriod(fieldValue))
	}
	return values
}

func concat(a, b string) string {
	return _concat(a, b, " ")
}

func concatTab(a, b string) string {
	return _concat(a, b, "\t")
}

func _concat(a, b, sep string) string {
	if a == "" && b == "" {
		return ""
	} else if a == "" && b != "" {
		return strings.TrimSpace(b)
	} else if a != "" && b == "" {
		return strings.TrimSpace(a)
	}
	return strings.TrimSpace(a) + sep + strings.TrimSpace(b)
}

func trimPeriod(s string) string {
	if s == "" || s == "." {
		return ""
	}
	if strings.HasSuffix(s, ".") {
		return strings.TrimSpace(s[:len(s)-1])
	}
	return s
}
