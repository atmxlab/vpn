package protocol

// NetPacket - пакет, который передается по сети (не тоннель)
type NetPacket struct {
	payload Payload
}
