package httpclient

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"time"
)

type Options func(*Client) error

func Timeout(t time.Duration) Options {
	return func(c *Client) error {
		c.timeout = t
		return nil
	}
}

func WithEncription(publicKeyPEM []byte) Options {
	return func(c *Client) error {
		publicKeyBlock, _ := pem.Decode(publicKeyPEM)
		publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
		if err != nil {
			return err
		}
		c.doMethod = &doWithEncryption{publicKey: publicKey.(*rsa.PublicKey)}
		return nil
	}
}
