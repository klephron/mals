package plane

import (
	"mals/internal/plane/controller"
	"mals/internal/plane/controller/lifecycle"
	"mals/internal/plane/controller/logging"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type Plane struct {
	state     *state.State
	bus       *event.EventBus
	Lifecycle controller.LifecycleController
	Log       controller.LogController
}

func New() *Plane {
	manager := &Plane{
		state: &state.State{
			Logs:      xsync.NewMap[string, *state.LogValue](),
			Listeners: xsync.NewMap[string, *state.ListenerValue](),
		},
		bus: event.NewEventBus(),
	}

	manager.Lifecycle = lifecycle.New(manager.state, manager.bus)
	manager.Log = logging.New(manager.state, manager.bus)

	return manager
}

func (s *Plane) Serve(onReady func()) {
	var wg sync.WaitGroup
	var wgReady sync.WaitGroup

	wgReady.Add(2)

	wg.Go(func() {
		s.Lifecycle.Serve(func() { wgReady.Done() })
	})
	wg.Go(func() {
		s.Log.Serve(func() { wgReady.Done() })
	})

	wgReady.Wait()
	onReady()

	wg.Wait()
}
