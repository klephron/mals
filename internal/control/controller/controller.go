package controller

import (
	"mals/internal/control/event"
	"mals/internal/control/state"
	"mals/internal/control/util"
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

func (s *Controller) Debug(msg string, args ...any) {
	util.Debug(s.state, msg, args...)
}

func (s *Controller) Info(msg string, args ...any) {
	util.Info(s.state, msg, args...)
}

func (s *Controller) Warn(msg string, args ...any) {
	util.Warn(s.state, msg, args...)
}

func (s *Controller) Error(msg string, args ...any) {
	util.Error(s.state, msg, args...)
}
