package document

import (
	"mals/internal/lsp/protocol"
	"strings"
)

type Document struct {
	uri     string
	version int32
	lines   []string
}

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

func textToLines(text *string) []string {
	if text != nil {
		return strings.Split(*text, "\n")
	}
	return []string{}
}

func New(uri string, text *string, version int32) *Document {
	return &Document{
		uri:     uri,
		version: version,
		lines:   textToLines(text),
	}
}

func (d *Document) Uri() string {
	return d.uri
}

func (d *Document) Version() int32 {
	return d.version
}

func (d *Document) LastNChars(line uint32, char uint32, n uint64) string {
	position := d.position(line, char)
	return d.positionLastNChars(position, n)
}

func (d *Document) Text() string {
	return strings.Join(d.lines, "\n")
}

func (d *Document) TextUpdate(version int32, changes []protocol.TextDocumentContentChangeEvent) bool {
	if d.version >= version {
		return false
	}
	d.version = version

	for _, change := range changes {
		if change.Range == nil {
			d.lines = textToLines(&change.Text)
			continue
		}

		startLine := int(change.Range.Start.Line)
		startChar := int(change.Range.Start.Character)
		endLine := int(change.Range.End.Line)
		endChar := int(change.Range.End.Character)

		// Bounds validation
		if startLine >= len(d.lines) || endLine >= len(d.lines) {
			return false
		}
		if startChar > len(d.lines[startLine]) || endChar > len(d.lines[endLine]) {
			return false
		}

		prefix := d.lines[startLine][:startChar]
		suffix := d.lines[endLine][endChar:]

		newLines := textToLines(&change.Text)
		if len(newLines) == 0 {
			newLines = []string{""}
		}
		newLines[0] = prefix + newLines[0]
		newLines[len(newLines)-1] = newLines[len(newLines)-1] + suffix

		result := make([]string, 0, len(d.lines)-(endLine-startLine+1)+len(newLines))
		result = append(result, d.lines[:startLine]...)
		result = append(result, newLines...)
		result = append(result, d.lines[endLine+1:]...)

		d.lines = result
	}

	return true
}

func (d *Document) TextBefore(line, char uint32) string {
	pos := d.position(line, char)
	content := d.Text()
	if pos > uint64(len(content)) {
		pos = uint64(len(content))
	}
	return content[:pos]
}

func (d *Document) TextAfter(line, char uint32) string {
	pos := d.position(line, char)
	content := d.Text()
	if pos >= uint64(len(content)) {
		return ""
	}
	return content[pos:]
}
