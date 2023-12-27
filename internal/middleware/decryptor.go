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
	buf        bytes.Buffer
}

func NewDecryptor(privatePemFilePath string) (*decryptor, error) {
	privateKeyPEM, err := os.ReadFile(privatePemFilePath)
	if err != nil {
		return nil, err
	}

	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	d := &decryptor{privateKey: privateKey}
	d.buf.Grow(512)

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

		_, err = d.buf.Write(decryptedBody)
		if err != nil {
			next.ServeHTTP(w, r)
		}
		defer d.buf.Reset()

		r.Body = io.NopCloser(&d.buf)
		r.ContentLength = int64(len(decryptedBody))

		next.ServeHTTP(w, r)
	})
}
