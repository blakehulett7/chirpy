package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	apiStatePointer := &apiState{}
	serverPointer := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	defaultHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", apiStatePointer.HitCounter(defaultHandler))
	mux.HandleFunc("/healthz", handler)
	mux.HandleFunc("/metrics", apiStatePointer.Handler)
	serverPointer.ListenAndServe()
}

func handler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}
