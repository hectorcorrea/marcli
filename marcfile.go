package main

import (
	"bufio"
	"io"
	"os"
)

type Processor interface {
	ProcessRecord(*MarcFile, Record)
	Header()
	Footer()
	Separator()
}

type MarcFile struct {
	Name        string
	f           *os.File
	records     int
	outputCount int
}

func NewMarcFile(filename string) (MarcFile, error) {
	f, err := os.Open(filename)
	if err != nil {
		return MarcFile{}, err
	}
	return MarcFile{Name: filename, f: f, records: 0}, nil
}

func (file *MarcFile) ReadAll(processor Processor, searchValue string) error {
	processor.Header()
	for {
		record, err := file.readRecord(processor)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		file.records += 1

		if record.IsMatch(searchValue) {
			if file.outputCount > 0 {
				processor.Separator()
			}
			processor.ProcessRecord(file, record)
			file.outputCount += 1
		}
	}
	file.f.Close()
	processor.Footer()
	return nil
}

func (file *MarcFile) readRecord(processor Processor) (Record, error) {
	leader, err := file.readLeader()
	if err != nil {
		return Record{}, err
	}

	directory, err := file.readDirectory()
	if err != nil {
		panic(err)
	}
	fields := file.readValues(directory)
	record := Record{
		Leader:    leader,
		Directory: directory,
		Fields:    fields,
		Pos:       file.records,
	}
	return record, nil
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

func (file *MarcFile) readDirectory() ([]DirEntry, error) {
	const RecordSeparator = 0x1e

	// Source: https://www.socketloop.com/references/golang-bufio-scanrunes-function-example
	offset := file.currentOffset()
	reader := bufio.NewReader(file.f)
	ss, err := reader.ReadString(RecordSeparator)
	if err != nil {
		return nil, err
	}
	count := (len(ss) - 1) / 12
	directory := make([]DirEntry, count)
	for i := 0; i < count; i++ {
		start := i * 12
		entry := ss[start : start+12]
		field, err := NewDirEntry(entry)
		if err != nil {
			return nil, err
		}
		directory[i] = field
	}
	// ReadString leaves the file pointer a bit further than we want to.
	// Force it to be exactly at the end of the directory.
	file.f.Seek(offset+int64(len(ss)), 0)
	return directory, nil
}

func (file *MarcFile) currentOffset() int64 {
	offset, _ := file.f.Seek(0, 1)
	return offset
}

func (file *MarcFile) readValues(directory []DirEntry) Fields {
	var fields Fields
	for _, entry := range directory {
		buffer := make([]byte, entry.Length)
		n, err := file.f.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}
		value := string(buffer[:n-1]) // -1 to exclude the record separator character (0x1e)
		field := NewField(entry.Tag, value)
		fields.Add(field)
	}

	eor := make([]byte, 1)
	n, err := file.f.Read(eor)
	if n != 1 {
		panic("End of record byte not found")
	}

	if err != nil {
		panic(err)
	}
	return fields
}
