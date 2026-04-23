package handler

import (
	"mals/internal/lsp/protocol"
	"mals/internal/middleware/document"
)

func (s *Handler) documentGet(workspace *workspace, uri string) *document.Document {
	document, ok := workspace.documents.Load(uri)
	if !ok {
		s.plane.Warnf("%T %v: workspace %s document %s is not present",
			s, s.Name(), workspace.name, document.Uri())
		return nil
	}
	return document
}

func (s *Handler) documentAdd(workspace *workspace, uri string, text *string, version int32) {
	document := document.New(uri, text, version)

	workspace.documents.Store(uri, document)

	s.plane.Infof("%T %v: workspace %s document %s added",
		s, s.Name(), workspace.name, document.Uri())
}

func (s *Handler) documentDelete(workspace *workspace, uri string) {
	document, ok := workspace.documents.LoadAndDelete(uri)

	if ok {
		s.plane.Infof("%T %v: workspace %s document %s deleted",
			s, s.Name(), workspace.name, document.Uri())
	} else {
		s.plane.Warnf("%T %v: workspace %s document %s is not present",
			s, s.Name(), workspace.name, uri)
	}
}

func (s *Handler) documentUpdate(workspace *workspace, document *document.Document, version int32, changes []protocol.TextDocumentContentChangeEvent) {
	ok := document.TextUpdate(version, changes)

	if ok {
		s.plane.Infof("%T %v: workspace %s document %s updated version %d",
			s, s.Name(), workspace.name, document.Uri(), document.Version())
	} else {
		s.plane.Warnf("%T %v: workspace %s document %s version %d >= %d",
			s, s.Name(), workspace.name, document.Uri(), document.Version(), version)
	}
}
