package main

import (
	"fmt"
	"net/http"
)

type apiState struct {
	serverHits int
}

func (state *apiState) HitCounter(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		state.serverHits++
		handler.ServeHTTP(writer, request)
	})
}

func (state *apiState) Handler(writer http.ResponseWriter, request *http.Request) {
	formattedString := fmt.Sprintf("Hits: %v", state.serverHits)
	writer.WriteHeader(200)
	writer.Write([]byte(formattedString))
}

func (state *apiState) Reset(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	state.serverHits = 0
}
