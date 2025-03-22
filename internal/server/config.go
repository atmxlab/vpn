package server

import (
	"net"
	"time"
)

type Config struct {
	// Адрес сервера - тоннеля - ip:port на котором будет стартовать сервер
	ServerAddr net.Addr
	// Кол-во байт, которые будут читаться из tun интерфейса и туннеля
	BufferSize uint16
	// Таймаут сохранения соединения после последнего получения keepalive сообщения
	PeerKeepAliveMissingTimeout time.Duration

	Tun struct {
		// 10.1.1.0/24 Подсеть TUN интерфейса
		// Из этой подсети будут выдаваться IP адреса клиентам
		Subnet net.IPNet
		// Maximum Transition Unit  - максимальная длина неделимого пакета
		MTU uint16
	}
}
