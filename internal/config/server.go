package config

import (
	"net"
	"time"
)

type ServerConfig struct {
	// Адрес сервера - тоннеля - ip:port на котором будет стартовать сервер
	ServerAddr net.Addr
	// Кол-во байт, которые будут читаться из tun интерфейса и туннеля
	BufferSize uint16
	// Время жизни пира до следующего keepalive запроса
	PeerKeepAliveTTL time.Duration

	Tun struct {
		// 10.1.1.0/24 Подсеть TUN интерфейса
		// Из этой подсети будут выдаваться IP адреса клиентам
		Subnet net.IPNet
		// Maximum Transition Unit  - максимальная длина неделимого пакета
		MTU uint16
		// Размер буфера канала,
		// в который будут складываться пакеты из тоннеля
		TunnelChanSize uint
		// Размер буфера канала,
		// в который будут складываться пакеты из TUN интерфейса
		TunChanSize uint
	}
}
