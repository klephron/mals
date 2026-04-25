package workspace

import (
	"fmt"
	"mals/internal/lsp/protocol"
	"mals/internal/middleware/document"

	"github.com/puzpuzpuz/xsync/v4"
)

type Workspace struct {
	uri       string
	name      string
	documents *xsync.Map[string, *document.Document]
}

func New(uri string, name string) *Workspace {
	return &Workspace{
		uri:       uri,
		name:      name,
		documents: xsync.NewMap[string, *document.Document](),
	}
}

func (s *Workspace) Uri() string {
	return s.uri
}

func (s *Workspace) Name() string {
	return s.name
}

func (s *Workspace) DocumentGet(uri string) (*document.Document, error) {
	document, ok := s.documents.Load(uri)
	if !ok {
		return nil, fmt.Errorf("workspace %s document %s is not present", s.name, document.Uri())
	}
	return document, nil
}

func (s *Workspace) DocumentAdd(uri string, text *string, version int32) {
	document := document.New(uri, text, version)

	s.documents.Store(uri, document)
}

func (s *Workspace) DocumentDelete(uri string) error {
	_, ok := s.documents.LoadAndDelete(uri)

	if !ok {
		return fmt.Errorf("workspace %s document %s is not present", s.name, uri)
	}
	return nil
}

func (s *Workspace) DocumentUpdate(document *document.Document, version int32, changes []protocol.TextDocumentContentChangeEvent) error {
	ok := document.TextUpdate(version, changes)

	if !ok {
		return fmt.Errorf("workspace %s document %s version %d >= %d",
			s.name, document.Uri(), document.Version(), version)
	}

	return nil
}
