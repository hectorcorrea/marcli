package main

import (
	"fmt"
	"io"
	"os"
)

func mrcProcessor(filename string, searchValue string, filters FieldFilters) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	marc := NewMarcFile(file)
	for marc.Scan() {
		r, err := marc.Record()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if r.Contains(searchValue) {
			fmt.Printf("%s", r.Raw())
		}
	}
	return marc.Err()
}
