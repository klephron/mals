package state

import (
	listener "mals/internal/listener"
)

func (s *StateImpl) EventLoop() {
	for {
		event, ok := <-s.EventChan
		if !ok {
			return
		}

		switch event := event.(type) {

		case *EventListenerDown:
			listening := false
			s.Listeners.Range(func(listener listener.Listener, value ListenerState) bool {
				if listener.Listening() {
					listening = true
					return false
				}
				return true
			})
			if !listening {
				return
			}

		case *EventClientListen:
			go event.client.Listen(event.ctx)
		}

	}
}
