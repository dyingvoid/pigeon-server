package authentication

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"github.com/dyingvoid/pigeon-server/internal/web/requests"
	"time"
)

type Authentication struct {
	challenges map[string]AuthChallenge
	nonceSize  int
	expiration time.Duration
}

// AddChallenge adds new challenge to collection and returns it
func (a *Authentication) AddChallenge() (AuthChallenge, error) {
	nonce := make([]byte, a.nonceSize)
	_, err := rand.Read(nonce)
	if err != nil {
		return AuthChallenge{},
			fmt.Errorf("could not generate nonce: %w", err)
	}

	session := AuthChallenge{
		Expiration: time.Now().Add(a.expiration),
	}
	a.challenges[string(nonce)] = session

	return session, nil
}

func (a *Authentication) GetChallenge(nonce []byte) (AuthChallenge, error) {
	challenge, ok := a.challenges[string(nonce)]
	if !ok {
		return AuthChallenge{}, fmt.Errorf("challenge not found")
	}

	return challenge, nil
}

// ValidateChallenge validates request's signature against corresponding json request's body
func (a *Authentication) ValidateChallenge(signature requests.RequestSignature, requestBody []byte) error {
	challenge, err := a.GetChallenge(signature.IssuedNonce)
	if err != nil {
		return a.authError(err)
	}

	if time.Now().After(challenge.Expiration) {
		return a.authError(fmt.Errorf("challenge expired"))
	}

	rsaPublicKey, err := a.parseRsaPublicKey(signature.PublicKey)
	if err != nil {
		return a.authError(err)
	}

	hash := sha256.Sum256(requestBody)
	err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hash[:], signature.Signature)
	if err != nil {
		return fmt.Errorf("could not verify signature: %w", err)
	}

	delete(a.challenges, string(signature.IssuedNonce))

	return nil
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

func (a *Authentication) authError(err error) error {
	return fmt.Errorf("authentication error: %w", err)
}

func NewSessionContainer() Authentication {
	return Authentication{challenges: make(map[string]AuthChallenge)}
}
