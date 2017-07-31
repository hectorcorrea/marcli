package main

import (
	"fmt"
	"strings"
)

type ConsoleProcessor struct {
	Fields []string
}

func (p ConsoleProcessor) Process(r Record) {
	if outputLeader(p.Fields) {
		fmt.Printf("=LDR  %s (%d, %d, %d)\n", r.Leader, r.Pos, r.Leader.Length, r.Leader.DataOffset)
	}
	for _, v := range r.Values {
		if outputField(p.Fields, v.Tag) {
			fmt.Printf("=%s  %s\r\n", v.Tag, v.Value)
		}
	}
	fmt.Printf("\r\n\r\n")
}

type ExtractProcessor struct {
	Fields []string
	Value  string // value to search
}

func (p ExtractProcessor) Process(r Record) {
	match := false
	for _, v := range r.Values {
		if strings.Contains(strings.ToLower(v.Value), p.Value) {
			match = true
			break
		}
	}

	if match {
		if outputLeader(p.Fields) {
			fmt.Printf("=LDR  %s (%d, %d, %d)\n", r.Leader, r.Pos, r.Leader.Length, r.Leader.DataOffset)
		}
		for _, v := range r.Values {
			if outputField(p.Fields, v.Tag) {
				fmt.Printf("=%s  %s\r\n", v.Tag, v.Value)
			}
		}
		fmt.Printf("\r\n\r\n")
	}
}

func outputLeader(arr []string) bool {
	return outputField(arr, "LDR")
}

func outputField(arr []string, value string) bool {
	if len(arr) == 0 {
		return true
	}
	for _, el := range arr {
		if value == el {
			return true
		}
	}
	return false
}
