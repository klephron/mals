package middleware

import "mals/internal/middleware/document"

func (s *Middleware) documentGet(workspace *Workspace, uri string) *document.Document {
	document, ok := workspace.documents.Load(uri)
	if !ok {
		s.plane.Warnf("%s: workspace %s document %s is not present", s.Name(), workspace.name, document.Uri())
		return nil
	}
	return document
}

func (s *Middleware) documentAdd(workspace *Workspace, uri string, text *string) {
	document := document.New(uri, text)

	workspace.documents.Store(uri, document)

	s.plane.Infof("%s: workspace %s document %s added", s.Name(), workspace.name, document.Uri())
}

func (s *Middleware) documentDelete(workspace *Workspace, uri string) {
	document, ok := workspace.documents.LoadAndDelete(uri)

	if !ok {
		s.plane.Warnf("%s: workspace %s document %s is not present", s.Name(), workspace.name, uri)
	}
	s.plane.Infof("%s: workspace %s document %s deleted", s.Name(), workspace.name, document.Uri())
}

func (s *Middleware) documentUpdateFull(workspace *Workspace, document *document.Document, text *string, version int32) {

	if document.SetVersion(version) {
		s.plane.Warnf("%s: workspace %s document %s version %d >= %d",
			s.Name(), workspace.name, document.Uri(), document.Version(), version)
		return
	}

	document.SetText(text)

	s.plane.Warnf("%s: workspace %s document %s updated version %d",
		s.Name(), workspace.name, document.Uri(), document.Version())
}
