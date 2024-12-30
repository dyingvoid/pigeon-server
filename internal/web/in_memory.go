package web

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"time"
)

type AuthChallenge struct {
	Expiration time.Time
}

// TODO safe len/size
const size = 256

// TODO cleanup unused sessions
type Authentication struct {
	challenges map[string]AuthChallenge
}

type ClientChallengeResponse struct {
	Nonce     []byte
	PublicKey []byte
	Signature []byte
}

func (a *Authentication) AddChallenge() (AuthChallenge, error) {
	nonceBytes := make([]byte, size)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return AuthChallenge{},
			fmt.Errorf("could not generate nonce: %w", err)
	}

	hash := sha256.Sum256(nonceBytes)
	nonce := hash[:]
	session := AuthChallenge{
		Expiration: time.Now().Add(time.Second * 10),
	}
	a.challenges[string(nonce)] = session

	return session, nil
}

func (a *Authentication) ValidateChallenge(r ClientChallengeResponse) (bool, error) {
	key := string(r.Nonce)
	challenge, ok := a.challenges[key]
	if !ok {
		return false, fmt.Errorf("challenge not found")
	}

	if time.Now().After(challenge.Expiration) {
		return false, fmt.Errorf("challenge expired")
	}

	rsaPublicKey, err := a.parseRsaPublicKey(r.PublicKey)
	if err != nil {
		return false, fmt.Errorf("could not get key: %w", err)
	}

	err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, r.Nonce, r.Signature)
	if err != nil {
		return false, fmt.Errorf("could not verify signature: %w", err)
	}

	delete(a.challenges, key)

	return true, nil
}

func (a *Authentication) parseRsaPublicKey(rsaBytes []byte) (*rsa.PublicKey, error) {
	key, err := x509.ParsePKIXPublicKey(rsaBytes)
	if err != nil {
		return nil,
			fmt.Errorf("could not parse RSA public key: %w", err)
	}

	// TODO get information about casting
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil,
			fmt.Errorf("key is not an RSA public key")
	}

	return rsaKey, nil
}

func NewSessionContainer() Authentication {
	return Authentication{challenges: make(map[string]AuthChallenge)}
}
