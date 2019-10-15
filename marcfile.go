package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	rt = 0x1d // End of record
	st = 0x1f // End of subfield
)

type MarcFile struct {
	scanner *bufio.Scanner
}

func NewMarcFile(f *os.File) MarcFile {
	scanner := bufio.NewScanner(f)
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

func (file *MarcFile) Next() bool {
	x := file.scanner.Scan()
	if x == false {
		fmt.Printf("scanner returned false")
	}
	return x
}

// Value returns the current Record on the MarcIterator.
func (file *MarcFile) Value() (Record, error) {
	return file.scanIntoRecord(file.scanner.Bytes())
}

func (file *MarcFile) scanIntoRecord(bytes []byte) (Record, error) {
	rec := Record{}
	rec.Data = append([]byte(nil), bytes...)
	rec.Leader = Leader{
		Status:        bytes[5],
		Type:          bytes[6],
		BibLevel:      bytes[7],
		Control:       bytes[8],
		EncodingLevel: bytes[17],
		Form:          bytes[18],
		Multipart:     bytes[19],
	}

	start, err := strconv.Atoi(string(bytes[12:17]))
	if err != nil {
		return rec, errors.New("Could not determine record start")
	}
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
