package ip

import (
	"math"
	"net"
)

func CountInMask(mask net.IPMask) (count int) {
	ones, bits := mask.Size()

	return int(math.Pow(2, float64(bits-ones)))
}
