package agenthttpclient

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/models"
)

type agentHTTPClient struct {
	client  *http.Client
	address string
}

func New(c *http.Client, s *config.AgentConfig) *agentHTTPClient {
	return &agentHTTPClient{
		client:  c,
		address: s.Address,
	}
}

func (ac *agentHTTPClient) ReportMetrics(metrics map[string]interface{}) error {
	for name, value := range metrics {
		var mType string

		switch t := value.(type) {
		case float64:
			mType = "gauge"
		case int64:
			mType = "counter"
		default:
			return fmt.Errorf("unknown metric type: %#v", t)
		}

		url := fmt.Sprintf("http://%s/update/%s/%s/%v", ac.address, mType, name, value)

		request, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, url, nil)
		if err != nil {
			return err
		}
		request.Header.Set("Content-Type", "text/plain")
		response, err := ac.client.Do(request)
		if err != nil {
			return err
		}
		response.Body.Close()
	}
	return nil
}

func (ac *agentHTTPClient) ReportMetricsJSON(metrics map[string]models.Metric) error {
	for _, metric := range metrics {
		var buf, reqBody bytes.Buffer

		url := fmt.Sprintf("http://%s/update/", ac.address)
		err := json.NewEncoder(&buf).Encode(metric)
		if err != nil {
			return err
		}

		gz := gzip.NewWriter(&reqBody)
		_, err = gz.Write(buf.Bytes())
		if err != nil {
			return err
		}

		err = gz.Close()
		if err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, url, &reqBody)
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Content-Encoding", "gzip")
		resp, err := ac.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}

		var p bytes.Buffer
		_, err = io.Copy(&p, reader)
		if err != nil {
			return err
		}
		defer reader.Close()

		// does something with decompressed response
		_ = p.String
	}
	return nil
}
