package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Chystik/runtime-metrics/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var nextIPValidatorHandler = func(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(nil)
	})
}

func makeIPValidatorRequest(t *testing.T, h http.Handler, ip string) int {
	req := httptest.NewRequest(http.MethodPost, "http://testing", nil)
	req.Header.Set("X-Real-IP", ip)

	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	res := rec.Result()
	res.Body.Close()

	return res.StatusCode
}

func Test_ipValidator_Validate(t *testing.T) {
	t.Parallel()

	type args struct {
		ip string
	}
	tests := []struct {
		name          string
		trustedSubnet string
		args          args
		wantStatus    int
	}{
		{
			name:          "valid and trusted ip",
			trustedSubnet: "127.0.0.0/8",
			args: args{
				ip: "127.0.0.1",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:          "empty ip",
			trustedSubnet: "127.0.0.0/8",
			args: args{
				ip: "",
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:          "wrong ip format",
			trustedSubnet: "127.0.0.0/8",
			args: args{
				ip: "wrong ip format",
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:          "valid and not trusted ip format",
			trustedSubnet: "127.0.0.0/8",
			args: args{
				ip: "192.168.0.1",
			},
			wantStatus: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &mocks.AppLogger{}
			v, err := NewIPValidator(tt.trustedSubnet, l)
			if err != nil {
				t.Error(err)
			}

			if tt.wantStatus != http.StatusOK {
				l.EXPECT().Info(mock.Anything).Return()
			}

			if got := makeIPValidatorRequest(
				t,
				v.Validate(nextIPValidatorHandler(t)),
				tt.args.ip,
			); got != tt.wantStatus {
				t.Errorf("ipValidator.Validate() = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func Test_NewIPValidator_WrongSubnet(t *testing.T) {
	t.Parallel()

	v, err := NewIPValidator("wrong subnet format", &mocks.AppLogger{})

	assert.Nil(t, v)
	assert.Error(t, err)
}
