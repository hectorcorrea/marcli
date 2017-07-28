package main

import (
	"errors"
	"fmt"
	"strconv"
)

type Record struct {
	Leader Leader
	Fields []Field
}

type Leader struct {
	raw        string
	Length     int
	DataOffset int
}

type Field struct {
	Tag      string
	Length   int
	StartsAt int
}

func NewLeader(value string) (Leader, error) {
	if len(value) != 24 {
		return Leader{}, errors.New("Incomplete leader")
	}
	l, _ := strconv.Atoi(value[0:5])
	o, _ := strconv.Atoi(value[12:17])
	return Leader{raw: value, Length: l, DataOffset: o}, nil
}

func (l Leader) String() string {
	return l.raw
}

func NewField(entry string) Field {
	if len(entry) != 12 {
		return Field{}
	}

	l, _ := strconv.Atoi(entry[3:7])
	s, _ := strconv.Atoi(entry[7:])
	dir := Field{
		Tag:      entry[0:3],
		Length:   l,
		StartsAt: s,
	}
	return dir
}

func (d Field) String() string {
	return fmt.Sprintf("tag: %s len: %d starts at: %d", d.Tag, d.Length, d.StartsAt)
}
