package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

const (
	RecordSeparator = 0x1e
	UnitSeparator   = 0x1f
)

type RecordProcessor interface {
	Process(Record, int)
}

type File struct {
	Name    string
	f       *os.File
	records int
}

type Record struct {
	Leader Leader
	Fields []Field
	Values []Value
	Pos    int
}

type Value struct {
	Tag   string
	Ind1  string
	Ind2  string
	Value string
}

func NewFile(filename string) (File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return File{}, err
	}
	return File{Name: filename, f: f, records: 0}, nil
}

func (file *File) ReadNext(processor func(Record)) (Record, error) {
	leader, err := file.readLeader()
	if err != nil {
		return Record{}, err
	}

	file.records += 1

	directory, err := file.readDirectory()
	if err != nil {
		panic(err)
	}
	// for i, entry := range directory {
	// 	fmt.Printf("(%d) %s\r\n", i, entry)
	// }
	values := file.readValues(directory)
	record := Record{
		Leader: leader,
		Fields: directory,
		Values: values,
		Pos:    file.records,
	}
	processor(record)
	return Record{Leader: leader, Fields: directory}, nil
}

func (file *File) Close() {
	file.f.Close()
}

func (file *File) readLeader() (Leader, error) {
	bytes := make([]byte, 24)
	_, err := file.f.Read(bytes)
	if err != nil {
		return Leader{}, err
	}
	return NewLeader(string(bytes))
}

func (file *File) readDirectory() ([]Field, error) {
	// Source: https://www.socketloop.com/references/golang-bufio-scanrunes-function-example
	offset := file.currentOffset()
	reader := bufio.NewReader(file.f)
	ss, err := reader.ReadString(RecordSeparator)
	if err != nil {
		return nil, err
	}
	count := (len(ss) - 1) / 12
	entries := make([]Field, count)
	for i := 0; i < count; i++ {
		start := i * 12
		entry := ss[start : start+12]
		field, err := NewField(entry)
		if err != nil {
			return nil, err
		}
		entries[i] = field
	}
	// ReadString leaves the file pointer a bit further than we want to.
	// Force it to be exactly at the end of the directory.
	file.f.Seek(offset+int64(len(ss)), 0)
	return entries, nil
}

func (file *File) currentOffset() int64 {
	offset, _ := file.f.Seek(0, 1)
	return offset
}

func (file *File) readValues(entries []Field) []Value {
	values := make([]Value, len(entries))
	for i, entry := range entries {
		buffer := make([]byte, entry.Length)
		n, err := file.f.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}
		value := string(buffer[:n])
		if entry.Tag > "009" {
			value = formatValue(value)
		}
		values[i].Tag = entry.Tag
		values[i].Value = value
		// fmt.Printf("=%s  %s\r\n", entry.Tag, value)
	}

	eor := make([]byte, 1)
	n, err := file.f.Read(eor)
	if n != 1 {
		panic("End of record byte not found")
	}

	if err != nil {
		panic(err)
	}
	return values
}

func formatValue(value string) string {
	formatted := ""
	formatted += formatIndicator(value[0])
	formatted += formatIndicator(value[1])
	formatted += string(value[2:])
	sep := string(byte(UnitSeparator))
	return strings.Replace(formatted, sep, "$", -1)
}

func formatIndicator(value byte) string {
	if value == ' ' {
		return "\\"
	}
	return string(value)
}
