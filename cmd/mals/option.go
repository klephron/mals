package main

import (
	"encoding/json"
	"mals/pkg/core"
	"os"
	"path/filepath"
)

type options struct {
	ConfigPath string
}

func newOptions(args args) *options {
	options := options{
		ConfigPath: args.Config,
	}

	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".config")
	}
	if options.ConfigPath == "" {
		options.ConfigPath = filepath.Join(configDir, core.AppName, "config.toml")
	}

	return &options
}

func (s *options) String() string {
	byte, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(byte)
}
