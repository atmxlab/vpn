package test

import (
	"encoding/binary"
)

func IPHeaderChecksum(b []byte) uint16 {
	csum := uint32(0)
	for i := 0; i < len(b); i += 2 {
		csum += uint32(binary.BigEndian.Uint16(b[i:min(i+2, len(b))]))
	}
	for csum > 0xffff {
		csum = (csum >> 16) + (csum & 0xffff)
	}
	return ^uint16(csum)
}
