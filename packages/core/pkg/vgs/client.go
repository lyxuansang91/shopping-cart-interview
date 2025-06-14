package vgs

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// Config holds the configuration for the VGS client
type Config struct {
	// Username is the VGS proxy username
	Username string
	// Password is the VGS proxy password
	Password string
	// ProxyHost is the VGS proxy host (e.g., tnttbbrd2hh.SANDBOX.verygoodproxy.com:8443)
	ProxyHost string
	// CertPEM is the VGS certificate in PEM format
	CertPEM string
	// InsecureSkipVerify controls whether to verify the server's certificate
	InsecureSkipVerify bool
}

// Client is a wrapper around http.Client that uses VGS proxy
type Client struct {
	*http.Client
	config Config
}

// NewClient creates a new VGS client with the given configuration
func NewClient(config Config) (*Client, error) {
	// Construct the proxy URL
	proxyURL := fmt.Sprintf("https://%s:%s@%s",
		url.QueryEscape(config.Username),
		url.QueryEscape(config.Password),
		config.ProxyHost,
	)

	// Set the HTTPS proxy environment variable
	if err := os.Setenv("HTTPS_PROXY", proxyURL); err != nil {
		return nil, fmt.Errorf("failed to set HTTPS_PROXY: %w", err)
	}

	// Create a new certificate pool and add the VGS certificate
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM([]byte(config.CertPEM)) {
		return nil, fmt.Errorf("failed to append certificate to pool")
	}

	// Create the HTTP client with the VGS proxy configuration
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				RootCAs:            caCertPool,
				InsecureSkipVerify: config.InsecureSkipVerify,
			},
		},
	}

	return &Client{
		Client: client,
		config: config,
	}, nil
}

// NewDefaultClient creates a new VGS client with default configuration
func NewDefaultClient(username, password, proxyHost, certPEM string) (*Client, error) {
	config := Config{
		Username:           username,
		Password:           password,
		ProxyHost:          proxyHost,
		CertPEM:            certPEM,
		InsecureSkipVerify: true, // Default to true for development
	}
	return NewClient(config)
}
