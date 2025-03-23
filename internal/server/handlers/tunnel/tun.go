package tunnel

import "io"

type Tun interface {
	io.Writer
}
