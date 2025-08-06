package main

import (
	"fmt"
	"io"
	"os"

	"github.com/hectorcorrea/marcli/pkg/marc"
)

// Produces output that looks like the one produced by that yaz-marcdump utility
// See: https://software.indexdata.com/yaz/doc/yaz-marcdump.html
func toYaz(params ProcessFileParams) error {
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

		if r.Contains(params.searchValue, params.searchRegEx, params.searchFields) && r.HasFields(params.hasFields) {
			recordCount += 1
			str := ""
			if params.filters.IncludeLeader() {
				str += fmt.Sprintf("%s%s", r.Leader.Raw(), params.NewLine())
			}
			for _, field := range r.Filter(params.filters, params.exclude) {
				if field.IsControlField() {
					str += fmt.Sprintf("%s %s%s", field.Tag, field.Value, params.NewLine())
				} else {
					str += fmt.Sprintf("%s %s%s ", field.Tag, field.Indicator1, field.Indicator2)
					for _, sub := range field.SubFields {
						str += fmt.Sprintf("$%s %s ", sub.Code, sub.Value)
					}
					str += params.NewLine()
				}
			}
			if str != "" {
				fmt.Printf("%s", str)
				if out++; out == count {
					break
				}
			}
		}
	}
	return marc.Err()
}
