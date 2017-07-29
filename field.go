package main

import (
	"errors"
	"fmt"
	"strconv"
)

type Field struct {
	Tag      string
	Length   int
	StartsAt int
}

func NewField(entry string) (Field, error) {
	if len(entry) != 12 {
		return Field{}, errors.New("Incomplete field definition")
	}

	l, _ := strconv.Atoi(entry[3:7])
	s, _ := strconv.Atoi(entry[7:])
	dir := Field{
		Tag:      entry[0:3],
		Length:   l,
		StartsAt: s,
	}
	return dir, nil
}

func (d Field) String() string {
	return fmt.Sprintf("tag: %s len: %d starts at: %d", d.Tag, d.Length, d.StartsAt)
}
