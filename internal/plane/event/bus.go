package event

import "sync"

type EventBus struct {
	mu          sync.RWMutex
	subscribers []chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make([]chan Event, 0),
	}
}

func (s *EventBus) Unicast(event Event, dst <-chan Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, subscriber := range s.subscribers {
		if dst == subscriber {
			subscriber <- event
			break
		}
	}
}

func (s *EventBus) Broadcast(event Event, src <-chan Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, subscriber := range s.subscribers {
		if src != subscriber {
			subscriber <- event
		}
	}
}

func (s *EventBus) Allcast(event Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, subscriber := range s.subscribers {
		subscriber <- event
	}
}

func (s *EventBus) Subscribe() <-chan Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan Event)
	s.subscribers = append(s.subscribers, ch)
	return ch
}

func (s *EventBus) Unsubscribe(ch <-chan Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, subscriber := range s.subscribers {
		if ch == subscriber {
			s.subscribers = append(s.subscribers[:i], s.subscribers[i+1:]...)
			close(subscriber)
			// drain
			for range ch {
			}
			return
		}
	}
}

func (s *EventBus) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, subscriber := range s.subscribers {
		close(subscriber)
	}
}
