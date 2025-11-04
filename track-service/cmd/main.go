package main

import (
	"log"
	"tracker/api/hook"
	"tracker/api/server"

	"github.com/go-chi/chi/v5"
)

const port = ":50052"

func main() {
	s := server.New()

	listener := hook.NewListener()
	r := chi.NewRouter()
	go listener.ListenAndServe(r)
	err := s.Listen(":50052")
	if err != nil {
		log.Fatal(err)
	}
}
