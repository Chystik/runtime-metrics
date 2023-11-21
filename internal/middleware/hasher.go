package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
)

type (
	hasherResponseWriter struct {
		http.ResponseWriter
		buf    bytes.Buffer
		status int
	}

	hasher struct {
		key        []byte
		headerName string
	}
)

func NewHasher(key string, headerName string) *hasher {
	return &hasher{
		key:        []byte(key),
		headerName: headerName,
	}
}

func (hw *hasherResponseWriter) Write(b []byte) (int, error) {
	return hw.buf.Write(b)
}

func (hw *hasherResponseWriter) WriteHeader(statusCode int) {
	hw.status = statusCode
}

func (h *hasher) WithHasher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read req
		containHash := r.Header.Get(h.headerName)
		ow := &hasherResponseWriter{
			ResponseWriter: w,
			status:         200,
		}

		if containHash != "" && len(h.key) > 0 {
			// get hash from header
			requestedHash, _ := base64.StdEncoding.DecodeString(r.Header.Get(h.headerName))

			// calculate hash
			hm := hmac.New(sha256.New, h.key)

			// read body to calculate hash
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			hm.Write(body)
			calculatedHash := hm.Sum(nil)

			if !hmac.Equal(requestedHash, calculatedHash) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// set new body, wich will contain same data we read
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		next.ServeHTTP(ow, r)

		// calculate hash
		if ow.status < 300 && len(h.key) > 0 {
			hs := hmac.New(sha256.New, h.key)
			hs.Write(ow.buf.Bytes())
			sign := hs.Sum(nil)
			calculatedHash := base64.StdEncoding.EncodeToString(sign)
			w.Header().Set(h.headerName, calculatedHash)
		}

		w.WriteHeader(ow.status)
		_, err := w.Write(ow.buf.Bytes())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
