package main

import (
	"log"
	"tracker/api/hook"
	"tracker/api/queue"
	"tracker/api/server"
	"tracker/config"
	"tracker/internal/infrastructure/kafka"
	"tracker/internal/usecase/queueservice"

	"github.com/go-chi/chi/v5"
)

const port = ":50052"

const trackStatusTopic = "track.status.get"

func main() {
	cfg := config.Load()
	//log.Println(cfg)
	prod := kafka.New([]string{cfg.Broker}, trackStatusTopic)
	srv := queueservice.New(prod)
	queue := queue.New(srv)
	sender := hook.NewSender(cfg.ApiKey)

	s := server.New(sender)

	listener := hook.NewListener(queue)
	r := chi.NewRouter()
	go listener.ListenAndServe(r)
	err := s.Listen(port)
	if err != nil {
		log.Fatal(err)
	}
}
