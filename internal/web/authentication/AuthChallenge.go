package authentication

import (
	"crypto/rand"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type Challenge struct {
	Nonce []byte `bson:"nonce"`
}

func (ac Challenge) MarshalBinary() ([]byte, error) {
	return bson.Marshal(ac)
}

func (ac Challenge) Key() string {
	return fmt.Sprintf("auth_challenge:%x", ac.Nonce)
}

func NewAuthChallenge(size int) (Challenge, error) {
	nonce := make([]byte, size)
	_, err := rand.Read(nonce)
	if err != nil {
		return Challenge{}, fmt.Errorf("could not generate nonce: %w", err)
	}

	return Challenge{
		Nonce: nonce,
	}, nil
}

type ChallengeResponse struct {
	IssuedChallenge Challenge `bson:"issued_challenge"`
	SignedData      []byte    `bson:"signed_data"`
	Signature       []byte    `bson:"signature"`
	PublicKey       []byte    `bson:"public_key_der"`
}
