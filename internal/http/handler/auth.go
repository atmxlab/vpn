package handler

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"

	"github.com/atmxlab/vpn/internal/domain/dto/usecase"
	hhttp "github.com/atmxlab/vpn/internal/http"
)

//go:generate mock Usecase
type Usecase interface {
	Auth(ctx context.Context, options usecase.AuthOptions) (*usecase.AuthResult, error)
}

type Auth struct {
	usecase Usecase
}

func NewAuth(usecase Usecase) *Auth {
	return &Auth{usecase: usecase}
}

type AuthRequest struct {
	IP  string `json:"ip"`
	Key string `json:"key"`
}

func (a *Auth) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// TODO: нужно проверить как это отработает если послать сюда ГБ, например
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		hhttp.BadRequestError(w)

		return
	}

	var req AuthRequest
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		hhttp.ServerError(err, w)
		return
	}

	res, err := a.usecase.Auth(ctx, usecase.AuthOptions{
		Key: []byte(req.Key),
		IP:  net.ParseIP(req.IP),
	})

	bytes, err = json.Marshal(res)
	if err != nil {
		hhttp.ServerError(err, w)
	}

	w.WriteHeader(200)
	_, err = w.Write(bytes)
	if err != nil {
		hhttp.ServerError(err, w)
	}
}

func (a *Auth) Pattern() string {
	return "GET auth"
}
