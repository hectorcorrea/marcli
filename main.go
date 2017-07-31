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
	flag.StringVar(&output, "o", "", "Comma delimited list of fields to output.")
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

	fields := strToFields(output)

	var processor RecordProcessor
	if extract != "" {
		processor = ExtractProcessor{
			Fields: fields,
			Value:  strings.ToLower(extract),
		}
	} else {
		processor = ConsoleProcessor{
			Fields: fields,
		}
	}

	err = file.ReadAll(processor)
	if err != nil {
		panic(err)
	}
}

func strToFields(str string) []string {
	if str == "" {
		return []string{}
	}
	values := strings.Split(str, ",")
	return values
}
