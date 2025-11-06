package hook

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"tracker/api/middleware"
	"tracker/api/queue"

	"github.com/go-chi/chi/v5"
)

const (
	newEvents        = "newEvents"
	trackerDelivered = "trackerDelivered"
)

type Listener struct {
	q *queue.Queue
}

func NewListener(q *queue.Queue) *Listener {
	return &Listener{q}
}

func (l *Listener) ListenAndServe(r *chi.Mux) {
	r.Use(middleware.WebhookEvent)
	r.Post("/hook/listen", func(w http.ResponseWriter, r *http.Request) {
		eventType := r.Header.Get("Event-Type")
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		switch eventType {

		case newEvents:
			var newEvent NewEvent
			err := json.NewDecoder(r.Body).Decode(&newEvent)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			status := fmt.Sprintf("%s: %s", newEvent.NewEvents.Events[0].Operation, newEvent.NewEvents.Events[0].Attribute)

			msg := &Message{
				Event:   eventType,
				Barcode: newEvent.NewEvents.Barcode,
				Status:  status,
			}

			byteMsg, err := json.Marshal(msg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = l.q.WriteMessages(ctx, byteMsg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Println(newEvent)

		case trackerDelivered:
			var trackerDelivered Delivered
			err := json.NewDecoder(r.Body).Decode(&trackerDelivered)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			msg := &Message{
				Event:   eventType,
				Barcode: trackerDelivered.TrackerDelivered.Barcode,
			}

			byteMsg, err := json.Marshal(msg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = l.q.WriteMessages(ctx, byteMsg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(200)
		w.Write([]byte("MoyaPosylkaWebhook"))
	})

	log.Println("Hook is listening on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Println(err)
	}
}
