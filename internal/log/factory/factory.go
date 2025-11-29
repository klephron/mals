package log

import (
	"fmt"
	"mals/internal/log"
	"mals/internal/log/file"
	"mals/pkg/config"
)

func OpenConfig(logConfig config.Log) (log.Log, error) {
	switch t := logConfig.(type) {
	case *config.LogFile:
		opened, err := file.Open(t.File, t.Level)
		if err != nil {
			return nil, err
		}
		return opened, nil
	default:
		return nil, fmt.Errorf("unhandled log type %T", t)
	}
}
