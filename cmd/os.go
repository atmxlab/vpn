package cmd

import (
	"context"
	"os/signal"
	"syscall"
)

func SignalCtx() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
}
