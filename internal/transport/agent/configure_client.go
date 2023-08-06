package httpclient

import (
	"net/http"

	"github.com/Chystik/runtime-metrics/config"
)

func NewHTTPClient(cfg config.HTTPServer) *http.Client {
	return &http.Client{}
}
