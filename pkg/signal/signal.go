package signal

import (
	"context"
	"time"
)

type Signal struct {
	signal chan struct{}
}

func NewSignal() *Signal {
	return &Signal{
		signal: make(chan struct{}),
	}
}

func (s *Signal) Wait() {
	<-s.signal
}

func (s *Signal) WaitWithTimeout(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
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

func (s *Signal) Signal() {
	s.signal <- struct{}{}
}
