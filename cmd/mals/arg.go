package main

import "github.com/alexflint/go-arg"

type args struct {
	Config string `arg:"-c" default:"" help:"config file path"`
}

func argParse() args {
	var args args
	arg.MustParse(&args)
	return args
}
