// Used to make API calls to the metrics server.
//
// Sends metrics to the metrics server with or without data compressing.
// Can send metrics in batch or one at a time.
package agentapiclient

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/Chystik/runtime-metrics/internal/service"
)

var (
	errBadStatusCode = "resp status code: %s"
)

type agentAPIClient struct {
	client  service.HTTPClient
	address string
	shaKey  string
}

// New creates new agent API client, wich sends http requests to the metrics server
func New(c service.HTTPClient, s *config.AgentConfig) *agentAPIClient {
	return &agentAPIClient{
		client:  c,
		address: s.Address,
		shaKey:  s.SHAkey,
	}
}

// ReportMetrics sends metrics one by one to the metrics server
func (ac *agentAPIClient) ReportMetrics(ctx context.Context, metrics map[string]interface{}) error {
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

		request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
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

// ReportMetricsJSON sends all metrics one by one to the metrics server using the JSON data format.
// It also compresses all data sent.
func (ac *agentAPIClient) ReportMetricsJSON(ctx context.Context, metrics map[string]models.Metric) error {
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

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &reqBody)
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

// ReportMetricsJSONBatch sends all metrics in batch to the metrics server using JSON data format.
// It also compresses all data sent.
func (ac *agentAPIClient) ReportMetricsJSONBatch(ctx context.Context, metrics map[string]models.Metric) error {
	var ms []models.Metric
	var buf, reqBody bytes.Buffer

	for _, m := range metrics {
		ms = append(ms, m)
	}

	url := fmt.Sprintf("http://%s/updates/", ac.address)
	err := json.NewEncoder(&buf).Encode(ms)
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &reqBody)
	if err != nil {
		return err
	}

	if ac.shaKey != "" {
		h := hmac.New(sha256.New, []byte(ac.shaKey))
		_, err = h.Write(reqBody.Bytes())
		if err != nil {
			return err
		}

		sign := h.Sum(nil)
		hVal := base64.StdEncoding.EncodeToString(sign)
		req.Header.Set("HashSHA256", hVal)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Encoding", "gzip")
	resp, err := ac.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(errBadStatusCode, resp.Status)
	}

	return nil
}
