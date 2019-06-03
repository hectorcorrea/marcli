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
	raw      string
}

func NewDirEntry(entry string) (DirEntry, error) {
	if len(entry) != 12 {
		return DirEntry{raw: entry}, errors.New("Incomplete field definition")
	}

	length, _ := strconv.Atoi(entry[3:7])
	if length == 0 {
		return DirEntry{raw: entry}, errors.New("Empty directory entry detected")
	}

	startsAt, _ := strconv.Atoi(entry[7:])
	dir := DirEntry{
		Tag:      entry[0:3],
		Length:   length,
		StartsAt: startsAt,
		raw:      entry,
	}
	return dir, nil
}

func (d DirEntry) String() string {
	if d.Tag == "" {
		return fmt.Sprintf("raw: %s", d.raw)
	}
	return fmt.Sprintf("tag: %s len: %d starts at: %d", d.Tag, d.Length, d.StartsAt)
}
