package event

import "mals/internal/log"

type EventLogAdd struct {
	EventGeneric
	Log log.Log
}

type EventLogDelete struct {
	EventGeneric
	Log log.Log
}

type EventLogStart struct {
	EventGeneric
	Log log.Log
}

type EventLogStop struct {
	EventGeneric
	Log log.Log
}
