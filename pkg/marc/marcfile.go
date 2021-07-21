package marc

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// See https://www.loc.gov/marc/specifications/specrecstruc.html
const (
	rt = 0x1d // End of record (MARC binary)
	st = 0x1f // End of subfield (MARC binary)
	ft = 0x1e // Field terminator (MARC binary)
)

// MarcFile represents a MARC file.
// The public interface more or less mimic Go's native Scanner (Scan, Err)
// but uses Record (instead of Text) to represent each MARC record.
type MarcFile struct {
	scanner *bufio.Scanner
	decoder *xml.Decoder
	isXML   bool
	element xml.StartElement
}

// NewMarcFile creates a struct to handle reading the MARC file.
func NewMarcFile(file *os.File) MarcFile {

	isXML := strings.HasSuffix(strings.ToLower(file.Name()), ".xml")
	if isXML {
		// For MARC XML files it uses a Decoder() to read one
		// MARC record at a time.
		decoder := xml.NewDecoder(file)
		return MarcFile{decoder: decoder, isXML: true}
	}

	// Assume MARC binary
	//
	// For MARC binary files uses a Scanner() to read the
	// contents of the file (stolen from https://github.com/MITLibraries/fml)
	scanner := bufio.NewScanner(file)

	// By default Scanner.Scan() returns "bufio.Scanner: token too long" if
	// the block to read is longer than 64K. Since MARC records can be up to
	// 100K we use a custom value. See https://stackoverflow.com/a/37455465/446681
	initialBuffer := make([]byte, 0, 64*1024)
	customMaxSize := 105 * 1024
	scanner.Buffer(initialBuffer, customMaxSize)

	scanner.Split(splitFunc)
	return MarcFile{scanner: scanner}
}

func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if atEOF {
		return len(data), data, nil
	}

	if i := bytes.IndexByte(data, rt); i >= 0 {
		return i + 1, data[0:i], nil
	}

	return 0, nil, nil
}

// Err returns the error in the scanner (if any)
func (file *MarcFile) Err() error {
	if file.isXML {
		return nil
	}
	return file.scanner.Err()
}

// Scan moves the scanner to the next record.
// Returns false when no more records can be read.
func (file *MarcFile) Scan() bool {

	if file.isXML {
		for {
			token, _ := file.decoder.Token()
			if token == nil {
				return false
			}
			// Find the next "<record>" element in the XML
			// and store it.
			element, ok := token.(xml.StartElement)
			if ok && element.Name.Local == "record" {
				file.element = element
				return true
			}
		}
	}

	return file.scanner.Scan()
}

// Record returns the current Record in the MarcFile.
func (file *MarcFile) Record() (Record, error) {
	rec := Record{}

	if file.isXML {
		// Decode the last element found in Scan() into an XML Record...
		var xmlRec XmlRecord
		file.decoder.DecodeElement(&xmlRec, &file.element)

		// Ignore error because a bad data offset is not a problem
		// in XML records.
		leader, _ := NewLeader([]byte(xmlRec.Leader))
		rec.Leader = leader
		rec.Data = []byte("Raw data not supported in XML format\n")

		// ...and then into a MARC Record.
		for _, control := range xmlRec.ControlFields {
			field := Field{Tag: control.Tag, Value: control.Value}
			rec.Fields = append(rec.Fields, field)
		}
		for _, data := range xmlRec.DataFields {
			field := Field{Tag: data.Tag, Indicator1: data.Ind1, Indicator2: data.Ind2}
			for _, sub := range data.SubFields {
				subfield := SubField{Code: sub.Code, Value: sub.Value}
				field.SubFields = append(field.SubFields, subfield)
			}
			rec.Fields = append(rec.Fields, field)
		}
		return rec, nil
	}

	// Parse the bytes from the scanner to create the MARC Record.
	recBytes := file.scanner.Bytes()
	rec.Data = append([]byte(nil), recBytes...)
	leader, err := NewLeader(recBytes[0:24])
	if err != nil {
		return rec, err
	}
	rec.Leader = leader

	start := leader.dataOffset
	if start <= 25 {
		return rec, errors.New("Bad data offset")
	} else if start > len(recBytes) {
		return rec, errors.New("Bad record length")
	}
	data := recBytes[start:]
	dirs := recBytes[24 : start-1]

	for len(dirs) >= 12 {
		tag := string(dirs[:3])
		length, err := strconv.Atoi(string(dirs[3:7]))
		if err != nil {
			return rec, errors.New("Could not determine length of field")
		}
		begin, err := strconv.Atoi(string(dirs[7:12]))
		if err != nil {
			return rec, errors.New("Could not determine field start")
		}
		if len(data) <= begin+length-1 {
			details := fmt.Sprintf("Tag: %s, len(data): %d, begin: %d, field length: %d",
				tag, len(data), begin, length)
			return rec, errors.New("Reported field length incorrect.\n" + details + "\n")
		}
		fdata := data[begin : begin+length-1] // length includes field terminator
		if len(fdata) > 4 {                   // ignore illegal data
			df, err := MakeField(tag, fdata)
			if err != nil {
				return rec, err
			}
			rec.Fields = append(rec.Fields, df)
		}
		dirs = dirs[12:]
	}
	return rec, nil
}
