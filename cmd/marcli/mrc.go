package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hectorcorrea/marcli/pkg/marc"
)

func toMrc(params ProcessFileParams) error {
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

		if r.Contains(params.searchValue, params.searchFields) && r.HasFields(params.hasFields) {
			fmt.Printf("%s", r.Raw())
			if out++; out == count {
				break
			}
		}
	}
	return marc.Err()
}
