package controller

import (
	"mals/internal/control/event"
	"mals/internal/control/state"
)

type Controller struct {
	state *state.State
	bus   *event.EventBus
}

func New(state *state.State, bus *event.EventBus) *Controller {
	return &Controller{
		state: state,
		bus:   bus,
	}
}
