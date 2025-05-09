package protocol

import (
	"errors"
	"net"
)

type Payload []byte

func (p Payload) IP() (net.IP, error) {
	if len(p) < 4 {
		return nil, errors.New("payload too short")
	}

	return net.IPv4(p[0], p[1], p[2], p[3]), nil
}

func (p Payload) Len() int {
	return len(p)
}
