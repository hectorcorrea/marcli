package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"io"
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

	leader, err := readLeader(f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("leader: %s\n", leader)

	directory, err := readDirectory(f)
	if err != nil {
		panic(err)
	}
	for i, entry := range directory {
		fmt.Printf("(%d) %s\r\n", i, entry)
	}

	readRecords(f, directory)
	f.Close()
}

func readLeader(f *os.File) (string, error) {
	bytes := make([]byte, 24)
	n, err := f.Read(bytes)
	if err != nil {
		return "", err
	}
	if n != 24 {
		return "", errors.New("Incomplete leader.")
	}
	return string(bytes), nil
}

// Using ReadString
func readDirectory(f *os.File) ([]DirectoryEntry, error) {
	// Source: https://www.socketloop.com/references/golang-bufio-scanrunes-function-example
	reader := bufio.NewReader(f)
	ss, err := reader.ReadString('^')
	if err != nil {
		return nil, err
	}
	fmt.Printf("len of directory: %d\r\n", len(ss))
	count := (len(ss) - 1) / 12
	entries := make([]DirectoryEntry, count)
	for i := 0; i < count; i++ {
		start := i * 12
		entry := ss[start : start+12]
		entries[i] = NewDirectoryEntryFromString(entry)
	}
	return entries, nil
}

// Using ReadBytes
func readDirectory00(f *os.File) ([]DirectoryEntry, error) {
	// Source: https://www.socketloop.com/references/golang-bufio-scanrunes-function-example
	reader := bufio.NewReader(f)
	bb, err := reader.ReadBytes('^')
	if err != nil {
		return nil, err
	}
	fmt.Printf("len of directory: %d", len(bb))
	return nil, errors.New("stop")
	count := (len(bb) - 1) / 12
	entries := make([]DirectoryEntry, count)
	for i := 0; i < count; i++ {
		start := i * 12
		entry := string(bb[start : start+12])
		entries[i] = NewDirectoryEntryFromString(entry)
	}
	return entries, nil
}

func readRecords(f *os.File, entries []DirectoryEntry) {
	offset := 24 + (len(entries) * 12)
	f.Seek(int64(offset), 0)
	for _, entry := range entries {
		buffer := make([]byte, entry.Length)
		n, err := f.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}
		fmt.Printf("%s=%s\r\n", entry.Tag, string(buffer[:n]))
	}
}


// func readDirectoryEntry(f *os.File) (string, error) {
// 	bytes := make([]byte, 12)
// 	n, err := f.Read(bytes)
// 	if err != nil {
// 		return "", err
// 	}
// 	if n != 12 {
// 		return "", errors.New("Incomplete directory entry.")
// 	}
// 	return string(bytes), nil
// }

// reader := bufio.NewReader(f)
// xx, err := reader.ReadBytes('^')
// fmt.Printf("%t", xx)
// fmt.Printf("%d", len(xx))

// One by one, how do we detect we are at the end?
// d1, err := readDirectoryEntry(f)
// fmt.Printf("dir: %s\n", d1)
//
// d2, err := readDirectoryEntry(f)
// fmt.Printf("dir: %s\n", d2)
