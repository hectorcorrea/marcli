package main

import (
	"errors"
	"fmt"
	"hectorcorrea/marcli/pkg/marc"
	"io"
	"os"
)

func toMrc(filename string, searchValue string, filters marc.FieldFilters, start int, count int) error {
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
			fmt.Printf("%s", r.Raw())
			if out++; out == count {
				break
			}
		}
	}
	return marc.Err()
}
