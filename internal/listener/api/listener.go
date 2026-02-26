package api

import (
	"context"
	"fmt"
	"mals/internal/plane"
	"net/http"
	"time"
)

type ListenerApi struct {
	name  string
	addr  string
	plane plane.Plane
}

func NewHttp(name string, port int, plane plane.Plane) (*ListenerApi, error) {
	l := &ListenerApi{
		name:  name,
		addr:  fmt.Sprintf(":%d", port),
		plane: plane,
	}
	return l, nil
}

func (s *ListenerApi) Name() string {
	return s.name
}

func (s *ListenerApi) Run(ctx context.Context) error {
	errCh := make(chan error)
	defer close(errCh)

	srv := s.newServer()

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		s.plane.Infof("%s: closed", s.Name())
	}()

	s.plane.Infof("%s: listen", s.Name())

	select {
	case err := <-errCh:
		s.plane.Errorf("%s: %v", s.Name(), err)
		return err

	case <-ctx.Done():
		srvCtx, srvCancel := context.WithTimeout(context.Background(), time.Second)
		defer srvCancel()
		return srv.Shutdown(srvCtx)
	}
}
