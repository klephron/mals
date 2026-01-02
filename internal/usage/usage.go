package usage

import "mals/pkg/config"

type ConditionFilter struct {
	Filetype *string
	Path     *string
}

type EventFilter struct {
	Event *config.Event
}
