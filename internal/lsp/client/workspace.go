package client

import (
	"github.com/puzpuzpuz/xsync/v4"
)

type Workspace struct {
	uri       string
	name      string
	documents *xsync.Map[string, *Document]
}

func newWorkspace(uri string, name string) *Workspace {
	return &Workspace{
		uri:       uri,
		name:      name,
		documents: xsync.NewMap[string, *Document](),
	}
}

func (s *ClientLsp) documentGet(workspace *Workspace, uri string) *Document {
	document, ok := workspace.documents.Load(uri)
	if !ok {
		s.plane.Warnf("%s: workspace %s document %s is not present", s.Name(), workspace.name, document.uri)
		return nil
	}
	return document
}

func (s *ClientLsp) documentAdd(workspace *Workspace, uri string, text *string) {
	document := newDocument(uri, text)

	workspace.documents.Store(uri, document)

	s.plane.Infof("%s: workspace %s document %s added", s.Name(), workspace.name, document.uri)
}

func (s *ClientLsp) documentDelete(workspace *Workspace, uri string) {
	document, ok := workspace.documents.LoadAndDelete(uri)

	if !ok {
		s.plane.Warnf("%s: workspace %s document %s is not present", s.Name(), workspace.name, uri)
	}
	s.plane.Infof("%s: workspace %s document %s deleted", s.Name(), workspace.name, document.uri)
}
