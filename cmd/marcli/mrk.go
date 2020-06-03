package main

import (
	"fmt"
	"hectorcorrea/marcli/pkg/marc"
	"io"
	"os"
)

func toMrk(filename string, searchValue string, filters marc.FieldFilters, start int, count int) error {
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
			str := ""
			if filters.IncludeLeader() {
				str += fmt.Sprintf("%s\r\n", r.Leader)
			}
			for _, field := range r.Filter(filters) {
				str += fmt.Sprintf("%s\r\n", field)
			}
			if str != "" {
				fmt.Printf("%s\r\n", str)
				if out++; out == count {
					break
				}
			}
		}
	}

	return marc.Err()
}
