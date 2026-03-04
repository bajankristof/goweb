package config

import (
	"crypto/ecdsa"
	"fmt"
	"os"

	"github.com/bajankristof/goweb/jwt"
)

type EncryptionConfig struct {
	PublicKeyPath string           `toml:"public_key_path" env:"PUBLIC_KEY_PATH"`
	PublicKeyPEM  string           `toml:"public_key_pem" env:"PUBLIC_KEY_PEM"`
	PublicKey     *ecdsa.PublicKey `toml:"-"`

	PrivateKeyPath string            `toml:"private_key_path" env:"PRIVATE_KEY_PATH"`
	PrivateKeyPEM  string            `toml:"private_key_pem" env:"PRIVATE_KEY_PEM"`
	PrivateKey     *ecdsa.PrivateKey `toml:"-"`
}

func (c *EncryptionConfig) LoadKeys() error {
	err := c.loadPublicKey()
	if err != nil {
		return fmt.Errorf("load public key: %w", err)
	}

	err = c.loadPrivateKey()
	if err != nil {
		return fmt.Errorf("load private key: %w", err)
	}

	return nil
}

func (c *EncryptionConfig) loadPublicKey() error {
	var pem []byte
	if c.PublicKeyPEM != "" {
		pem = []byte(c.PublicKeyPEM)
	} else {
		var err error
		pem, err = os.ReadFile(c.PublicKeyPath)
		if err != nil {
			return fmt.Errorf("read %s: %w", c.PublicKeyPath, err)
		}
	}

	key, err := jwt.ParsePublicKeyFromPEM(pem)
	if err != nil {
		return fmt.Errorf("parse from PEM: %w", err)
	}

	c.PublicKey = key
	return nil
}

func (c *EncryptionConfig) loadPrivateKey() error {
	var pem []byte
	if c.PrivateKeyPEM != "" {
		pem = []byte(c.PrivateKeyPEM)
	} else {
		var err error
		pem, err = os.ReadFile(c.PrivateKeyPath)
		if err != nil {
			return fmt.Errorf("read %s: %w", c.PrivateKeyPath, err)
		}
	}

	key, err := jwt.ParsePrivateKeyFromPEM(pem)
	if err != nil {
		return fmt.Errorf("parse from PEM: %w", err)
	}

	c.PrivateKey = key
	return nil
}
