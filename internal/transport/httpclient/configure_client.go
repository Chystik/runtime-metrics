package httpclient

import (
	"net/http"
	"time"

	"github.com/Chystik/runtime-metrics/config"
)

func NewHTTPClient(cfg *config.AgentConfig) *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}