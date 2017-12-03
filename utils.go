package main

import (
	"fmt"
	"strings"
)

func pad(str string) string {
	if len(str) > 40 {
		return str[0:40]
	}
	return fmt.Sprintf("%-40s", str)
}

func concat(a, b string) string {
	return _concat(a, b, " ")
}

func concatTab(a, b string) string {
	return _concat(a, b, "\t")
}

func _concat(a, b, sep string) string {
	if a == "" && b == "" {
		return ""
	} else if a == "" && b != "" {
		return strings.TrimSpace(b)
	} else if a != "" && b == "" {
		return strings.TrimSpace(a)
	}
	return strings.TrimSpace(a) + sep + strings.TrimSpace(b)
}

func concat3(a, b, c string) string {
	return concat(concat(a, b), c)
}

func removeSpaces(s string) string {
	return strings.Replace(s, " ", "", -1)
}
