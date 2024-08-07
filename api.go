package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"internal/database"
	"net/http"
	"strings"
)

type apiState struct {
	serverHits int
	db         *database.Database
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

func (state *apiState) PostChirp(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	error := decoder.Decode(&params)
	if error != nil {
		responseBody := returnError{
			Error: errors.New("something went wrong"),
		}
		responseData, _ := json.Marshal(responseBody)
		JsonResponse(responseData, writer, 500)
		return
	}
	params.Body, error = ChirpValidator(params)
	if error != nil {
		responseBody := returnError{
			Error: error,
		}
		responseData, _ := json.Marshal(responseBody)
		JsonResponse(responseData, writer, 500)
		return
	}
	responseBody := state.db.CreateChirp(params.Body)
	responseData, error := json.Marshal(responseBody)
	if error != nil {
		responseBody := returnError{
			Error: errors.New("something went wrong"),
		}
		responseData, _ := json.Marshal(responseBody)
		JsonResponse(responseData, writer, 400)
		return
	}
	JsonResponse(responseData, writer, 201)
}

func JsonResponse(responseData []byte, writer http.ResponseWriter, statusCode int) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	writer.Write(responseData)
}

func ChirpValidator(params parameters) (body string, error error) {
	if len(params.Body) > 140 {
		return "", errors.New("chirp is too long")
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
	return params.Body, nil
}
