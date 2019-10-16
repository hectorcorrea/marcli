package export

import (
	"fmt"
	"io"
	"marcli/marc"
	"os"
)

func ToMrk(filename string, searchValue string, filters marc.FieldFilters) error {
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
			str := ""
			if filters.IncludeLeader() {
				str += fmt.Sprintf("%s\r\n", r.Leader)
			}
			// if filters.IncludeRecordInfo() {
			// 	str += fmt.Sprintf("=RIN  pos=%d, length=%d, data offset=%d\r\n", r.Pos, r.Leader.Length, r.Leader.DataOffset)
			// }
			// if filters.IncludeFileInfo() {
			// 	str += fmt.Sprintf("=FIN  %s\r\n", f.Name)
			// }
			for _, field := range r.Filter(filters) {
				str += fmt.Sprintf("%s\r\n", field)
			}
			if str != "" {
				fmt.Printf("%s\r\n", str)
			}
		}
	}

	return marc.Err()
}
