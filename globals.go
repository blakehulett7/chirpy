package main

var metricsHTML string = `<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`

type parameters struct {
	Body string `json:"body"`
}

type returnError struct {
	Error error `json:"error"`
}

type userParams struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
}

type responseUser struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}
