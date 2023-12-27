package middleware

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Chystik/runtime-metrics/internal/models"
	_ "github.com/Chystik/runtime-metrics/pkg/cert"
)

var (
	certDir        string = "./cert"
	publicKeyFile  string = "public_key.pem"
	privateKeyFile string = "private_key.pem"
)

var nextDecryptorHandler = func(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m models.Metric
		var body bytes.Buffer

		err := json.NewEncoder(&body).Encode(m)
		if err != nil {
			t.Error(err)
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		w.Write(nil)
	})
}

func makeDecryptorRequest(t *testing.T, h http.Handler, b []byte, k *rsa.PublicKey) int {
	var body bytes.Buffer

	encryptedBody, err := rsa.EncryptPKCS1v15(rand.Reader, k, b)
	if err != nil {
		t.Error(err)
	}

	_, err = body.Write(encryptedBody)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, "http://testing", &body)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)
	res := rec.Result()

	defer res.Body.Close()
	return res.StatusCode
}

func Test_decryptor_WithDecryptor(t *testing.T) {
	keyFile, err := os.ReadFile(fmt.Sprintf("%s/%s", certDir, publicKeyFile))
	if err != nil {
		t.Error(err)
	}

	publicKeyBlock, _ := pem.Decode(keyFile)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		t.Error(err)
	}

	d, err := NewDecryptor(fmt.Sprintf("%s/%s", certDir, privateKeyFile))
	if err != nil {
		t.Error(err)
	}

	type args struct {
		next     http.Handler
		testData []models.Metric
	}
	tests := []struct {
		name       string
		d          *decryptor
		args       args
		wantStatus int
	}{
		{
			name: "decryptor returns 200",
			d:    d,
			args: args{
				next:     nextDecryptorHandler(t),
				testData: generateMetrics(10),
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			for _, d := range tt.args.testData {
				var buf bytes.Buffer
				if err := json.NewEncoder(&buf).Encode(d); err != nil {
					t.Errorf(err.Error())
				}

				if got := makeDecryptorRequest(t, tt.d.WithDecryptor(tt.args.next), buf.Bytes(), publicKey.(*rsa.PublicKey)); got != tt.wantStatus {
					t.Errorf("decryptor.WithDecryptor() = %v, want %v", got, tt.wantStatus)
				}
			}
		})
	}

	os.RemoveAll(certDir)
}
