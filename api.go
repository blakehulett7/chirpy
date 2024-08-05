package main

import "net/http"

type apiState struct {
	serverHits int
}

func (state *apiState) HitCounter(handler http.Handler) http.Handler {
	state.serverHits++
	return handler
}
