package protocol

// Flag - флаг передается в заголовке пакета
// Флаги работают как ручки для сервера и клиента
// В зависимости от флага, сервер или клиент по-разному обрабатывают пакеты
type Flag byte

const (
	// FlagSYN - synchronize - установление соединения
	FlagSYN Flag = iota
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
