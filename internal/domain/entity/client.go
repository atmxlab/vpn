package entity

import "github.com/atmxlab/vpn/internal/domain/client"

type Client struct {
	id  client.ID
	key client.Key
}

func (c Client) ID() client.ID {
	return c.id
}

func (c Client) Key() client.Key {
	return c.key
}
