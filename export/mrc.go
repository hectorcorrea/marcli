package export

import (
	"fmt"
	"io"
	"marcli/marc"
	"os"
)

func ToMrc(filename string, searchValue string, filters marc.FieldFilters) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	marc := marc.NewMarcFile(file)
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
