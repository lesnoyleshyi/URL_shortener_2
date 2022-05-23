package handlers

import (
	"URL_shortener_2/internal/services"
	"net/http"
)

type Handler interface {
	http.Handler
}

type handler struct {
	service services.Service
}

func New(s services.Service) Handler {
	return &handler{service: s}

}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		processShort()
	case http.MethodGet:
		processLong()
	}
}

func processShort() {

}

func processLong() {

}
