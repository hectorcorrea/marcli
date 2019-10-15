package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strconv"
)

const (
	rt = 0x1d // End of record
	st = 0x1f // End of subfield
)

// MarcFile represents a MARC file and mimic Go's native Scanner
// interface (Scan, Err, Text)
type MarcFile struct {
	scanner *bufio.Scanner
}

// NewMarcFile creates a scanner to manage reading the contents
// of the MARC file using Go's native Scanner interface.
// (stolen from https://github.com/MITLibraries/fml)
func NewMarcFile(file *os.File) MarcFile {
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

	return
}

// Err returns the error in the scanner (if any)
func (file *MarcFile) Err() error {
	return file.scanner.Err()
}

// Scan moves the scanner to the next record.
// Returns false when no more records can be read.
func (file *MarcFile) Scan() bool {
	return file.scanner.Scan()
}

// Record returns the current Record in the MarcFile.
func (file *MarcFile) Record() (Record, error) {
	bytes := file.scanner.Bytes()
	rec := Record{}
	rec.Data = append([]byte(nil), bytes...)

	leader, err := NewLeader(bytes[0:24])
	if err != nil {
		return rec, err
	}
	rec.Leader = leader

	start := leader.dataOffset
	data := bytes[start:]
	dirs := bytes[24 : start-1]

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
			return rec, errors.New("Reported field length incorrect")
		}
		fdata := data[begin : begin+length-1] // length includes field terminator
		df, err := MakeField(tag, fdata)
		if err != nil {
			return rec, err
		}
		rec.Fields = append(rec.Fields, df)
		dirs = dirs[12:]
	}
	return rec, nil
}
