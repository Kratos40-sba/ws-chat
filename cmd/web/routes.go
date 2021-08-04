package main

import (
	"github.com/Kratos40-sba/ws-chat/internal/handlers"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func routes() http.Handler {
	r := httprouter.New()
	r.Handler(http.MethodGet, "/", http.HandlerFunc(handlers.Home))
	r.Handler(http.MethodGet, "/ws", http.HandlerFunc(handlers.WsEndpoint))
	return r
}
