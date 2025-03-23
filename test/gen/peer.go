package gen

import "github.com/atmxlab/vpn/internal/server"

func RandPeer() *server.Peer {
	return server.NewPeer(
		RandIP(),
		RandAddr(),
	)
}
