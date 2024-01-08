package httpclient

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout = 20 * time.Second
)

var httpClietn *http.Client

type doMethod interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	timeout  time.Duration
	doMethod doMethod
	req      *http.Request
}

func NewClient(opts ...Options) (*Client, error) {
	httpClietn = &http.Client{
		Timeout: defaultTimeout,
	}

	// by default using Do method without encryption
	client := &Client{doMethod: &doWithoutEncryption{}}

	for _, opt := range opts {
		err := opt(client)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	c.req = req
	return c.doMethod.Do(c.req)
}

type doWithEncryption struct {
	publicKey *rsa.PublicKey
}

func (c *doWithEncryption) Do(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// encript data with public key
	encryptedBody, err := rsa.EncryptPKCS1v15(rand.Reader, c.publicKey, body)
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewBuffer(encryptedBody))
	req.ContentLength = int64(len(encryptedBody))

	return httpClietn.Do(req)
}

type doWithoutEncryption struct{}

func (c *doWithoutEncryption) Do(req *http.Request) (*http.Response, error) {
	return httpClietn.Do(req)
}
