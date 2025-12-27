package lsp

type Document struct {
	lines []string
}

func newDocument(lines []string) *Document {
	return &Document{
		lines: lines,
	}
}
