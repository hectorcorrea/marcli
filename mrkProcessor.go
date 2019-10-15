package main

import (
	"fmt"
	"io"
	"os"
)

func mrkProcessor(filename string, searchValue string, filters FieldFilters) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	marc := NewMarcFile(file)
	for marc.Next() {

		r, err := marc.Value()
		if err == io.EOF {
			fmt.Printf("EOF\r\n")
			break
		}
		if err != nil {
			fmt.Printf("Error %s\r\n", err)
			return err
		}

		if r.ControlNum() == "ocm02647229" {
			fmt.Printf("read ocm02647229\r\n")
		}

		if r.IsMatch(searchValue) {
			str := ""
			// if filters.IncludeLeader() {
			// 	str += fmt.Sprintf("%s\r\n", r.Leader)
			// }
			// if filters.IncludeRecordInfo() {
			// 	str += fmt.Sprintf("=RIN  pos=%d, length=%d, data offset=%d\r\n", r.Pos, r.Leader.Length, r.Leader.DataOffset)
			// }
			// if filters.IncludeFileInfo() {
			// 	str += fmt.Sprintf("=FIN  %s\r\n", f.Name)
			// }
			// filteredFields := filters.Apply(r.Fields)
			// for _, field := range filteredFields.All() {
			// 	str += fmt.Sprintf("%s\r\n", field)
			// }
			for _, f := range r.Fields {
				str += fmt.Sprintf("%s\r\n", f)
			}
			if str != "" {
				fmt.Printf("%s\r\n", str)
			}
		}
	}

	return marc.Err()
}
