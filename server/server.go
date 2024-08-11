package server

import (
	"github.com/isaquecsilva/diskstatus/templ"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

func CreateNewServer(addr string, logStream io.Writer, handler http.Handler) *http.Server {
	logHandler := slog.NewJSONHandler(logStream, nil)
	logger := slog.NewLogLogger(logHandler, slog.LevelError)

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, strings.NewReader(templ.ChartTemplate))
	})

	http.Handle("GET /data", handler)

	return &http.Server{
		Addr:     addr,
		Handler:  nil,
		ErrorLog: logger,
	}
}
