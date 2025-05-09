package router

type config struct {
	// Кол-во байт, которые будут читаться из tun интерфейса и туннеля
	bufferSize uint16
	// Размер буфера канала
	tunnelChanSize uint
	// Размер буфера канала
	tunChanSize uint
}
