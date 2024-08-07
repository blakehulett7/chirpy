package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

func (state *apiState) Metrics(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/html")
	writer.WriteHeader(200)
	formattedHTML := fmt.Sprintf(metricsHTML, state.serverHits)
	writer.Write([]byte(formattedHTML))
}

func (state *apiState) Reset(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	state.serverHits = 0
}

func (state *apiState) Validate(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnError struct {
		Error string `json:"error"`
	}
	type returnCleaned struct {
		CleanedBody string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	error := decoder.Decode(&params)
	if error != nil {
		responseBody := returnError{
			Error: "something went wrong",
		}
		responseData, _ := json.Marshal(responseBody)
		JsonResponse(responseData, writer, 500)
		return
	}
	if len(params.Body) > 140 {
		responseBody := returnError{
			Error: "chirp is too long",
		}
		responseData, _ := json.Marshal(responseBody)
		JsonResponse(responseData, writer, 400)
		return
	}
	normalizedBody := strings.ToLower(params.Body)
	profanities := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	for _, profanity := range profanities {
		if strings.Contains(normalizedBody, profanity) {
			words := strings.Split(params.Body, " ")
			cleanedWords := []string{}
			for _, word := range words {
				if strings.ToLower(word) == profanity {
					cleanedWords = append(cleanedWords, "****")
					continue
				}
				cleanedWords = append(cleanedWords, word)
			}
			params.Body = strings.Join(cleanedWords, " ")
		}
	}
	responseBody := returnCleaned{
		CleanedBody: params.Body,
	}
	responseData, error := json.Marshal(responseBody)
	if error != nil {
		responseBody := returnError{
			Error: "something went wrong",
		}
		responseData, _ := json.Marshal(responseBody)
		JsonResponse(responseData, writer, 500)
		return
	}
	JsonResponse(responseData, writer, 200)
}

func JsonResponse(responseData []byte, writer http.ResponseWriter, statusCode int) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	writer.Write(responseData)
}
