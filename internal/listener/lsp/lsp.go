package lsp

import (
	"context"
	"errors"
	"fmt"
	"mals/internal/listener/common"
	"mals/internal/state"
	"net"
)

type ListenerLsp struct {
	common.Listener
	state     *state.State
	addr      string
	listening bool
}

func Type() string {
	return "lsp"
}

func New(state *state.State, port int) (*ListenerLsp, error) {
	l := &ListenerLsp{
		state:     state,
		addr:      fmt.Sprintf(":%d", port),
		listening: false,
	}
	return l, nil
}

func (l *ListenerLsp) Type() string {
	return Type()
}

func (l *ListenerLsp) ListenAndServe(ctx context.Context) error {
	l.listening = true
	err := l.listen(ctx)
	l.listening = false
	return err
}

func (l *ListenerLsp) Listening() bool {
	return l.listening
}

func (l *ListenerLsp) listen(ctx context.Context) error {
	listener, err := net.Listen("tcp", l.addr)

	if err != nil {
		l.state.Error(fmt.Sprintf("%s: %v", l.logPrefix(), err))
		return err
	}

	defer listener.Close()

	l.state.Info(fmt.Sprintf("%s: listen", l.logPrefix()))

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()

		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				l.state.Info(fmt.Sprintf("%s: closed", l.logPrefix()))
				return nil
			}
			l.state.Warn(fmt.Sprintf("%s: %v", l.logPrefix(), err))
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			l.state.Info(fmt.Sprintf("%s: connection %v established", l.logPrefix(), c))
		}(conn)
	}
}

func (l *ListenerLsp) logPrefix() string {
	return fmt.Sprintf("listener[%s%s]", l.Type(), l.addr)
}
