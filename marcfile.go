package main

import (
	"bufio"
	"io"
	"os"
)

const (
	RecordSeparator = 0x1e
)

type RecordProcessor interface {
	Process(Record)
}

type MarcFile struct {
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

func NewMarcFile(filename string) (MarcFile, error) {
	f, err := os.Open(filename)
	if err != nil {
		return MarcFile{}, err
	}
	return MarcFile{Name: filename, f: f, records: 0}, nil
}

func (file *MarcFile) ReadAll(processor RecordProcessor) error {
	for {
		_, err := file.readRecord(processor)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	file.f.Close()
	return nil
}

func (file *MarcFile) readRecord(processor RecordProcessor) (Record, error) {
	leader, err := file.readLeader()
	if err != nil {
		return Record{}, err
	}

	file.records += 1

	directory, err := file.readDirectory()
	if err != nil {
		panic(err)
	}
	values := file.readValues(directory)
	record := Record{
		Leader: leader,
		Fields: directory,
		Values: values,
		Pos:    file.records,
	}
	processor.Process(record)
	return Record{Leader: leader, Fields: directory}, nil
}

func (file *MarcFile) Close() {
	file.f.Close()
}

func (file *MarcFile) readLeader() (Leader, error) {
	bytes := make([]byte, 24)
	_, err := file.f.Read(bytes)
	if err != nil {
		return Leader{}, err
	}
	return NewLeader(string(bytes))
}

func (file *MarcFile) readDirectory() ([]Field, error) {
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

func (file *MarcFile) currentOffset() int64 {
	offset, _ := file.f.Seek(0, 1)
	return offset
}

func (file *MarcFile) readValues(entries []Field) []Value {
	values := make([]Value, len(entries))
	for i, entry := range entries {
		buffer := make([]byte, entry.Length)
		n, err := file.f.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}
		value := string(buffer[:n])
		values[i] = NewValue(entry.Tag, value)
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
