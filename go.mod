module github.com/blakehulett/chirpy

go 1.22.3

require internal/database v1.0.0

require golang.org/x/crypto v0.26.0 // indirect

replace internal/database => ./internal/database/
