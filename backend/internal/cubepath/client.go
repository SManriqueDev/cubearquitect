package cubepath

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) request(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "CubeArchitect/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	resBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cubepath api error: %d - %s", resp.StatusCode, string(resBody))
	}

	return resBody, nil
}

func (c *Client) Get(path string) (json.RawMessage, error) {
	return c.request(context.Background(), http.MethodGet, path, nil)
}

func (c *Client) Post(path string, body interface{}) (json.RawMessage, error) {
	return c.request(context.Background(), http.MethodPost, path, body)
}

func (c *Client) Put(path string, body interface{}) (json.RawMessage, error) {
	return c.request(context.Background(), http.MethodPut, path, body)
}

func (c *Client) Patch(path string, body interface{}) (json.RawMessage, error) {
	return c.request(context.Background(), http.MethodPatch, path, body)
}

func (c *Client) Delete(path string) (json.RawMessage, error) {
	return c.request(context.Background(), http.MethodDelete, path, nil)
}
