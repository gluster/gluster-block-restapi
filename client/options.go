package client

import (
	"fmt"
	"net/http"
	"time"
)

// OptFuncs receives a Client and overrides its members
// It will set the optional parameter of Client while creating
type OptFuncs func(*Client) error

// WithHTTPClient overrides http Client with specified one.
func WithHTTPClient(httpClient *http.Client) OptFuncs {
	return func(client *Client) error {
		if httpClient != nil {
			client.httpClient = httpClient
		}
		return nil
	}
}

// WithTLSConfig applies tls config to underlying http.Client Transport
func WithTLSConfig(tlsOpts *TLSOptions) OptFuncs {
	return func(client *Client) error {
		tlsConfig, err := NewTLSConfig(tlsOpts)
		if err != nil {
			return fmt.Errorf("failed to create tlsconfig, err: %s", err.Error())
		}
		if transport, ok := client.httpClient.Transport.(*http.Transport); ok {
			transport.TLSClientConfig = tlsConfig
			return nil
		}
		return fmt.Errorf("failed to apply tlsconfig on Transport : %T", client.httpClient.Transport)
	}
}

// WithAuth set username and password in client to be used for rest authentication
func WithAuth(username string, password string) OptFuncs {
	return func(client *Client) error {
		client.auth.username = username
		client.auth.password = password
		return nil
	}
}

// WithTimeOut overrides Client timeout with specified one
func WithTimeOut(timeout time.Duration) OptFuncs {
	return func(client *Client) error {
		client.httpClient.Timeout = timeout
		return nil
	}
}
