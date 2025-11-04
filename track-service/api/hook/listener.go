package hook

import (
	"encoding/json"
	"fmt"
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
		res := make(map[string]string)
		json.NewDecoder(r.Body).Decode(&res)
		fmt.Println(res)
		w.WriteHeader(200)
		w.Write([]byte("MoyaPosylkaWebhook"))
	})

	log.Println("Hook is listening on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Println(err)
	}
}
