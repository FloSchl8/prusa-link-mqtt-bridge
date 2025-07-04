package prusalink

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Client struct {
	Host       string
	APIKey     string
	HttpClient *http.Client
}

func NewClient(host, apiKey string) *Client {
	return &Client{
		Host:   host,
		APIKey: apiKey,
		HttpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Get(path string) ([]byte, error) {
	slog.Debug("Attempting to GET from PrusaLink", "host", c.Host, "path", path)
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/api/v1/%s", c.Host, path), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", c.APIKey)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	slog.Debug("Received response from PrusaLink API",
		"status_code", resp.StatusCode,
		"body", string(body),
	)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return body, nil
}
