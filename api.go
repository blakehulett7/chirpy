package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"internal/database"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type apiState struct {
	serverHits int
	db         *database.Database
	jwtSecret  []byte
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
	bearerToken := request.Header.Get("Authorization")
	givenToken, _ := strings.CutPrefix(bearerToken, "Bearer ")
	claims := jwt.RegisteredClaims{}
	token, error := jwt.ParseWithClaims(givenToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return state.jwtSecret, nil
	})
	if error != nil {
		fmt.Println(error)
		JsonResponse([]byte("Unauthorized"), writer, 401)
		return
	}
	tokenClaims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		fmt.Println(ok)
		return
	}
	userId, _ := strconv.Atoi(tokenClaims.Subject)
	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	error = decoder.Decode(&params)
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
	responseBody := state.db.CreateChirp(params.Body, userId)
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

func (state *apiState) GetChirpy(writer http.ResponseWriter, request *http.Request) {
	chirpArray := state.db.GetChirps()
	responseData, _ := json.Marshal(chirpArray)
	JsonResponse(responseData, writer, 200)
}

func (state *apiState) GetaBitChirpy(writer http.ResponseWriter, request *http.Request) {
	chirpArray := state.db.GetChirps()
	indexString := request.PathValue("id")
	index, _ := strconv.Atoi(indexString)
	index--
	if index < 0 || index >= len(chirpArray) {
		JsonResponse([]byte("Not found"), writer, 404)
		return
	}
	responseBody := chirpArray[index]
	responseData, _ := json.Marshal(responseBody)
	JsonResponse(responseData, writer, 200)
}

func (state *apiState) CreateUser(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	userParameters := userParams{}
	decoder.Decode(&userParameters)
	_, userExists := state.db.GetUser(userParameters.Email)
	if userExists {
		fmt.Println("user already exists!")
		return
	}
	user := state.db.CreateUser(userParameters.Email, userParameters.Password)
	responseBody := responseUser{Id: user.Id, Email: user.Email}
	responseData, _ := json.Marshal(responseBody)
	JsonResponse(responseData, writer, 201)
}

func (state *apiState) Login(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	userParameters := LoginParams{}
	decoder.Decode(&userParameters)
	user, userExists := state.db.GetUser(userParameters.Email)
	if !userExists {
		JsonResponse([]byte("Unauthorized, user does not exist"), writer, 401)
		return
	}
	passwordsMatch := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userParameters.Password))
	if passwordsMatch != nil {
		JsonResponse([]byte("Unauthorized, password incorrect"), writer, 401)
		return
	}
	expires := time.Now().Add(time.Duration(60*60) * time.Second)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expires),
		Subject:   strconv.Itoa(user.Id),
	})
	signedToken, _ := token.SignedString(state.jwtSecret)
	state.db.UpdateUserToken(user.Email, signedToken)
	responseLogin := responseLogin{Id: user.Id, Email: user.Email, Token: signedToken, RefreshToken: user.RefreshToken.Token}
	responseData, _ := json.Marshal(responseLogin)
	JsonResponse(responseData, writer, 200)
}

func (state *apiState) UpdateUser(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	userParameters := LoginParams{}
	decoder.Decode(&userParameters)
	bearerToken := request.Header.Get("Authorization")
	givenToken, _ := strings.CutPrefix(bearerToken, "Bearer ")
	claims := jwt.RegisteredClaims{}
	token, error := jwt.ParseWithClaims(givenToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return state.jwtSecret, nil
	})
	if error != nil {
		fmt.Println(error)
		JsonResponse([]byte("Unauthorized"), writer, 401)
		return
	}
	tokenClaims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		fmt.Println(ok)
		return
	}
	userId, _ := strconv.Atoi(tokenClaims.Subject)
	updatedUser := state.db.UpdateUserCredentials(userId, userParameters.Email, userParameters.Password)
	responseBody := responseUser{
		Id:    updatedUser.Id,
		Email: updatedUser.Email,
	}
	responseData, _ := json.Marshal(responseBody)
	JsonResponse(responseData, writer, 200)
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

func (state *apiState) Refresh(writer http.ResponseWriter, request *http.Request) {
	bearerToken := request.Header.Get("Authorization")
	givenToken, _ := strings.CutPrefix(bearerToken, "Bearer ")
	valid, user := state.db.RefreshTokenIsValid(givenToken)
	if !valid {
		JsonResponse([]byte("Unauthorized, bad refresh token"), writer, 401)
		return
	}
	token := state.CreateToken(user)
	state.db.UpdateUserToken(user.Email, token)
	responseBody := responseRefresh{Token: token}
	responseData, _ := json.Marshal(responseBody)
	JsonResponse(responseData, writer, 200)
}

func (state *apiState) CreateToken(user database.User) string {
	expires := time.Now().Add(time.Duration(1) * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expires),
		Subject:   strconv.Itoa(user.Id),
	})
	signedToken, _ := token.SignedString(state.jwtSecret)
	return signedToken
}

func (state *apiState) Revoke(writer http.ResponseWriter, request *http.Request) {
	bearerToken := request.Header.Get("Authorization")
	givenToken, _ := strings.CutPrefix(bearerToken, "Bearer ")
	valid, user := state.db.RefreshTokenIsValid(givenToken)
	if !valid {
		JsonResponse([]byte("Unauthorized, bad refresh token"), writer, 401)
		return
	}
	state.db.RevokeRefreshToken(user)
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(204)
}
