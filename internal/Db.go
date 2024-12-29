package internal

import (
	crand "crypto/rand"
	"math/rand/v2"
	"strconv"
)

// TODO name is not index, it could be optional thing to I believe
type Db struct {
	Users map[string]User
}

func (db *Db) AddUser(user User) {
	db.Users[user.Name] = user
}

func (db *Db) GetAllUsers() []User {
	users := make([]User, 0, len(db.Users))

	for _, user := range db.Users {
		users = append(users, user)
	}

	return users
}

func NewUser(keysize int) User {
	num := rand.Int()
	name := strconv.Itoa(num)
	bytes := make([]byte, keysize)
	_, err := crand.Read(bytes)
	if err != nil {
		panic(err)
	}

	return User{
		Name: name,
		PublicKey: PublicKey{
			Data: bytes,
		},
	}
}

type SignedUser struct {
	User      User `json:"user"`
	Signature `json:"signature"`
}

type User struct {
	Name      string    `json:"name"`
	PublicKey PublicKey `json:"public_key"`
}

type PublicKey struct {
	Data []byte `json:"der_data"`
}

type Signature struct {
	Data []byte `json:"data"`
}
