package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type Database struct {
	Path  string
	Mutex *sync.RWMutex
}

type DatabaseStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func (databaseAddress *Database) EnsureDatabase() {
	_, error := os.ReadFile(databaseAddress.Path)
	if errors.Is(error, os.ErrNotExist) {
		db := DatabaseStructure{Chirps: map[int]Chirp{}}
		data, _ := json.Marshal(db)
		os.WriteFile(databaseAddress.Path, data, 0666)
	}
}

func (databaseAddress *Database) LoadDatabase() DatabaseStructure {
	dbData, _ := os.ReadFile(databaseAddress.Path)
	db := DatabaseStructure{}
	json.Unmarshal(dbData, &db)
	return db
}

func (databaseAddress *Database) SaveDatabase(db DatabaseStructure) {
	data, _ := json.Marshal(db)
	os.WriteFile(databaseAddress.Path, data, 0666)
}

func (databaseAddress *Database) CreateChirp(body string) {
	db := databaseAddress.LoadDatabase()
	id := len(db.Chirps) + 1
	fmt.Println(id)
	db.Chirps[id] = Chirp{
		Id:   id,
		Body: body,
	}
	databaseAddress.SaveDatabase(db)
}

func (databaseAddress *Database) GetChirps() []Chirp {
	db := databaseAddress.LoadDatabase()
	chirpArray := []Chirp{}
	for i := 1; i <= len(db.Chirps); i++ {
		chirpArray = append(chirpArray, db.Chirps[i])
	}
	return chirpArray
}
