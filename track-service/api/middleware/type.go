package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
	newEvents        = "newEvents"
	trackerDelivered = "trackerDelivered"
)

// check which event has been delivered
func WebhookEvent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawData, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		log.Println(rawData)
		rawMap := make(map[string]map[string]any)
		err = json.Unmarshal(rawData, &rawMap)
		log.Println(rawMap)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		for k := range rawMap {
			switch k {
			case newEvents:
				r.Header.Add("Event-Type", newEvents)
				r.Body = io.NopCloser(bytes.NewBuffer(rawData))
			case trackerDelivered:
				r.Header.Add("Event-Type", trackerDelivered)
				r.Body = io.NopCloser(bytes.NewBuffer(rawData))
			}
		}

		next.ServeHTTP(w, r)
	})
}
