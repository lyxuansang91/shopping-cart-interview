package vgs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RequestOptions holds the options for making a request
type RequestOptions struct {
	Method      string
	URL         string
	Headers     map[string]string
	Body        interface{}
	ContentType string
}

// DoRequest makes a request through VGS with the given options
func (c *Client) DoRequest(opts RequestOptions) ([]byte, error) {
	var bodyReader io.Reader

	if opts.Body != nil {
		bodyBytes, err := json.Marshal(opts.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(opts.Method, opts.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default content type if not provided
	if opts.ContentType == "" {
		opts.ContentType = "application/json"
	}
	req.Header.Set("Content-Type", opts.ContentType)

	// Set additional headers
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// DoJSONRequest makes a request and unmarshals the response into the provided result
func (c *Client) DoJSONRequest(opts RequestOptions, result interface{}) error {
	body, err := c.DoRequest(opts)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
