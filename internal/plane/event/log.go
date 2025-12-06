package event

import "mals/internal/log"

type EventLog struct {
	EventGeneric
	Level   log.Level
	Pattern string
	Args    []any
}
