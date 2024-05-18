package server

import (
	"net/http"
	"slices"
	"time"
)

var TextStreamHTTPHandler = func(buf []byte, sleepTime time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "localhost:8080")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)	

		for {
			w.Write(slices.Concat([]byte("data:"), buf, []byte("\n\n")))
			w.(http.Flusher).Flush()
			time.Sleep(sleepTime)
		}
	})
}
