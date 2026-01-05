package main

import (
	"github.com/alexflint/go-arg"
)

type Params struct {
	Target  string `arg:"-o" default:"internal/lsp/protocol"`
	Package string `arg:"-p" default:"protocol"`
	Suffix  string `arg:"-s" default:"gen"`
}

func argParse() Params {
	var params Params
	arg.MustParse(&params)
	return params
}
