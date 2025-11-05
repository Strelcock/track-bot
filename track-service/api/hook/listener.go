package hook

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Listener struct {
}

func NewListener() *Listener {
	return &Listener{}
}

func (l *Listener) ListenAndServe(r *chi.Mux) {
	r.Post("/hook/listen", func(w http.ResponseWriter, r *http.Request) {
		res := make(map[string]struct{})
		err := json.NewDecoder(r.Body).Decode(&res)
		if err != nil {
			log.Println("first-step decoding: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, ok := res["newEvents"]; ok {
			var newEvent Event
			err = json.NewDecoder(r.Body).Decode(&newEvent)
			if err != nil {
				log.Println("new event decoding: ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Println(newEvent)
		} else if _, ok := res["trackerDelivered"]; ok {
			var trackDelivered Delivered
			err = json.NewDecoder(r.Body).Decode(&trackDelivered)
			if err != nil {
				log.Println("tracker delivered decoding: ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Println(trackDelivered)
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
