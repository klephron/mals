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
