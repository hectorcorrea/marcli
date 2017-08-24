package main

import (
	"errors"
	"fmt"
	"strconv"
)

type DirEntry struct {
	Tag      string
	Length   int
	StartsAt int
}

func NewDirEntry(entry string) (DirEntry, error) {
	if len(entry) != 12 {
		return DirEntry{}, errors.New("Incomplete field definition")
	}

	l, _ := strconv.Atoi(entry[3:7])
	s, _ := strconv.Atoi(entry[7:])
	dir := DirEntry{
		Tag:      entry[0:3],
		Length:   l,
		StartsAt: s,
	}
	return dir, nil
}

func (d DirEntry) String() string {
	return fmt.Sprintf("tag: %s len: %d starts at: %d", d.Tag, d.Length, d.StartsAt)
}
