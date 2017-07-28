package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	RecordSeparator = 0x1e
	UnitSeparator   = 0x1f
)

func main() {
	if len(os.Args) < 2 {
		panic("Must provide name of MARC file to process")
	}

	fileName := os.Args[1]
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	i := 0
	for {
		leader, err := readLeader(f)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		i += 1
		fmt.Printf("=LDR  %s (%d, %d)\n", leader, leader.Length, leader.DataOffset)

		directory, err := readDirectory(f)
		if err != nil {
			panic(err)
		}
		// for i, entry := range directory {
		// 	fmt.Printf("(%d) %s\r\n", i, entry)
		// }

		readRecord(f, directory)
		fmt.Printf("\r\n\r\n")
	}
	f.Close()
}

func readLeader(f *os.File) (Leader, error) {
	bytes := make([]byte, 24)
	_, err := f.Read(bytes)
	if err != nil {
		return Leader{}, err
	}
	return NewLeader(string(bytes))
}

func readDirectory(f *os.File) ([]Field, error) {
	// Source: https://www.socketloop.com/references/golang-bufio-scanrunes-function-example
	offset := currentOffset(f)
	reader := bufio.NewReader(f)
	ss, err := reader.ReadString(RecordSeparator)
	if err != nil {
		return nil, err
	}
	count := (len(ss) - 1) / 12
	entries := make([]Field, count)
	for i := 0; i < count; i++ {
		start := i * 12
		entry := ss[start : start+12]
		entries[i] = NewField(entry)
	}
	// ReadString leaves the file pointer a bit further than we want to.
	// Force it to be exactly at the end of the directory.
	f.Seek(offset+int64(len(ss)), 0)
	return entries, nil
}

func currentOffset(f *os.File) int64 {
	offset, _ := f.Seek(0, 1)
	return offset
}

func readRecord(f *os.File, entries []Field) {
	for _, entry := range entries {
		buffer := make([]byte, entry.Length)
		n, err := f.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}
		value := string(buffer[:n])
		if entry.Tag > "009" {
			value = formatValue(value)
		}
		fmt.Printf("=%s  %s\r\n", entry.Tag, value)
	}

	eor := make([]byte, 1)
	n, err := f.Read(eor)
	if n != 1 {
		panic("End of record byte not found")
	}

	if err != nil {
		panic(err)
	}
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
