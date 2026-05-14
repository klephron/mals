package config

import (
	"fmt"
	"mals/pkg/config"
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
	case config.LogLevelError:
		o.Level = LogLevelError
	case config.LogLevelWarn:
		o.Level = LogLevelWarn
	case config.LogLevelInfo:
		o.Level = LogLevelInfo
	case config.LogLevelDebug:
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
		c.Level = config.LogLevelError
	case LogLevelWarn:
		c.Level = config.LogLevelWarn
	case LogLevelInfo:
		c.Level = config.LogLevelInfo
	case LogLevelDebug:
		c.Level = config.LogLevelDebug
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
