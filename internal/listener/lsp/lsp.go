package lsp

import (
	"context"
	"errors"
	"fmt"
	"mals/internal/control/controller"
	"mals/internal/listener"
	"net"
)

type ListenerLsp struct {
	listener.Listener
	controller *controller.Controller
	addr       string
	listening  bool
}

func Type() string {
	return "lsp"
}

func New(controller *controller.Controller, port int) (*ListenerLsp, error) {
	l := &ListenerLsp{
		controller: controller,
		addr:       fmt.Sprintf(":%d", port),
		listening:  false,
	}
	return l, nil
}

func (s *ListenerLsp) Type() string {
	return Type()
}

func (s *ListenerLsp) Listening() bool {
	return s.listening
}

func (s *ListenerLsp) Listen(ctx context.Context) error {
	s.listening = true
	defer func() {
		s.listening = false
	}()

	listener, err := net.Listen("tcp", s.addr)

	if err != nil {
		s.controller.Error(fmt.Sprintf("%s: %v", s.logPrefix(), err))
		return err
	}

	defer listener.Close()

	s.controller.Info(fmt.Sprintf("%s: listen", s.logPrefix()))

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()

		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				s.controller.Info(fmt.Sprintf("%s: closed", s.logPrefix()))
				return nil
			}
			s.controller.Warn(fmt.Sprintf("%s: %v", s.logPrefix(), err))
			continue
		}

		s.controller.Info(fmt.Sprintf("%s: accepted %v", s.logPrefix(), conn.RemoteAddr()))

		s.controller.Warn("TODO: listener connect")
	}
}

func (s *ListenerLsp) logPrefix() string {
	return fmt.Sprintf("listener[%s%s]", s.Type(), s.addr)
}
