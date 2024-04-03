package strutil

import (
	"bufio"
	"bytes"
	"strings"
)

// ToLines parses a multiline text of string into a slice of line item.
func ToLines(text string) []string {
	lines := []string{}
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

// ScanLines just like bufio.ScanLines, bufio.ScanLines treats `\r\n` or `\n` as line break,
// while ScanLines also treat `\r` as line break.
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		if data[i] == '\n' {
			// We have a line terminated by single newline.
			return i + 1, data[0:i], nil
		}
		advance = i + 1
		if len(data) > i+1 && data[i+1] == '\n' {
			advance += 1
		}
		return advance, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
