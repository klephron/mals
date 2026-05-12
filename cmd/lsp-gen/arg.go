package main

import (
	"github.com/alexflint/go-arg"
)

type args struct {
	Target  string `arg:"-o" default:"third_party/lsp"`
	Package string `arg:"-p" default:"lsp"`
	Suffix  string `arg:"-s" default:"gen"`
}

func argParse() args {
	var args args
	arg.MustParse(&args)
	return args
}
