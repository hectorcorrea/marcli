package main

import (
	"errors"
	"strconv"
)

type Leader struct {
	raw        string
	Length     int
	DataOffset int
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
