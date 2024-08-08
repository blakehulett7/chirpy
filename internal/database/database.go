package database

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"sync"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Database struct {
	Path  string
	Mutex *sync.RWMutex
}

type DatabaseStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"email"`
}

func (databaseAddress *Database) EnsureDatabase() {
	db := DatabaseStructure{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
	}
	data, _ := json.Marshal(db)
	os.WriteFile(databaseAddress.Path, data, 0666)
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

func (databaseAddress *Database) CreateChirp(body string) Chirp {
	db := databaseAddress.LoadDatabase()
	id := len(db.Chirps) + 1
	chirp := Chirp{
		Id:   id,
		Body: body,
	}
	db.Chirps[id] = chirp
	databaseAddress.SaveDatabase(db)
	return chirp
}

func (databaseAddress *Database) GetChirps() []Chirp {
	db := databaseAddress.LoadDatabase()
	chirpArray := []Chirp{}
	for i := 1; i <= len(db.Chirps); i++ {
		chirpArray = append(chirpArray, db.Chirps[i])
	}
	return chirpArray
}

func (databaseAddress *Database) CreateUser(email string, password string) User {
	db := databaseAddress.LoadDatabase()
	id := len(db.Users) + 1
	passwordHash, error := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if error != nil {
		fmt.Println(error)
		return User{}
	}
	user := User{
		Id:       id,
		Email:    email,
		Password: string(passwordHash),
	}
	db.Users[id] = user
	databaseAddress.SaveDatabase(db)
	return user
}

func (databaseAddress *Database) GetUser(email string) (User, bool) {
	db := databaseAddress.LoadDatabase()
	for _, user := range db.Users {
		if email == user.Email {
			return user, true
		}
	}
	return User{}, false
}
