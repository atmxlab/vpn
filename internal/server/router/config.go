package router

import (
	"net"
)

type config struct {
	// Кол-во байт, которые будут читаться из tun интерфейса и туннеля
	bufferSize uint16

	tun struct {
		// 10.1.1.0/24 Подсеть TUN интерфейса
		// Из этой подсети будут выдаваться IP адреса клиентам
		subnet net.IPNet
		// Maximum Transition Unit  - максимальная длина неделимого пакета
		mtu uint16
	}

	// Размер буфера канала
	tunnelChanSize uint
	// Размер буфера канала
	tunChanSize uint
}
