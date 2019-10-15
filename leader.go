package main

import (
	"fmt"
)

// Leader contains a subset of the bytes in the record leader. Omitted are
// bytes specifying the length of parts of the record and bytes which do
// not vary from record to record.
type Leader struct {
	Status        byte // 05 byte position
	Type          byte // 06
	BibLevel      byte // 07
	Control       byte // 08
	EncodingLevel byte // 17
	Form          byte // 18
	Multipart     byte // 19
}

// type Leader struct {
// 	raw        string
// 	Length     int
// 	DataOffset int
// }

// func NewLeader(value string) (Leader, error) {
// 	if len(value) != 24 {
// 		return Leader{}, errors.New("Incomplete leader")
// 	}
// 	l, _ := strconv.Atoi(value[0:5])
// 	o, _ := strconv.Atoi(value[12:17])
// 	return Leader{raw: value, Length: l, DataOffset: o}, nil
// }

func (l Leader) String() string {
	return fmt.Sprintf("=LDR  %s", "TODO")
}
