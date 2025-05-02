package engine

import (
	"time"

	"github.com/atmxlab/vpn/test"
)

func WAIT(dur time.Duration) test.Action {
	return newSimpleAction(func(_ test.App) {
		time.Sleep(dur)
	})
}
