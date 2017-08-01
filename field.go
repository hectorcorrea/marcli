package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

func NewFieldsFromString(valueStr string) []SubFieldValue {
	var values []SubFieldValue
	// valueStr comes with the indicators, we skip them:
	//   value[0] indicator 1
	// 	 value[0] indicator 2
	// 	 value[0] separator (ascii 31)
	tokens := strings.Split(valueStr[3:], string(UnitSeparator))
	for _, token := range tokens {
		value := SubFieldValue{
			SubField: string(token[0]),
			Value:    token[1:],
		}
		values = append(values, value)
	}
	return values
}
