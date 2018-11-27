package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

// TLSOptions holds the TLS configurations information needed to create Gluster-Block rest client .
type TLSOptions struct {
	CaCertFile         string
	InsecureSkipVerify bool
}

// NewTLSConfig returns TLS configuration meant to be used by Gluster-Block rest client
func NewTLSConfig(opts *TLSOptions) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: opts.InsecureSkipVerify, // nolint
	}

	if opts.CaCertFile == "" || tlsConfig.InsecureSkipVerify {
		return tlsConfig, nil
	}

	caCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(opts.CaCertFile)

	if err != nil {
		return nil, err
	}

	if !caCertPool.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("failed to append cert from PEM file : %s", opts.CaCertFile)
	}

	tlsConfig.RootCAs = caCertPool
	return tlsConfig, nil
}
