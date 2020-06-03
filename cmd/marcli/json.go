package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"hectorcorrea/marcli/pkg/marc"
	"io"
	"os"
)

// TODO: Add support for JSONL (JSON line delimited) format that makes JSON
// easier to parse with Unix tools like grep, tail, and so on.
func toJson(filename string, searchValue string, filters marc.FieldFilters, start int, count int) error {
	if len(filters.Fields) > 0 {
		return errors.New("filters not supported for this format")
	}

	if count == 0 {
		return nil
	}

	file, err := os.Open(filename)
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
		if r.Contains(searchValue) {
			if out > 0 {
				fmt.Printf(",\r\n")
			} else {
				fmt.Printf("\r\n")
			}
			b, err := json.Marshal(r.Filter(filters))
			if err != nil {
				fmt.Printf("%s\r\n", err)
			}
			// fmt.Printf("{ \"record\": %s}\r\n", b)
			fmt.Printf("%s", b)
			if out++; out == count {
				break
			}
		}
	}
	fmt.Printf("\r\n]\r\n")

	return marc.Err()
}
