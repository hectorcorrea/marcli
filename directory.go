package main

import (
	"fmt"
	"strconv"
)

type DirectoryEntry struct {
	Tag string
	Length int
	StartsAt int
}

func NewDirectoryEntryFromString(entry string) DirectoryEntry {
  if len(entry) != 12 {
    return DirectoryEntry{}
  }

	l, _ := strconv.Atoi(entry[3:7])
	s, _ := strconv.Atoi(entry[7:])
	dir := DirectoryEntry{
		Tag: entry[0:3],
		Length: l,
		StartsAt: s,
	}
	return dir
}

func (d DirectoryEntry) String() string {
	return fmt.Sprintf("tag: %s len: %d starts at: %d", d.Tag, d.Length, d.StartsAt)
}
