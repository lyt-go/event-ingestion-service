package main

import (
	"eventingestion/internal/app"
	"eventingestion/internal/handler"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, handler.New(app.NewService())); err != nil {
		log.Fatal(err)
	}
}
