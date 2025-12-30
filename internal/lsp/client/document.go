package client

import "strings"

type Document struct {
	uri     string
	version int32
	lines   []string
}

func documentTextToLines(text *string) []string {
	return strings.Split(*text, "\n")
}

func newDocument(uri string, text *string) *Document {
	return &Document{
		uri:     uri,
		version: 0,
		lines:   documentTextToLines(text),
	}
}

func (s *ClientLsp) documentUpdateFull(workspace *Workspace, document *Document, text *string, version int32) {
	if document.version >= version {
		s.plane.Log().Warnf("%s: workspace %s document %s version %d >= %d",
			s.Name(), workspace.name, document.uri, document.version, version)
		return
	}

	document.version = version
	document.lines = documentTextToLines(text)

	s.plane.Log().Warnf("%s: workspace %s document %s updated version %d",
		s.Name(), workspace.name, document.uri, document.version)
}

func (d *Document) text() string {
	return strings.Join(d.lines, "\n")
}

func (d *Document) positionGet(line, char uint32) uint64 {
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

func (d *Document) lastNCharsByPosition(position uint64, n uint64) string {
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

func (d *Document) lastNChars(line uint32, char uint32, n uint64) string {
	position := d.positionGet(line, char)
	return d.lastNCharsByPosition(position, n)
}
