package http_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeForm = "application/x-www-form-urlencoded"
)

type Client struct {
	client      *http.Client
	baseURL     string
	headers     map[string]string
	contentType string
}

// NewClient 创建 HTTP Client
func NewClient(baseURL string) *Client {

	baseURL = strings.TrimRight(baseURL, "/")

	return &Client{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:     baseURL,
		headers:     make(map[string]string),
		contentType: ContentTypeJSON,
	}
}

// SetTimeout 设置请求超时时间
func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

// SetTransport 设置 Transport
func (c *Client) SetTransport(t http.RoundTripper) {
	c.client.Transport = t
}

// SetHeader 设置 Header
func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

// RemoveHeader 删除 Header
func (c *Client) RemoveHeader(key string) {
	delete(c.headers, key)
}

// ClearHeaders 清空 Header
func (c *Client) ClearHeaders() {
	c.headers = make(map[string]string)
}

// SetContentType 设置 Content-Type
func (c *Client) SetContentType(contentType string) {
	c.contentType = contentType
}

// ---------------- GET ----------------

func (c *Client) Get(path string, params map[string]string, result any) error {

	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return err
	}

	query := u.Query()

	for k, v := range params {
		query.Set(k, v)
	}

	u.RawQuery = query.Encode()

	return c.request(
		http.MethodGet,
		u.String(),
		nil,
		result,
	)
}

// ---------------- POST JSON ----------------

func (c *Client) Post(path string, body any, result any) error {

	var reader io.Reader

	switch v := body.(type) {

	case nil:

	case []byte:
		reader = bytes.NewReader(v)

	case string:
		reader = strings.NewReader(v)

	default:

		data, err := json.Marshal(v)
		if err != nil {
			return err
		}

		reader = bytes.NewReader(data)
	}

	return c.request(
		http.MethodPost,
		c.baseURL+path,
		reader,
		result,
	)
}

// ---------------- PUT ----------------

func (c *Client) Put(path string, body any, result any) error {

	var reader io.Reader

	switch v := body.(type) {

	case nil:

	case []byte:
		reader = bytes.NewReader(v)

	case string:
		reader = strings.NewReader(v)

	default:

		data, err := json.Marshal(v)
		if err != nil {
			return err
		}

		reader = bytes.NewReader(data)
	}

	return c.request(
		http.MethodPut,
		c.baseURL+path,
		reader,
		result,
	)
}

// ---------------- DELETE ----------------

func (c *Client) Delete(path string, result any) error {

	return c.request(
		http.MethodDelete,
		c.baseURL+path,
		nil,
		result,
	)
}

// ---------------- POST FORM ----------------

func (c *Client) PostForm(path string, values url.Values, result any) error {

	old := c.contentType

	c.contentType = ContentTypeForm

	defer func() {
		c.contentType = old
	}()

	return c.request(
		http.MethodPost,
		c.baseURL+path,
		strings.NewReader(values.Encode()),
		result,
	)
}

// ---------------- 核心请求 ----------------

func (c *Client) request(
	method string,
	url string,
	body io.Reader,
	result any,
) error {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	return c.Do(req, result)
}

// Do 发送自定义请求
func (c *Client) Do(req *http.Request, result any) error {

	// 设置默认 Header
	for k, v := range c.headers {
		if req.Header.Get(k) == "" {
			req.Header.Set(k, v)
		}
	}

	// 默认 Content-Type
	if req.Body != nil &&
		req.Header.Get("Content-Type") == "" {

		req.Header.Set("Content-Type", c.contentType)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {

		return errors.New(string(data))
	}

	if result == nil {
		return nil
	}

	switch v := result.(type) {

	case *[]byte:
		*v = data
		return nil

	case *string:
		*v = string(data)
		return nil

	default:
		return json.Unmarshal(data, result)
	}
}
