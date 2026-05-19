package server

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	http *http.Server
}

func New() (*Server, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{"status": "ok"}
		res, _ := json.Marshal(data)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	return &Server{
		http: server,
	}, nil
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}
