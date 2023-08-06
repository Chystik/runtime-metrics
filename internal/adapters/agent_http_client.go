package adapters

import (
	"net/http"

	"github.com/Chystik/runtime-metrics/config"
	agenthttpclient "github.com/Chystik/runtime-metrics/internal/adapters/agent_http_client"
)

type AgentClient interface {
	agenthttpclient.AgentHTTPClient
}

type agentClient struct {
	agenthttpclient.AgentHTTPClient
}

func NewAgentClient(hc *http.Client, cfg config.HTTPServer) AgentClient {
	return &agentClient{
		agenthttpclient.New(hc, cfg),
	}
}
