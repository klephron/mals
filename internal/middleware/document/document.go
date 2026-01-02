package document

import "strings"

type Document struct {
	uri     string
	version int32
	lines   []string
}

func New(uri string, text *string) *Document {
	return &Document{
		uri:     uri,
		version: 0,
		lines:   textToLines(text),
	}
}

func (d *Document) Uri() string {
	return d.uri
}

func (d *Document) Version() int32 {
	return d.version
}

func (d *Document) SetVersion(version int32) bool {
	if d.version >= version {
		return false
	}
	d.version = version
	return true
}

func (d *Document) LastNChars(line uint32, char uint32, n uint64) string {
	position := d.position(line, char)
	return d.positionLastNChars(position, n)
}

func (d *Document) Text() string {
	return strings.Join(d.lines, "\n")
}

func (d *Document) SetText(text *string) {
	d.lines = textToLines(text)
}
