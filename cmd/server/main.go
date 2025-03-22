package server

import (
	"github.com/atmxlab/vpn/internal/domain/usecase/auth"
	"github.com/atmxlab/vpn/internal/http"
	"github.com/atmxlab/vpn/internal/http/handler"
	"github.com/atmxlab/vpn/pkg/errors"
)

func main() {
	authHandler := handler.NewAuth(auth.New())

	httpServer := http.New(authHandler)

	err := httpServer.ListenAndServe(":8080")
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		errors.Fatalf(err.Error())
	}
}
