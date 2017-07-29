package main

import (
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
		_, err := f.ReadNext()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}
	f.Close()
}
