package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("./one.mrc")
	if err != nil {
		panic(err)
	}

	leader, err := readLeader(f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("leader: %s\n", leader)

	directory, err := readDirectory(f)
	for i, entry := range directory {
		fmt.Printf("(%d) %s\r\n", i, entry)
	}

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

func readDirectory(f *os.File) ([]string, error) {
	// Source: https://www.socketloop.com/references/golang-bufio-scanrunes-function-example
	reader := bufio.NewReader(f)
	bb, err := reader.ReadBytes('^')
	if err != nil {
		return nil, err
	}
	count := (len(bb) - 1) / 12
	entries := make([]string, count)
	for i := 0; i < count; i++ {
		start := i * 12
		entry := string(bb[start : start+12])
		entries[i] = entry
	}
	return entries, nil
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
