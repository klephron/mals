package lsp

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
