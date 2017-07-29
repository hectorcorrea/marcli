package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("Must provide name of MARC file to process")
	}

	fileName := os.Args[1]
	f, err := NewFile(fileName)
	if err != nil {
		panic(err)
	}

	for {
		_, err := f.ReadNext(PrintIt)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}
	f.Close()
}

func PrintIt(r Record) {
	fmt.Printf("=LDR  %s (%d, %d, %d)\n", r.Leader, r.Pos, r.Leader.Length, r.Leader.DataOffset)
	for _, v := range r.Values {
		fmt.Printf("=%s  %s\r\n", v.Tag, v.Value)
	}
	fmt.Printf("\r\n\r\n")
}
