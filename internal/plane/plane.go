package plane

import (
	"mals/internal/listener"
	"mals/internal/log"
	"mals/internal/plane/controller/factory"
	"mals/internal/plane/controller/lifecycle"
	"mals/internal/plane/controller/logging"
	"mals/internal/plane/controller/ownership"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type Plane struct {
	state     *state.State
	bus       *event.EventBus
	Factory   *factory.FactoryController
	Lifecycle *lifecycle.LifecycleController
	Logging   *logging.LogController
	Ownership *ownership.OwnershipController
}

func New() *Plane {
	manager := &Plane{
		state: &state.State{
			Listeners: xsync.NewMap[listener.Listener, *state.ListenerValue](),
			Logs:      xsync.NewMap[log.Log, *state.LogValue](),
		},
		bus: event.NewEventBus(),
	}

	manager.Factory = factory.NewController(manager.state, manager.bus)
	manager.Lifecycle = lifecycle.NewController(manager.state, manager.bus)
	manager.Logging = logging.NewController(manager.state, manager.bus)
	manager.Ownership = ownership.NewController(manager.state, manager.bus)

	return manager
}

func (s *Plane) Serve() {
	var wg sync.WaitGroup

	wg.Go(func() {
		s.Factory.Serve()
	})
	wg.Go(func() {
		s.Lifecycle.Serve()
	})
	wg.Go(func() {
		s.Logging.Serve()
	})
	wg.Go(func() {
		s.Ownership.Serve()
	})

	wg.Wait()
}
