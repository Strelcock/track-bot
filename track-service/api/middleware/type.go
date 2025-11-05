package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const (
	newEvents        = "newEvents"
	trackerDelivered = "trackerDelivered"
)

func WebhookEvent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawData := []byte{}
		err := json.NewDecoder(r.Body).Decode(&rawData)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		rawMap := make(map[string]any)
		err = json.Unmarshal(rawData, &rawMap)
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
