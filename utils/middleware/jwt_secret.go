package middleware

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
)

type JwtSecret interface {
	GetSign() any
	GetVerify() any
}

type EdDSASecret struct {
	PrivKey ed25519.PrivateKey
	PubKey  ed25519.PublicKey
}

func NewEdDSASecret(privateKeyBase64 string) (*EdDSASecret, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return nil, errors.New("invalid private key encoding")
	}

	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid private key size")
	}

	privKey := ed25519.PrivateKey(privateKeyBytes)
	pubKey := privKey.Public().(ed25519.PublicKey)

	return &EdDSASecret{
		PrivKey: privKey,
		PubKey:  pubKey,
	}, nil
}

func (s *EdDSASecret) GetSign() any {
	return s.PrivKey
}

func (s *EdDSASecret) GetVerify() any {
	return s.PubKey
}
