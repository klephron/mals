package state

import (
	"context"
	listener "mals/internal/listener/common"
	"sync"
)

func (s *State) ListenerAdd(listener listener.Listener) {
	s.Listeners = append(s.Listeners, listener)
}

func (s *State) ListenerListenAndServeSnapshot(ctx context.Context) {
	var wg sync.WaitGroup

	for _, l := range s.Listeners {
		if l.Listening() {
			continue
		}

		wg.Add(1)
		go func(l listener.Listener) {
			defer wg.Done()
			err := l.ListenAndServe(ctx)
			if err != nil {
				s.LogContext().Error(err.Error())
			}
		}(l)
	}

	wg.Wait()
}
