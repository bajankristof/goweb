package config

import (
	"crypto/ecdsa"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bajankristof/goweb/jwt"
)

type Duration time.Duration

func (d *Duration) UnmarshalText(text []byte) error {
	pd, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}

	*d = Duration(pd)
	return nil
}

type PrivateKey ecdsa.PrivateKey

func (k *PrivateKey) UnmarshalText(text []byte) error {
	pem := text
	f := string(text)
	if !strings.HasPrefix(f, "-----BEGIN") {
		var err error
		pem, err = os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read private key file %q: %w", f, err)
		}
	}

	pk, err := jwt.ParsePrivateKeyFromPEM(pem)
	if err != nil {
		return fmt.Errorf("parse private key: %w", err)
	}

	*k = PrivateKey(*pk)
	return nil
}

func (k *PrivateKey) Unwrap() *ecdsa.PrivateKey {
	return (*ecdsa.PrivateKey)(k)
}

type URL url.URL

func (u *URL) String() string {
	return (*url.URL)(u).String()
}

func (u *URL) UnmarshalText(text []byte) error {
	pu, err := url.Parse(string(text))
	if err != nil {
		return err
	}

	*u = URL(*pu)
	return nil
}

func (u *URL) Unwrap() *url.URL {
	return (*url.URL)(u)
}
