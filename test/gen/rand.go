package gen

import (
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

func RandString() string {
	return uuid.NewString()
}

func RandInt64() int64 {
	return rand.Int63()
}

func RandInt32() int32 {
	return rand.Int31()
}

func RandInt() int {
	return rand.Intn(math.MaxInt16)
}

func RandIntInRange(from, to int) int {
	return rand.Intn(to-from+1) + from
}

func RandUInt64() uint64 {
	return rand.Uint64()
}

func RandUInt16() uint16 {
	return uint16(rand.Uint32() % math.MaxUint16)
}

func RandUInt8() uint8 {
	return uint8(rand.Uint32() % math.MaxUint8)
}

func RandByte() byte {
	return RandUInt8()
}

func RandDuration() time.Duration {
	return time.Duration(rand.Intn(1000))
}

// RandElement возвращает рандомный элемент.
func RandElement[T any](elements ...T) T {
	if len(elements) == 0 {
		return lo.Empty[T]()
	}

	return elements[RandIntInRange(0, len(elements)-1)]
}
