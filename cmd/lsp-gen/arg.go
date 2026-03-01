package main

import (
	"github.com/alexflint/go-arg"
)

type args struct {
	Target  string `arg:"-o" default:"internal/lsp/protocol"`
	Package string `arg:"-p" default:"protocol"`
	Suffix  string `arg:"-s" default:"gen"`
}

func argParse() args {
	var args args
	arg.MustParse(&args)
	return args
}
