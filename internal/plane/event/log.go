package event

import "mals/internal/log"

type EventLog struct {
	EventGeneric
	Level log.Level
	Msg   string
}
