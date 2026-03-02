package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Time = jwt.NumericDate
type Claims = jwt.Claims
type RegisteredClaims = jwt.RegisteredClaims
type VerifyOption = jwt.ParserOption

var WithAudience = jwt.WithAudience
var WithAllAudiences = jwt.WithAllAudiences
var WithIssuedAt = jwt.WithIssuedAt
var WithIssuer = jwt.WithIssuer
var WithSubject = jwt.WithSubject
var WithoutClaimsValidation = jwt.WithoutClaimsValidation

var ParsePublicKeyFromPEM = jwt.ParseECPublicKeyFromPEM
var ParsePrivateKeyFromPEM = jwt.ParseECPrivateKeyFromPEM

func Now() *Time {
	return jwt.NewNumericDate(time.Now())
}

func After(d time.Duration) *Time {
	return jwt.NewNumericDate(time.Now().Add(d))
}
