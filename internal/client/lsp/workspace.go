package lsp

import "github.com/puzpuzpuz/xsync/v4"

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
