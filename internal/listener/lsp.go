package listener

import (
	"context"
	"errors"
	"fmt"
	"mals/internal/log"
	"net"
)

type ListenerLsp struct {
	Listener
	log  *log.Log
	addr string
}

func listenerLspType() string {
	return "lsp"
}

func newListenerLsp(port int, log *log.Log) (*ListenerLsp, error) {
	l := &ListenerLsp{
		log:  log,
		addr: fmt.Sprintf(":%d", port),
	}
	return l, nil
}

func (l *ListenerLsp) logPrefix() string {
	return fmt.Sprintf("listener[%s, %s]", l.Type(), l.addr)
}

func (l *ListenerLsp) Type() string {
	return listenerLspType()
}

func (l *ListenerLsp) ListenAndServe(ctx context.Context) error {
	listener, err := net.Listen("tcp", l.addr)
	if err != nil {
		l.log.Error("%s: %v", l.logPrefix(), err)
		return err
	}
	defer listener.Close()

	l.log.Info("%s: listen", l.logPrefix())

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()

		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				l.log.Info("%s: closed", l.logPrefix())
				return nil
			}
			l.log.Warn("%s: %v", l.logPrefix(), err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			l.log.Info("%s: connection %v established", l.logPrefix(), c)
		}(conn)
	}
}
