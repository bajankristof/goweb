package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Signer struct {
	publicKey  *ecdsa.PublicKey
	privateKey *ecdsa.PrivateKey
	method     jwt.SigningMethod
}

// NewSigner creates a new Signer instance from the provided ECDSA key pair.
func NewSigner(publicKey *ecdsa.PublicKey, privateKey *ecdsa.PrivateKey) (*Signer, error) {
	var method jwt.SigningMethod
	switch publicKey.Curve {
	case elliptic.P256():
		method = jwt.SigningMethodES256
	case elliptic.P384():
		method = jwt.SigningMethodES384
	case elliptic.P521():
		method = jwt.SigningMethodES512
	default:
		return nil, fmt.Errorf(
			"unsupported elliptic curve %s",
			publicKey.Curve.Params().Name,
		)
	}

	return &Signer{
		publicKey:  publicKey,
		privateKey: privateKey,
		method:     method,
	}, nil
}

// Sign returns a signed JWT token string based on the provided claims.
func (s *Signer) Sign(claims Claims) (string, error) {
	token := jwt.NewWithClaims(s.method, claims)

	return token.SignedString(s.privateKey)
}

// Verify checks the validity of the provided JWT token string against the provided claims.
func (s *Signer) Verify(token string, claims Claims, opts ...VerifyOption) error {
	_, err := jwt.ParseWithClaims(token, claims, s.keyfunc, opts...)

	return err
}

// keyfunc is a helper function that returns the public key for verifying the token's signature.
func (s *Signer) keyfunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
		return nil, jwt.ErrTokenUnverifiable
	}

	return s.publicKey, nil
}
