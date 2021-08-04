package main

import (
	"github.com/Kratos40-sba/ws-chat/internal/handlers"
	"log"
	"net/http"
)

func main() {
	r := routes()
	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
