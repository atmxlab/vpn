package protocol

const (
	// HeaderSize - длина заголовка
	HeaderSize = 1
)

// Header заголовок пакета
// Служебные данные в пакетах передаются через такой заголовок
type Header struct {
	flag Flag
}

func NewHeader(flag Flag) *Header {
	return &Header{
		flag: flag,
	}
}

func (h *Header) Flag() Flag {
	return h.flag
}
