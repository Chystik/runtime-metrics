package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Chystik/runtime-metrics/internal/models"
)

const (
	shaHeaderName = "HashSHA256"
	shaKey        = "secret key"
)

var nextHasherHandler = func(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hashHeader := r.Header.Get(shaHeaderName)
		containHash := hashHeader != "" && len(hashHeader) > 0

		if !containHash {
			t.Error("no hash in the header")
		}

		var m models.Metric

		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			t.Error(err)
		}

		var body bytes.Buffer

		err = json.NewEncoder(&body).Encode(m)
		if err != nil {
			t.Error(err)
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		w.Write(nil)
	})
}

func makeHasherRequest(t *testing.T, h http.Handler, b []byte) int {
	var body bytes.Buffer
	_, _ = body.Write(b)

	hash := hmac.New(sha256.New, []byte(shaKey))
	_, err := hash.Write(b)
	if err != nil {
		t.Error(err)
	}

	sign := hash.Sum(nil)
	hVal := base64.StdEncoding.EncodeToString(sign)

	req := httptest.NewRequest(http.MethodPost, "http://testing", &body)
	req.Header.Set("HashSHA256", hVal)

	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	res := rec.Result()

	defer res.Body.Close()

	return res.StatusCode
}

func Test_hasher_WithHasher(t *testing.T) {
	type args struct {
		handlerToTest http.Handler
		testData      []models.Metric
	}
	tests := []struct {
		name       string
		h          *hasher
		args       args
		wantStatus int
	}{
		{
			name: "hasher returns 200",
			h:    NewHasher(shaKey, shaHeaderName),
			args: args{
				handlerToTest: nextHasherHandler(t),
				testData:      generateMetrics(10),
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Parallel()

		t.Run(tt.name, func(t *testing.T) {
			for _, d := range tt.args.testData {
				var buf bytes.Buffer
				if err := json.NewEncoder(&buf).Encode(d); err != nil {
					t.Errorf(err.Error())
				}

				if got := makeHasherRequest(t, tt.h.WithHasher(tt.args.handlerToTest), buf.Bytes()); got != tt.wantStatus {
					t.Errorf("hasher.WithHasher() = %v, want %v", got, tt.wantStatus)
				}
			}
		})
	}
}
