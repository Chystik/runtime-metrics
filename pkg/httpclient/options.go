package httpclient

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net"
	"net/http"
	"os"
	"time"
)

type Options func(*client) error

func Timeout(t time.Duration) Options {
	return func(c *client) error {
		c.timeout = t
		return nil
	}
}

func WithEncryption(publicKeyFilePath string) Options {
	return func(c *client) error {
		pemKey, err := os.ReadFile(publicKeyFilePath)
		if err != nil {
			return err
		}

		publicKeyBlock, _ := pem.Decode(pemKey)
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
	return func(c *client) error {

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
