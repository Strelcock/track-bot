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
			w.Write([]byte("read body: " + err.Error()))
			return
		}
		log.Println(string(rawData))
		// if len(rawData) == 0 {
		// 	next.ServeHTTP(w, r)
		// }
		//log.Println(rawData)

		rawMap := make(map[string]any)
		err = json.Unmarshal(rawData, &rawMap)
		log.Println(rawMap)
		if err != nil {
			log.Printf("unmarshal into map: %s\n", err.Error())
			w.WriteHeader(500)
			w.Write([]byte("unmarshal into map: " + err.Error()))
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
