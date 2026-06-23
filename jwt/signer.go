package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Signer struct {
	method jwt.SigningMethod
	key    *ecdsa.PrivateKey
}

type SignerOption func(*Signer)

func NewSigner(key *ecdsa.PrivateKey, opts ...SignerOption) (*Signer, error) {
	var method jwt.SigningMethod
	switch key.Curve {
	case elliptic.P256():
		method = jwt.SigningMethodES256
	case elliptic.P384():
		method = jwt.SigningMethodES384
	case elliptic.P521():
		method = jwt.SigningMethodES512
	default:
		return nil, fmt.Errorf(
			"unsupported elliptic curve %s",
			key.Curve.Params().Name,
		)
	}

	s := &Signer{
		method: method,
		key:    key,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

func (s *Signer) Sign(claims Claims) (string, error) {
	token := jwt.NewWithClaims(s.method, claims)

	return token.SignedString(s.key)
}

func (s *Signer) Verify(token string, claims Claims, opts ...VerifyOption) error {
	_, err := jwt.ParseWithClaims(token, claims, s.keyfunc, opts...)

	return err
}

func (s *Signer) keyfunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
		return nil, jwt.ErrTokenUnverifiable
	}

	return s.key.Public(), nil
}
