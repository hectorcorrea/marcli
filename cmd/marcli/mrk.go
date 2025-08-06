package main

import (
	"fmt"
	"io"
	"os"

	"github.com/hectorcorrea/marcli/pkg/marc"
)

// Mnemonic MARC, a human readable version
// See: https://librarycarpentry.org/lc-marcedit/03-working-with-MARC-files.html
func toMrk(params ProcessFileParams) error {
	if count == 0 {
		return nil
	}

	file, err := os.Open(params.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var i, out, recordCount int
	marc := marc.NewMarcFile(file)
	for marc.Scan() {

		r, err := marc.Record()
		if err == io.EOF {
			break
		}

		if err != nil {
			str := "== RECORD WITH ERROR STARTS HERE" + params.NewLine()
			str += "ERROR:" + params.NewLine() + err.Error() + params.NewLine()
			str += r.DebugString() + params.NewLine()
			str += "== RECORD WITH ERROR ENDS HERE" + params.NewLine() + params.NewLine()
			fmt.Print(str)
			if params.debug {
				continue
			}
			return err
		}

		if i++; i < start {
			continue
		}

		if r.Contains(params.searchValue, params.searchRegEx, params.searchFields) && r.HasFields(params.hasFields) {
			recordCount += 1
			str := ""
			if params.filters.IncludeLeader() {
				str += fmt.Sprintf("%s%s", r.Leader, params.NewLine())
			}
			for _, field := range r.Filter(params.filters, params.exclude) {
				str += fmt.Sprintf("%s%s", field, params.NewLine())
			}
			if str != "" {
				// Print the details of the record
				if params.format == "mrk" {
					fmt.Printf("%s%s", str, params.NewLine())
				}
				if out++; out == count {
					break
				}
			}
		}
	}

	// Print the count of records only
	if params.format == "count-only" {
		fmt.Printf("%d%s", recordCount, params.NewLine())
	}
	return marc.Err()
}
