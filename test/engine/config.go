package engine

import (
	"time"

	"github.com/atmxlab/vpn/internal/config"
)

type Config struct {
	// Позволяет дожидаться доставки данных по каналам и обработки данных
	actionDelay  time.Duration
	serverConfig *config.ServerConfig
}

func WithActionDelay(dur time.Duration) func(c *Config) {
	return func(c *Config) {
		c.actionDelay = dur
	}
}

func WithPeerKeepAliveTTL(ttl time.Duration) func(c *Config) {
	return func(c *Config) {
		c.serverConfig.PeerKeepAliveTTL = config.Duration(ttl)
	}
}
