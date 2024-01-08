package httpclient

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net"
	"net/http"
	"time"
)

type Options func(*Client) error

func Timeout(t time.Duration) Options {
	return func(c *Client) error {
		c.timeout = t
		return nil
	}
}

func WithEncryption(publicKeyPEM []byte) Options {
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

// ExtractOutboundIP extracts preferred outbound IP when http.Client dialing remote server.
// It sets the header entry associated with key headerName to the string representation of IP value.
func ExtractOutboundIP(headerName string) Options {
	return func(c *Client) error {

		// ip extraction function
		dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
			d := &net.Dialer{}

			conn, err := d.DialContext(ctx, network, addr)
			if err == nil {
				la := conn.LocalAddr().(*net.TCPAddr)
				c.req.Header.Set(headerName, la.IP.String())
			}
			return conn, err
		}

		httpClietn.Transport = &http.Transport{
			DialContext: dialContext,
		}
		return nil
	}
}
