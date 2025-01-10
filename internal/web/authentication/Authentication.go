package authentication

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Authentication struct {
	redisClient *redis.Client
	nonceSize   int
	expiration  time.Duration
}

func (a *Authentication) CreateChallenge(ctx context.Context) (Challenge, error) {
	challenge, err := NewAuthChallenge(a.nonceSize)
	if err != nil {
		return challenge, fmt.Errorf("could not create auth challenge: %w", err)
	}

	err = a.redisClient.Set(ctx, challenge.Key(), challenge, a.expiration).Err()
	if err != nil {
		return challenge, fmt.Errorf("could not set redis auth challenge: %w", err)
	}

	return challenge, nil
}

func (a *Authentication) ValidateChallengeResponse(ctx context.Context, response ChallengeResponse) error {
	_, err := a.redisClient.Get(ctx, response.IssuedChallenge.Key()).Result()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("challenge not found")
	}
	if err != nil {
		return fmt.Errorf("could not get challenge: %w", err)
	}

	parsedKey, err := x509.ParsePKIXPublicKey(response.PublicKey)
	if err != nil {
		return fmt.Errorf("could not parse public key: %w", err)
	}

	hash := sha256.Sum256(response.SignedData)
	slice := hash[:]
	switch publicKey := parsedKey.(type) {
	case *rsa.PublicKey:
		// TODO VerifyPSS is modern and has mathematical proof
		err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, slice, response.Signature)
	case *ecdsa.PublicKey:
		res := ecdsa.VerifyASN1(publicKey, slice, response.Signature)
		if res == false {
			err = fmt.Errorf("ecdsa signature verification failed")
		}
	case ed25519.PublicKey:
		if len(publicKey) != ed25519.PublicKeySize {
			err = fmt.Errorf("invalid ed25519 public key size: %d", len(publicKey))
		}

		res := ed25519.Verify(publicKey, response.SignedData, response.Signature)
		if res == false {
			err = fmt.Errorf("ed25519 signature verification failed")
		}
	default:
		err = fmt.Errorf("unknown public key type: %T", publicKey)
	}

	if err != nil {
		_ = a.redisClient.Del(ctx, response.IssuedChallenge.Key()).Err()
	}

	return err
}

func NewAuthentication(redisClient *redis.Client, nonceSize int, expiration time.Duration) Authentication {
	return Authentication{
		redisClient: redisClient,
		nonceSize:   nonceSize,
		expiration:  expiration,
	}
}
