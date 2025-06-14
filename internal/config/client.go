package config

import "net"

type ClientConfig struct {
	// Кол-во байт, которые будут читаться из tun интерфейса и туннеля
	BufferSize uint16 `json:"bufferSize,omitempty"`
	// Время жизни пира до следующего keepalive запроса
	KeepAliveTickDuration Duration `json:"keepAliveTickDuration,omitempty"`
	// Шлюз для выхода в интернет (e.g. Адрес роутера)
	GatewayIP string       `json:"gatewayIP,omitempty"`
	Tun       ClientTun    `json:"tun"`
	Tunnel    ClientTunnel `json:"tunnel"`
}

func (cc ClientConfig) GetGatewayIP() net.IP {
	return net.ParseIP(cc.GatewayIP)
}

type ClientTun struct {
	// Маска подсети для адреса tun интерфейса
	IPMask string `json:"ipMask,omitempty"`
	// Maximum Transition Unit  - максимальная длина неделимого пакета
	MTU uint16 `json:"mtu,omitempty"`
	// Размер буфера канала,
	// в который будут складываться пакеты из TUN интерфейса
	TunChanSize uint `json:"tunChanSize,omitempty"`
}

type ClientTunnel struct {
	// Размер буфера канала,
	// в который будут складываться пакеты из тоннеля
	TunnelChanSize uint `json:"tunnelChanSize,omitempty"`
	// Протокол, в котором будет работать тоннель
	Network string `json:"network,omitempty"`
	// IP адрес, который будет слушать тоннель
	IP string `json:"ip,omitempty"`
	// Port, который будет слушать тоннель
	Port uint16 `json:"port,omitempty"`
	// IP адрес сервера, к которому подключаемся
	ServerIP string `json:"serverIP,omitempty"`
	// Port сервера, к которому подключаемся
	ServerPort uint16 `json:"serverPort,omitempty"`
	// Timeout соединения к серверу
	ServerConnectionTimeout Duration `json:"serverConnectionTimeout,omitempty"`
}

func (ct ClientTunnel) GetServerIP() net.IP {
	return net.ParseIP(ct.ServerIP)
}
