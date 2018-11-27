package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// Client is concrete instance implementing the GlusterBlockClient interface
type Client struct {
	httpClient *http.Client
	verb       string
	baseURL    string
	version    string
	path       string
	auth       struct {
		username string
		password string
	}
}

// NewClientWithOpts initializes a default GlusterBlockClient with host addr set as provided baseURL.
// It takes functors to modify its members while creating.
func NewClientWithOpts(baseURL string, optFuncs ...OptFuncs) (*Client, error) {
	client := &Client{
		httpClient: defaultHTTPClient(),
		baseURL:    baseURL,
	}

	for _, fn := range optFuncs {
		if err := fn(client); err != nil {
			return client, err
		}
	}
	client.wrapHTTPTransport()
	return client, nil
}

// wrapHTTPTransport wraps a http Roundtripper with all required round trippers
func (c *Client) wrapHTTPTransport() {
	var rt = c.httpClient.Transport

	if c.auth.username != "" && c.auth.password != "" {
		rt = NewBearerAuthRoundTrippers(c.auth.username, c.auth.password, rt)
	}
	c.httpClient.Transport = rt
}

// Verb will set http Method for a request
func (c *Client) Verb(method string) *Client {
	c.verb = method
	return c
}

// Path sets http path for a request
func (c *Client) Path(path string) *Client {
	c.path = path
	return c
}

// APIVersion sets api version for a request
func (c *Client) APIVersion(ver string) *Client {
	c.version = ver
	return c
}

// Do sends an HTTP request and decode the response in successV if status
// code lies in [200,400).
// An error is returned if server responds with error or status code
// does not lies in [200,400)
func (c *Client) Do(reqBody interface{}, successV interface{}) error {
	req, err := c.buildRequest(reqBody)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	defer closeRespBody(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		if successV == nil {
			return nil
		}
		return json.NewDecoder(resp.Body).Decode(successV)
	}

	return checkResponseError(resp)
}

func (c *Client) buildRequest(reqBody interface{}) (*http.Request, error) {
	var (
		body  io.Reader
		parts = strings.Split(c.baseURL, "://")
	)

	if len(parts) == 1 {
		c.baseURL = "http://" + c.baseURL
	}

	if reqBody != nil {
		reqBody, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(reqBody)
	}

	url := c.baseURL + "/" + c.version + c.path

	req, err := http.NewRequest(c.verb, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	return req, nil
}
