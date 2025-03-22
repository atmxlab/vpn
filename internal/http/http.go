package http

import (
	"github.com/atmxlab/vpn/pkg/errors"
	"net/http"
)

// ResponseWriter need for mock
type ResponseWriter interface {
	http.ResponseWriter
}

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
