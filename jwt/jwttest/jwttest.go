package jwttest

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/bajankristof/goweb/jwt"
)

// NewSigner returns a *jwt.Signer with a freshly generated P-256 ECDSA key pair.
func NewSigner(t testing.TB) *jwt.Signer {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("jwttest: generate ECDSA key: %v", err)
	}

	signer, err := jwt.NewSigner(&key.PublicKey, key)
	if err != nil {
		t.Fatalf("jwttest: create signer: %v", err)
	}

	return signer
}
