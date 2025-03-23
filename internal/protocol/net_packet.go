package protocol

// TunPacket - пакет, который передается по сети (не тоннель)
type TunPacket struct {
	payload Payload
}

func NewTunPacket(payload Payload) *TunPacket {
	return &TunPacket{payload: payload}
}

func (t TunPacket) Payload() Payload {
	return t.payload
}
