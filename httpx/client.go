package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	defaultPostHeader = map[string]string{"Content-Type": "application/json"}
	defaultTimeout    = time.Second * 30
)

type Client struct {
	// Method specifies the HTTP method (GET, POST, PUT, etc.).
	// For client requests, an empty string means GET.
	Method string
	// Header contains the request header fields either received
	Header http.Header
	// Body is the request's body.
	Body io.Reader
	// Timeout specifies a time limit for requests made by this
	Timeout time.Duration
	// Adding query parameters to URL
	Query map[string]interface{}
}

type ClientOptionFunc func(*Client) error

// SetMethod sets the request method
func SetMethod(method string) ClientOptionFunc {
	return func(c *Client) error {
		c.Method = strings.ToUpper(method)
		return nil
	}
}

// SetHeader sets the request header
func SetHeader(header interface{}) ClientOptionFunc {
	return func(c *Client) error {
		switch h := header.(type) {
		case map[string]string:
			for k, v := range h {
				c.Header.Set(k, v)
			}
		case http.Header:
			for k := range h {
				c.Header.Set(k, h.Get(k))
			}
		default:
			return errors.New("set header unknown type")
		}
		return nil
	}
}

// SetBody sets the request body
func SetBody(body interface{}) ClientOptionFunc {
	return func(c *Client) error {
		switch b := body.(type) {
		case []byte:
			c.Body = bytes.NewReader(b)
		default:
			t, err := json.Marshal(b)
			if err != nil {
				return fmt.Errorf("body encode error, %s", err)
			}
			c.Body = bytes.NewReader(t)
		}
		return nil
	}
}

// SetTimeout sets the request timeout
func SetTimeout(timeout time.Duration) ClientOptionFunc {
	return func(c *Client) error {
		c.Timeout = timeout
		return nil
	}
}

// SetQuery sets the request query parameters
func SetQuery(query map[string]interface{}) ClientOptionFunc {
	return func(c *Client) error {
		c.Query = query
		return nil
	}
}

// Call sends an HTTP request and returns an HTTP response
// Note: When resp is nil, you have to do it manually `response.Body.Close()`
func Call(url string, respBody interface{}, options ...ClientOptionFunc) (r *http.Response, err error) {
	// Set up the client
	c := &Client{Header: make(http.Header), Timeout: defaultTimeout}

	// Run the options on it
	for _, option := range options {
		if err = option(c); err != nil {
			return
		}
	}

	// Set default header
	if c.Method == "POST" {
		for k, v := range defaultPostHeader {
			if len(c.Header.Get(k)) == 0 {
				c.Header.Set(k, v)
			}
		}
	}

	// New request
	req, err := http.NewRequest(c.Method, url, c.Body)
	if err != nil {
		return
	}

	// Set request header
	for k := range c.Header {
		req.Header.Set(k, c.Header.Get(k))
	}

	// Set query parameters to URL
	if req.URL != nil && c.Query != nil {
		q := req.URL.Query()
		for k, v := range c.Query {
			q.Add(k, fmt.Sprintf("%v", v))
		}
		req.URL.RawQuery = q.Encode()
	}

	// Client
	client := http.Client{Timeout: c.Timeout}
	r, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	if respBody == nil {
		return
	}
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err == nil && len(b) > 0 {
		err = json.Unmarshal(b, respBody)
	}
	return
}
