package models

import (
	"fmt"
)

// TODO probably needs bson flags too
type User struct {
	Name      string    `json:"name"`
	PublicKey PublicKey `json:"public_key"`
}

func NewUser(name string, derKey []byte) (User, error) {
	pk, err := NewPublicKey(derKey)
	if err != nil {
		return User{}, fmt.Errorf("could not create user: %w", err)
	}

	return User{
		Name:      name,
		PublicKey: pk,
	}, nil
}
