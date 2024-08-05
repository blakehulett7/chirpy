package main

import (
	"fmt"
	"net/http"
)

type apiState struct {
	serverHits int
}

func (state *apiState) HitCounter(handler http.Handler) http.Handler {
	state.serverHits++
	return handler
}

func (state *apiState) Handler(writer http.ResponseWriter, request *http.Request) {
	formattedString := fmt.Sprintf("Hits: %v", state.serverHits)
	writer.WriteHeader(200)
	writer.Write([]byte(formattedString))
}

func 
