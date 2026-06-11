package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
	apiKey  string
}

type Option func(*Client)

func New(baseURL string, opts ...Option) *Client {
	c := &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 10 * time.Second}}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		if h != nil {
			c.http = h
		}
	}
}

func WithAPIKey(key string) Option {
	return func(c *Client) { c.apiKey = key }
}

func (c *Client) Index(ctx context.Context, index, id string, doc any) error {
	method := http.MethodPost
	path := "/" + index + "/_doc"
	if id != "" {
		method = http.MethodPut
		path += "/" + id
	}
	_, err := c.do(ctx, method, path, doc)
	return err
}

func (c *Client) Search(ctx context.Context, index string, query any, out any) error {
	b, err := c.do(ctx, http.MethodGet, "/"+index+"/_search", query)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}

func (c *Client) Get(ctx context.Context, index, id string, out any) error {
	b, err := c.do(ctx, http.MethodGet, "/"+index+"/_doc/"+id, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}

func (c *Client) Delete(ctx context.Context, index, id string) error {
	_, err := c.do(ctx, http.MethodDelete, "/"+index+"/_doc/"+id, nil)
	return err
}

func (c *Client) do(ctx context.Context, method, path string, body any) ([]byte, error) {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "ApiKey "+c.apiKey)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("es: %s: %s", resp.Status, string(b))
	}
	return b, nil
}
