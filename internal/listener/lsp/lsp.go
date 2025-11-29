package lsp

import (
	"context"
	"errors"
	"fmt"
	"mals/internal/listener"
	"mals/internal/state"
	"net"
)

type ListenerLsp struct {
	listener.Listener
	state     state.State
	addr      string
	listening bool
}

func Type() string {
	return "lsp"
}

func New(state state.State, port int) (*ListenerLsp, error) {
	l := &ListenerLsp{
		state:     state,
		addr:      fmt.Sprintf(":%d", port),
		listening: false,
	}
	return l, nil
}

func (s *ListenerLsp) Type() string {
	return Type()
}

func (s *ListenerLsp) Listen(ctx context.Context) error {
	s.listening = true
	err := s.listen(ctx)
	s.listening = false
	return err
}

func (s *ListenerLsp) Listening() bool {
	return s.listening
}

func (s *ListenerLsp) listen(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)

	if err != nil {
		s.state.Error(fmt.Sprintf("%s: %v", s.logPrefix(), err))
		return err
	}

	defer listener.Close()

	s.state.Info(fmt.Sprintf("%s: listen", s.logPrefix()))

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()

		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				s.state.Info(fmt.Sprintf("%s: closed", s.logPrefix()))
				return nil
			}
			s.state.Warn(fmt.Sprintf("%s: %v", s.logPrefix(), err))
			continue
		}

		s.state.Info(fmt.Sprintf("%s: accepted %v", s.logPrefix(), conn.RemoteAddr()))

		s.state.ListenerAddConn(s, ctx, conn)
	}
}

func (s *ListenerLsp) logPrefix() string {
	return fmt.Sprintf("listener[%s%s]", s.Type(), s.addr)
}
