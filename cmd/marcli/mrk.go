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
			str := "== RECORD WITH ERROR STARTS HERE\n"
			str += "ERROR:\n" + err.Error() + "\n"
			str += r.DebugString() + "\n"
			str += "== RECORD WITH ERROR ENDS HERE\n\n"
			fmt.Print(str)
			if params.debug {
				continue
			}
			return err
		}

		if i++; i < start {
			continue
		}

		if r.Contains(params.searchValue, params.searchFields) && r.HasFields(params.hasFields) {
			str := ""
			if params.filters.IncludeLeader() {
				str += fmt.Sprintf("%s\r\n", r.Leader)
			}
			for _, field := range r.Filter(params.filters, params.exclude) {
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
