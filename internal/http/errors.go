package http

import (
	"log"
	"net/http"
)

func ServerError(errOrigin error, w http.ResponseWriter) {
	log.Printf("server error: %s", errOrigin.Error())

	w.WriteHeader(500)
	_, err := w.Write([]byte("server error"))
	if err != nil {
		log.Printf("server error: error write response: %s", err.Error())
	}

	return
}

func BadRequestError(w http.ResponseWriter) {
	w.WriteHeader(404)

	_, err := w.Write([]byte("bad request"))
	if err != nil {
		ServerError(err, w)
	}

	return
}
