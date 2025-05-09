package signal

import (
	"context"
	"time"

	"github.com/atmxlab/vpn/pkg/errors"
)

type Signal struct {
	signal chan struct{}
}

func NewSignal() *Signal {
	return &Signal{
		signal: make(chan struct{}, 1),
	}
}

func (s *Signal) Wait() {
	<-s.signal
}

func (s *Signal) WaitWithTimeout(timeout time.Duration) error {
	select {
	case <-time.After(timeout):
		return errors.DeadlineExceeded("waiting signal")
	case <-s.signal:
		return nil
	}
}

func (s *Signal) Close() {
	close(s.signal)
}

func (s *Signal) After(ctx context.Context, callback func(context.Context) error) error {
	s.Wait()
	return callback(ctx)
}

func (s *Signal) Signal(_ context.Context) error {
	select {
	case s.signal <- struct{}{}:
		return nil
	default:
		return errors.AlreadyExists("signal buffer overflow")
	}
}
