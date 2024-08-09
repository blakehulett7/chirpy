package database

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"sync"
	"time"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id           int          `json:"id"`
	Email        string       `json:"email"`
	Password     string       `json:"password"`
	Token        string       `json:"token"`
	RefreshToken RefreshToken `json:"refresh_token"`
}

type RefreshToken struct {
	Token   string
	Expires time.Time
}

type Database struct {
	Path  string
	Mutex *sync.RWMutex
}

type DatabaseStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"email"`
}

type UserParams struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
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
	databaseAddress.Mutex.Lock()
	defer databaseAddress.Mutex.Unlock()
	dbData, _ := os.ReadFile(databaseAddress.Path)
	db := DatabaseStructure{}
	json.Unmarshal(dbData, &db)
	return db
}

func (databaseAddress *Database) SaveDatabase(db DatabaseStructure) {
	databaseAddress.Mutex.Lock()
	defer databaseAddress.Mutex.Unlock()
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
	hexData := make([]byte, 32)
	rand.Read(hexData)
	hexString := hex.EncodeToString(hexData)
	expires := time.Now().AddDate(0, 0, 60)
	refreshToken := RefreshToken{
		Token:   hexString,
		Expires: expires,
	}
	user := User{
		Id:           id,
		Email:        email,
		Password:     string(passwordHash),
		Token:        "",
		RefreshToken: refreshToken,
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

func (databaseAddress *Database) UpdateUserToken(email string, token string) {
	db := databaseAddress.LoadDatabase()
	for _, user := range db.Users {
		if email == user.Email {
			user.Token = token
			databaseAddress.SaveDatabase(db)
			return
		}
	}
	fmt.Println("Email Not Found")
}

func (databaseAddress *Database) UpdateUserCredentials(id int, email string, password string) User {
	db := databaseAddress.LoadDatabase()
	oldUserInfo := db.Users[id]
	passwordHash, error := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if error != nil {
		fmt.Println(error)
		return User{}
	}
	updatedUserInfo := User{
		Id:           id,
		Email:        email,
		Password:     string(passwordHash),
		Token:        oldUserInfo.Token,
		RefreshToken: oldUserInfo.RefreshToken,
	}
	db.Users[id] = updatedUserInfo
	databaseAddress.SaveDatabase(db)
	return updatedUserInfo
}

func (databaseAddress *Database) RefreshTokenIsValid(RefreshToken string) (bool, User) {
	db := databaseAddress.LoadDatabase()
	for _, user := range db.Users {
		fmt.Println(RefreshToken, "<div>", user.RefreshToken.Token)
		fmt.Println("Now: ", time.Now(), "Expires: ", user.RefreshToken.Expires)
		if RefreshToken == user.RefreshToken.Token {
			fmt.Println("Found refresh token")
			if time.Now().Before(user.RefreshToken.Expires) {
				return true, user
			}
		}
	}
	return false, User{}
}
