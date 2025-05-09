package auth

import (
	"context"
	"net"

	"github.com/atmxlab/vpn/internal/domain/client"
	"github.com/atmxlab/vpn/internal/domain/dto/usecase"
	"github.com/atmxlab/vpn/pkg/errors"
)

type ClientRepository interface {
	ExistsByKey(ctx context.Context, key client.Key) (bool, error)
}

type IPDistributor interface {
	AllocateIP(ctx context.Context) (net.IP, error)
}

type Usecase struct {
	clientRepository ClientRepository
	ipDistributor    IPDistributor
}

func New() *Usecase {
	return &Usecase{}
}

func (u *Usecase) Auth(ctx context.Context, options usecase.AuthOptions) (*usecase.AuthResult, error) {
	exists, err := u.clientRepository.ExistsByKey(ctx, options.Key)
	if err != nil {
		return nil, errors.Wrap(err, "get client by key")
	}
	if !exists {
		return nil, errors.NotFoundf("client not found")
	}

	dedicatedIP, err := u.ipDistributor.AllocateIP(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "acquire ip pool")
	}

	return &usecase.AuthResult{
		DedicatedIP: dedicatedIP,
	}, nil
}
