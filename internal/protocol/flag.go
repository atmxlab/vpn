package protocol

// Flag - флаг передается в заголовке пакета
// Флаги работают как ручки для сервера и клиента
// В зависимости от флага, сервер или клиент по-разному обрабатывают пакеты
type Flag byte

const (
	// FlagUNK - unknown - неизвестный флаг
	FlagUNK Flag = iota
	// FlagSYN - synchronize - установление соединения
	FlagSYN
	// FlagACK - acknowledgement - подтверждение соединения
	FlagACK
	// FlagFIN - finish -завершение соединения
	FlagFIN
	// FlagPSH -push - передача данных
	FlagPSH
	// FlagKPA - keepalive - поддержка соединения
	FlagKPA
)

func (f Flag) Is(ff Flag) bool {
	return f == ff
}

func (f Flag) Byte() byte {
	return byte(f)
}

func (f Flag) String() string {
	switch f {
	case FlagSYN:
		return "synchronize"
	case FlagACK:
		return "acknowledgement"
	case FlagFIN:
		return "finish"
	case FlagPSH:
		return "push"
	case FlagKPA:
		return "keepalive"
	default:
		return "Unknown"
	}
}

func Flags() []Flag {
	return []Flag{FlagSYN, FlagACK, FlagFIN, FlagPSH, FlagKPA}
}
