package main

import (
	"encoding/json"
	"github.com/alexflint/go-arg"
)

type Params struct {
	Config string `arg:"-c" default:"" help:"config file path"`
}

func argParse() Params {
	var params Params
	arg.MustParse(&params)
	return params
}

func (p Params) String() string {
	byte, err := json.Marshal(p)
	if err != nil {
		return err.Error()
	}
	return string(byte)
}
