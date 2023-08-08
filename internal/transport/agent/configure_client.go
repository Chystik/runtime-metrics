package httpclient

import (
	"net/http"

	"github.com/Chystik/runtime-metrics/config"
)

func NewHTTPClient(cfg *config.AgentConfig) *http.Client {
	return &http.Client{}
}
