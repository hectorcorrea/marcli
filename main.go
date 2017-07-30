package main

import (
	"flag"
	"fmt"
	"strings"
)

var fileName, extract string

func init() {
	flag.StringVar(&fileName, "f", "", "MARC file to process. Required.")
	flag.StringVar(&extract, "x", "", "Extract record where value is found.")
	flag.Parse()
}

func main() {
	if fileName == "" {
		flag.PrintDefaults()
		return
	}

	file, err := NewFile(fileName)
	if err != nil {
		panic(err)
	}

	var processor RecordProcessor
	if extract != "" {
		processor = ExtractProcessor{Value: strings.ToLower(extract)}
	} else {
		processor = ConsoleProcessor{}
	}

	err = file.ReadAll(processor)
	if err != nil {
		panic(err)
	}
}

type ConsoleProcessor struct {
}

func (p ConsoleProcessor) Process(r Record) {
	fmt.Printf("=LDR  %s (%d, %d, %d)\n", r.Leader, r.Pos, r.Leader.Length, r.Leader.DataOffset)
	for _, v := range r.Values {
		fmt.Printf("=%s  %s\r\n", v.Tag, v.Value)
	}
	fmt.Printf("\r\n\r\n")
}

type ExtractProcessor struct {
	Value string // value to search
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
			fmt.Printf("=%s  %s\r\n", v.Tag, v.Value)
		}
		fmt.Printf("\r\n\r\n")
	}
}
