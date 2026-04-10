package main

import "strings"

func getLine(content string, lineNum int) string {
	lines := strings.Split(content, "\n")
	if lineNum < len(lines) {
		return strings.TrimRight(lines[lineNum], "\r")
	}
	return ""
}

func lastWord(s string) string {
	s = strings.TrimRight(s, " \t")
	for i := len(s) - 1; i >= 0; i-- {
		if !isIdentChar(s[i]) {
			return s[i+1:]
		}
	}
	return s
}

func dotReceiver(s string) string {
	s = strings.TrimRight(s, ". \t")
	return lastWord(s)
}

func wordAt(line string, col int) string {
	if col > len(line) {
		col = len(line)
	}
	start, end := col, col
	for start > 0 && isIdentChar(line[start-1]) {
		start--
	}
	for end < len(line) && isIdentChar(line[end]) {
		end++
	}
	return line[start:end]
}

func isIdentChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') || c == '_'
}
