package lsp

import "github.com/puzpuzpuz/xsync/v4"

type State struct {
	initialized bool
	workspaces  *xsync.Map[string, *Workspace]
}

type Workspace struct {
}
