package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hectorcorrea/marcli/pkg/marc"
)

// TODO: Add support for JSONL (JSON line delimited) format that makes JSON
// easier to parse with Unix tools like grep, tail, and so on.
func toJson(params ProcessFileParams) error {
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
		if r.Contains(params.searchValue, params.searchRegEx, params.searchFields) && r.HasFields(params.hasFields) {
			if out > 0 {
				fmt.Printf(",%s", params.NewLine())
			} else {
				fmt.Printf("%s", params.NewLine())
			}
			b, err := json.Marshal(r.Filter(params.filters, params.exclude))
			if err != nil {
				fmt.Printf("%s%s", err, params.NewLine())
			}
			fmt.Printf("%s", b)
			if out++; out == count {
				break
			}
		}
	}
	fmt.Printf("%s]%s", params.NewLine(), params.NewLine())

	return marc.Err()
}
