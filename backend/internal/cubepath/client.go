package cubepath

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func (c *Client) request(method, path string, body interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, c.baseURL+path, &buf)
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
	defer resp.Body.Close()

	resBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cubepath api error: %d - %s", resp.StatusCode, string(resBody))
	}

	return resBody, nil
}

func (c *Client) Get(path string) (json.RawMessage, error) {
	return c.request(http.MethodGet, path, nil)
}

func (c *Client) Post(path string, body interface{}) (json.RawMessage, error) {
	return c.request(http.MethodPost, path, body)
}

func (c *Client) Put(path string, body interface{}) (json.RawMessage, error) {
	return c.request(http.MethodPut, path, body)
}

func (c *Client) Patch(path string, body interface{}) (json.RawMessage, error) {
	return c.request(http.MethodPatch, path, body)
}

func (c *Client) Delete(path string) (json.RawMessage, error) {
	return c.request(http.MethodDelete, path, nil)
}

func (c *Client) GetProjects() (ProjectResponse, error) {
	res, err := c.Get("/projects/")
	if err != nil {
		log.Printf("Error fetching projects: %v", err)
		return nil, err
	}
	var projects ProjectResponse
	err = json.Unmarshal(res, &projects)
	return projects, err
}

func (c *Client) CreateVPS(projectID int, req VPSCreateRequest) (*VPS, error) {
	res, err := c.Post(fmt.Sprintf("/vps/create/%d", projectID), req)
	if err != nil {
		log.Printf("Error creating VPS: %v", err)
		return nil, err
	}
	var vps VPS
	err = json.Unmarshal(res, &vps)
	log.Printf("VPS created: %+v", vps)
	return &vps, err
}

func (c *Client) GetPricing() (json.RawMessage, error) {
	res, err := c.Get("/pricing/")
	if err != nil {
		log.Printf("Error fetching pricing: %v", err)
		return nil, err
	}
	return res, err
}
