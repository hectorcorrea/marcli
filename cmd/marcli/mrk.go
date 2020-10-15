package main

import (
	"fmt"
	"io"
	"os"

	"github.com/hectorcorrea/marcli/pkg/marc"
)

func toMrk(params ProcessFileParams) error {
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

		if r.Contains(params.searchValue) && r.HasFields(params.hasFields) {
			str := ""
			if params.filters.IncludeLeader() {
				str += fmt.Sprintf("%s\r\n", r.Leader)
			}
			for _, field := range r.Filter(params.filters) {
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
