package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	mux.Handle("/", http.FileServer(http.Dir(".")))
	server.ListenAndServe()
}
