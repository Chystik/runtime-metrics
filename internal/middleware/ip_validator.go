package middleware

import (
	"fmt"
	"net"
	"net/http"

	"github.com/Chystik/runtime-metrics/internal/service"
)

type ipValidator struct {
	trustedSubnet *net.IPNet
	logger        service.AppLogger
}

func NewIPValidator(trustedSubnet string, logger service.AppLogger) (*ipValidator, error) {
	v := new(ipValidator)
	var err error

	_, v.trustedSubnet, err = net.ParseCIDR(trustedSubnet)
	if err != nil {
		return nil, err
	}

	v.logger = logger

	return v, nil
}

func (v *ipValidator) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("X-Real-IP")
		if h == "" {
			v.logger.Info("can't find X-Real-IP header")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ip := net.ParseIP(h)
		if ip == nil {
			v.logger.Info(fmt.Sprintf("can't parse IP from X-Real-IP header: %s", h))
			w.WriteHeader(http.StatusForbidden)
			return
		}

		trusted := v.trustedSubnet.Contains(ip)

		if !trusted {
			v.logger.Info(fmt.Sprintf("ip address %s not valid", ip.String()))
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
