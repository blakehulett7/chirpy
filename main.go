package main

import (
	"internal/database"
	"net/http"
)

func main() {
	db := &database.Database{
		Path: "./database.json",
	}
	db.EnsureDatabase()
	mux := http.NewServeMux()
	apiStateAddress := &apiState{
		serverHits: 0,
		db:         db,
	}
	serverAddress := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	defaultHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", apiStateAddress.HitCounter(defaultHandler))
	mux.HandleFunc("/api/reset", apiStateAddress.Reset)
	mux.HandleFunc("POST /api/chirps", apiStateAddress.PostChirp)
	mux.HandleFunc("GET /api/chirps", apiStateAddress.GetChirpy)
	mux.HandleFunc("GET /api/chirps/{id}", apiStateAddress.GetaBitChirpy)
	mux.HandleFunc("GET /api/healthz", handler)
	mux.HandleFunc("POST /api/users", apiStateAddress.CreateUser)
	mux.HandleFunc("GET /admin/metrics", apiStateAddress.Metrics)
	serverAddress.ListenAndServe()
}

func handler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}
