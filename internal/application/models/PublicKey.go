package models

import (
	"crypto/x509"
	"fmt"
)

type PublicKey struct {
	DerData []byte `json:"der_data"`
}

func NewPublicKey(derData []byte) (PublicKey, error) {
	_, err := x509.ParsePKIXPublicKey(derData)
	if err != nil {
		return PublicKey{}, fmt.Errorf("could not create public key: %w", err)
	}

	return PublicKey{
		DerData: derData,
	}, nil
}
