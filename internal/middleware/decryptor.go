package middleware

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net/http"
	"os"
)

type decryptor struct {
	privateKey *rsa.PrivateKey
}

func NewDecryptor(privatePemFile string) (*decryptor, error) {
	privateKeyPEM, err := os.ReadFile(privatePemFile)
	if err != nil {
		return nil, err
	}

	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	d := &decryptor{privateKey: privateKey}

	return d, nil
}

func (d *decryptor) WithDecryptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			next.ServeHTTP(w, r)
		}

		decryptedBody, err := rsa.DecryptPKCS1v15(rand.Reader, d.privateKey, body)
		if err != nil {
			next.ServeHTTP(w, r)
		}

		r.Body = io.NopCloser(bytes.NewBuffer(decryptedBody))
		r.ContentLength = int64(len(decryptedBody))

		next.ServeHTTP(w, r)
	})
}
