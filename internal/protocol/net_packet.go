package protocol

// TunPacket - пакет, который передается по сети (не тоннель)
type TunPacket struct {
	payload Payload
}
