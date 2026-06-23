package session

import "errors"

var (
	ErrNotFound = errors.New("session: not found")
	ErrRevoked  = errors.New("session: revoked")
	ErrExpired  = errors.New("session: expired")
)
