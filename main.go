package main

import (
	"flag"
	"fmt"
	"strings"
)

var fileName, extract, output string

func init() {
	flag.StringVar(&fileName, "f", "", "MARC file to process. Required.")
	flag.StringVar(&extract, "x", "", "Extract record where string is found.")
	flag.StringVar(&output, "o", "all", "Fields to output.")
	flag.Parse()
}

func main() {
	if fileName == "" {
		fmt.Printf("marcli syntax:\r\n")
		flag.PrintDefaults()
		return
	}

	file, err := NewFile(fileName)
	if err != nil {
		panic(err)
	}

	var processor RecordProcessor
	if extract != "" {
		processor = ExtractProcessor{
			Fields: output,
			Value:  strings.ToLower(extract),
		}
	} else {
		processor = ConsoleProcessor{
			Fields: output,
		}
	}

	err = file.ReadAll(processor)
	if err != nil {
		panic(err)
	}
}

type ConsoleProcessor struct {
	Fields string
}

func (p ConsoleProcessor) Process(r Record) {
	fmt.Printf("=LDR  %s (%d, %d, %d)\n", r.Leader, r.Pos, r.Leader.Length, r.Leader.DataOffset)
	for _, v := range r.Values {
		// TODO: handle multiple comma delimited fields
		if p.Fields == "all" || p.Fields == v.Tag {
			fmt.Printf("=%s  %s\r\n", v.Tag, v.Value)
		}
	}
	fmt.Printf("\r\n\r\n")
}

type ExtractProcessor struct {
	Fields string
	Value  string // value to search
}

func (p ExtractProcessor) Process(r Record) {
	match := false
	for _, v := range r.Values {
		if strings.Contains(strings.ToLower(v.Value), p.Value) {
			match = true
			break
		}
	}

	if match {
		fmt.Printf("=LDR  %s (%d, %d, %d)\n", r.Leader, r.Pos, r.Leader.Length, r.Leader.DataOffset)
		for _, v := range r.Values {
			// TODO: handle multiple comma delimited fields
			if p.Fields == "all" || p.Fields == v.Tag {
				fmt.Printf("=%s  %s\r\n", v.Tag, v.Value)
			}
		}
		fmt.Printf("\r\n\r\n")
	}
}
