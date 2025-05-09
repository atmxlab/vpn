package signal

import (
	"context"
	"time"

	"github.com/atmxlab/vpn/pkg/errors"
)

type Signaller struct {
	signal chan struct{}
}

func NewSignaller() *Signaller {
	return &Signaller{
		signal: make(chan struct{}, 1),
	}
}

func (s *Signaller) Wait() {
	<-s.signal
}

func (s *Signaller) WaitWithTimeout(timeout time.Duration) error {
	select {
	case <-time.After(timeout):
		return errors.DeadlineExceeded("waiting signal")
	case <-s.signal:
		return nil
	}
}

func (s *Signaller) Close() {
	close(s.signal)
}

func (s *Signaller) After(ctx context.Context, callback func(context.Context) error) error {
	s.Wait()
	return callback(ctx)
}

func (s *Signaller) Signal(_ context.Context) error {
	select {
	case s.signal <- struct{}{}:
		return nil
	default:
		return errors.AlreadyExists("signal buffer overflow")
	}
}
