package tun

import (
	"context"
	"io"
	"sync"

	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

// TunIface - виртуальный сетевой интерфейс
// В Linux можно создавать два интерфейса - TUN и TAP
// Нам необходим только TUN, так как работаем на L3 уровне, а не на L2
type TunIface interface {
	io.ReadWriteCloser
	Name() string
}

type Tun struct {
	tun TunIface
}

func NewTun(tun TunIface) *Tun {
	return &Tun{tun: tun}
}

// ReadWithContext - необходим, чтобы учитывать отмену контекста при чтении из потока
func (t *Tun) ReadWithContext(ctx context.Context, buf []byte) (int, error) {
	type result struct {
		n   int
		err error
	}

	resultChan := make(chan result, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(resultChan)

		n, err := t.tun.Read(buf)
		resultChan <- result{n, err}
	}()

	select {
	case <-ctx.Done():
		logrus.Warnf("Context canceled: %v", ctx.Err())
		if err := t.tun.Close(); err != nil {
			return 0, errors.Join(err, ctx.Err())
		}

		logrus.Warn("Waiting ending read from Tunnel...")
		wg.Wait()

		return 0, ctx.Err()
	case res := <-resultChan:
		return res.n, res.err
	}
}

func (t *Tun) Close() error {
	return t.tun.Close()
}
