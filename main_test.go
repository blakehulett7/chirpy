package main

import (
	"fmt"
	"internal/database"
	"os"
	"testing"
)

func TestCreateChirp(test *testing.T) {
	fmt.Println("Testing CreateChirp...")
	db := &database.Database{Path: "test_database.json"}
	db.EnsureDatabase()
	db.CreateChirp("Christ Is King!")
	db.CreateChirp("Jesus Is Lord!")
	haveChirpArray := db.GetChirps()
	wantChirpArray := []database.Chirp{
		{
			Id:   1,
			Body: "Christ Is King!",
		},
		{
			Id:   2,
			Body: "Jesus Is Lord!",
		},
	}
	os.Remove(db.Path)
	for index := range haveChirpArray {
		if haveChirpArray[index] != wantChirpArray[index] {
			fmt.Printf("have: %v", haveChirpArray)
			fmt.Printf("want: %v", wantChirpArray)
			test.Fatalf("CreateChirp not working as intended.")
		}
	}
}
