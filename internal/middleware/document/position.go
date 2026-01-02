package document

import "strings"

func (d *Document) position(line, char uint32) uint64 {
	if line >= uint32(len(d.lines)) {
		// Calculate total length including newlines
		var total uint64
		for i := range d.lines {
			total += uint64(len(d.lines[i]))
			if i < len(d.lines)-1 {
				total++ // newline
			}
		}
		return total
	}

	var position uint64
	// Add lengths of all previous lines (including newlines)
	for i := range line {
		position += uint64(len(d.lines[i])) + 1 // +1 for newline
	}

	// Add column position (0-based)
	if char <= uint32(len(d.lines[line])) {
		position += uint64(char)
	} else {
		position += uint64(len(d.lines[line]))
	}

	return position
}

func (d *Document) positionLastNChars(position uint64, n uint64) string {
	content := strings.Join(d.lines, "\n")
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
