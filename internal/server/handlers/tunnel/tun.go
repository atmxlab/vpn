package tunnel

type Tun interface {
	Write(data []byte) (int, error)
}
