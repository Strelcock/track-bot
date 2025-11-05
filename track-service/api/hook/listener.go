package hook

import (
	"encoding/json"
	"log"
	"net/http"
	"tracker/api/middleware"

	"github.com/go-chi/chi/v5"
)

const (
	newEvents        = "newEvents"
	trackerDelivered = "trackerDelivered"
)

type Listener struct {
}

func NewListener() *Listener {
	return &Listener{}
}

func (l *Listener) ListenAndServe(r *chi.Mux) {
	r.Use(middleware.WebhookEvent)
	r.Post("/hook/listen", func(w http.ResponseWriter, r *http.Request) {
		eventType := r.Header.Get("Event-Type")

		switch eventType {

		case newEvents:
			var newEvent Event
			err := json.NewDecoder(r.Body).Decode(&newEvent)
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
			log.Println(trackerDelivered)
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
