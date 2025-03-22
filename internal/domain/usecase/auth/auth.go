package auth

import (
	"context"

	"github.com/atmxlab/vpn/internal/domain/dto/usecase"
)

type Usecase struct {
}

func New() *Usecase {
	return &Usecase{}
}

func (receiver Usecase) Auth(ctx context.Context, options usecase.AuthOptions) (*usecase.AuthResult, error) {
	return &usecase.AuthResult{}, nil
}
