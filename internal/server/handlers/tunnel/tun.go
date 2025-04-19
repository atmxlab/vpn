package tunnel

import "io"

//go:generate mock Tun
type Tun interface {
	io.Writer
}
