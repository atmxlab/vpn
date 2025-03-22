package usecase

import "net"

type AuthOptions struct {
	Key []byte
	IP  net.IP
}

type AuthResult struct {
	DedicatedIP net.IP
}
