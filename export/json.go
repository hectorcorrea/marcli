package export

import (
	"encoding/json"
	"fmt"
	"io"
	"marcli/marc"
	"os"
)

// TODO: Add support for JSONL (JSON line delimited) format that makes JSON
// easier to parse with Unix tools like grep, tail, and so on.
func ToJson(filename string, searchValue string, filters marc.FieldFilters) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	count := 0
	marc := marc.NewMarcFile(file)

	fmt.Printf("[\r\n")
	for marc.Scan() {
		r, err := marc.Record()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if r.Contains(searchValue) {
			if count > 0 {
				fmt.Printf(",\r\n")
			} else {
				fmt.Printf("\r\n")
			}
			b, err := json.Marshal(r.Filter(filters))
			if err != nil {
				fmt.Printf("%s\r\n", err)
			}
			// fmt.Printf("{ \"record\": %s}\r\n", b)
			fmt.Printf("%s", b)
			count++
		}
	}
	fmt.Printf("]\r\n")

	return marc.Err()
}
