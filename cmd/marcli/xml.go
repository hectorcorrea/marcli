package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/hectorcorrea/marcli/pkg/marc"
)

type controlField struct {
	Tag   string `xml:"tag,attr"`
	Value string `xml:",chardata"`
}

type subfield struct {
	Code  string `xml:"code,attr"`
	Value string `xml:",chardata"`
}
type dataField struct {
	Tag       string     `xml:"tag,attr"`
	Ind1      string     `xml:"ind1,attr"`
	Ind2      string     `xml:"ind2,attr"`
	Subfields []subfield `xml:"subfield"`
}

type xmlRecord struct {
	XMLName       xml.Name       `xml:"record"`
	Leader        string         `xml:"leader"`
	ControlFields []controlField `xml:"controlfield"`
	DataFields    []dataField    `xml:"datafield"`
}

const xmlProlog = `<?xml version="1.0" encoding="UTF-8"?>`
const xmlRootBegin = `<collection xmlns="http://www.loc.gov/MARC21/slim" xmlns:marc="http://www.loc.gov/MARC21/slim">`
const xmlRootEnd = `</collection>`

func toXML(params ProcessFileParams) error {
	if count == 0 {
		return nil
	}

	file, err := os.Open(params.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Printf("%s\n%s\n", xmlProlog, xmlRootBegin)

	var i, out int
	marc := marc.NewMarcFile(file)
	for marc.Scan() {

		r, err := marc.Record()
		if err == io.EOF {
			break
		}

		if err != nil {
			printError(r, "PARSE ERROR", err)
			if params.debug {
				continue
			}
			return err
		}

		if i++; i < start {
			continue
		}

		if r.Contains(params.searchValue) && r.HasFields(params.hasFields) {
			str, err := recordToXML(r, params.debug)
			if err != nil {
				if params.debug {
					printError(r, "XML PARSE ERROR", err)
					continue
				}
				panic(err)
			}
			fmt.Printf("%s\r\n", str)
			if out++; out == count {
				break
			}
		}
	}
	fmt.Printf("%s\n", xmlRootEnd)

	return marc.Err()
}

func recordToXML(r marc.Record, debug bool) (string, error) {
	x := xmlRecord{
		Leader: r.Leader.Raw(),
	}

	for _, f := range r.Fields {
		if f.IsControlField() {
			x.ControlFields = append(x.ControlFields, controlField{Tag: f.Tag, Value: f.Value})
		} else {
			df := dataField{Tag: f.Tag, Ind1: f.Indicator1, Ind2: f.Indicator2}
			for _, s := range f.SubFields {
				df.Subfields = append(df.Subfields, subfield{Code: s.Code, Value: s.Value})
			}
			x.DataFields = append(x.DataFields, df)
		}
	}

	indent := ""
	if debug {
		indent = " "
	}
	b, err := xml.MarshalIndent(x, indent, indent)
	return string(b), err
}

func printError(r marc.Record, errType string, err error) {
	str := "== RECORD WITH ERROR STARTS HERE\n"
	str += fmt.Sprintf("%s:\n%s\n", errType, err.Error())
	str += r.DebugString() + "\n"
	str += "== RECORD WITH ERROR ENDS HERE\n\n"
	fmt.Print(str)
}
