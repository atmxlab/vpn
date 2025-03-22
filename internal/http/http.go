package http

import (
	"net/http"

	"github.com/atmxlab/vpn/pkg/errors"
)

// ResponseWriter need for mock
//
//go:generate mock ResponseWriter
type ResponseWriter interface {
	http.ResponseWriter
}

//go:generate mock Handler
type Handler interface {
	Pattern() string
	Handle(writer http.ResponseWriter, request *http.Request)
}

type HTTP struct {
	handlers []Handler
}

func New(handlers ...Handler) *HTTP {
	return &HTTP{handlers: handlers}
}

func (h *HTTP) ListenAndServe(addr string) error {
	for _, handler := range h.handlers {
		http.HandleFunc(handler.Pattern(), handler.Handle)
	}

	if err := http.ListenAndServe(addr, nil); err != nil {
		return errors.Wrap(err, "http.ListenAndServe")
	}

	return nil
}
