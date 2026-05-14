package config

import (
	"fmt"
	"mals/pkg/config"
	"mals/pkg/core"
)

type Log struct {
	Name   string     `mapstructure:"name"`
	Level  LogLevel   `mapstructure:"level"`
	Output *LogOutput `mapstructure:"output"`
}

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type LogOutput struct {
	Kind LogOutputKind `mapstructure:"kind"`
	File *string       `mapstructure:"file"`
}

type LogOutputKind string

const (
	LogOutputKindFile LogOutputKind = "file"
)

func (o *Log) Wire(c *config.Log) error {
	o.Name = c.Name

	switch c.Level {
	case core.LogLevelError:
		o.Level = LogLevelError
	case core.LogLevelWarn:
		o.Level = LogLevelWarn
	case core.LogLevelInfo:
		o.Level = LogLevelInfo
	case core.LogLevelDebug:
		o.Level = LogLevelDebug
	default:
		return fmt.Errorf("unknown log level")
	}

	switch k := c.Output.(type) {
	case *config.LogOutputFile:
		o.Output = &LogOutput{
			Kind: LogOutputKindFile,
			File: &k.File,
		}
	default:
		return fmt.Errorf("unknown log kind")
	}

	return nil
}

func (o *Log) Unwire() (*config.Log, error) {
	c := &config.Log{
		Name: o.Name,
	}

	switch o.Level {
	case LogLevelError:
		c.Level = core.LogLevelError
	case LogLevelWarn:
		c.Level = core.LogLevelWarn
	case LogLevelInfo:
		c.Level = core.LogLevelInfo
	case LogLevelDebug:
		c.Level = core.LogLevelDebug
	default:
		return nil, fmt.Errorf("unknown log level")
	}

	switch o.Output.Kind {
	case LogOutputKindFile:
		output := &config.LogOutputFile{}
		if o.Output.File != nil {
			output.File = *o.Output.File
		}
		c.Output = output
	default:
		return nil, fmt.Errorf("unknown log kind: %v", o.Output.Kind)
	}

	return c, nil
}
