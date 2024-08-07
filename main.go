package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	apiStateAddress := &apiState{
		serverHits: 0,
	}
	serverAddress := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	defaultHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", apiStateAddress.HitCounter(defaultHandler))
	mux.HandleFunc("GET /api/healthz", handler)
	mux.HandleFunc("GET /api/metrics", apiStateAddress.Handler)
	mux.HandleFunc("/api/reset", apiStateAddress.Reset)
	serverAddress.ListenAndServe()
}

func handler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}
