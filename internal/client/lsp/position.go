package lsp

import "strings"

func positionGet(content string, line, char uint32) uint64 {
	lines := strings.Split(content, "\n")

	if line >= uint32(len(lines)) {
		return uint64(len(content))
	}

	var position uint64
	// Add lengths of all previous lines (including newlines)
	for i := range line {
		position += uint64(len(lines[i])) + 1 // +1 for newline
	}

	// Add column position (0-based)
	if char <= uint32(len(lines[line])) {
		position += uint64(char)
	} else {
		position += uint64(len(lines[line]))
	}

	return position
}

func positionLastNCharsByPos(content string, position uint64, n uint64) string {
	contentLen := uint64(len(content))
	if position > contentLen {
		position = contentLen
	}

	var start uint64
	if position > n {
		start = position - n
	} else {
		start = 0
	}

	return content[start:position]
}

func positionLastNChars(content string, line, char uint32, n uint64) string {
	position := positionGet(content, line, char)
	return positionLastNCharsByPos(content, position, n)
}
