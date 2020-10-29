package main

import "strings"

func trimString(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.TrimRight(s, "\x00")
	return s
}
