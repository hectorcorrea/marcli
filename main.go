package main

import (
	"flag"
	"fmt"
	"strings"
)

var fileName, extract, output string

func init() {
	flag.StringVar(&fileName, "f", "", "MARC file to process. Required.")
	flag.StringVar(&extract, "x", "", "Extract record where string is found (case insensitive).")
	flag.StringVar(&output, "o", "all", "Fields to output (comma delimited).")
	flag.Parse()
}

func main() {
	if fileName == "" {
		fmt.Printf("marcli syntax:\r\n")
		flag.PrintDefaults()
		return
	}

	file, err := NewMarcFile(fileName)
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
