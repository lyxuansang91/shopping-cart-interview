package vgs

// See https://www.verygoodsecurity.com/docs/guides/outbound-connection for more information

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	VGS_USERNAME   = "USfGoGJ6TARLMEuW21MrF76V"
	VGS_PASSWORD   = "9d2c9b0b-4e83-4f8c-8dc3-37b97cdf17dd"
	VGS_PROXY_HOST = "tnttbbrd2hh.SANDBOX.verygoodproxy.com:8443"
	VGS_PEM_CERT   = `-----BEGIN CERTIFICATE-----
MIID2TCCAsGgAwIBAgIHAN4Gs/LGhzANBgkqhkiG9w0BAQ0FADB5MSQwIgYDVQQD
DBsqLnNhbmRib3gudmVyeWdvb2Rwcm94eS5jb20xITAfBgNVBAoMGFZlcnkgR29v
ZCBTZWN1cml0eSwgSW5jLjEuMCwGA1UECwwlVmVyeSBHb29kIFNlY3VyaXR5IC0g
RW5naW5lZXJpbmcgVGVhbTAgFw0xNjAyMDkyMzUzMzZaGA8yMTE3MDExNTIzNTMz
NloweTEkMCIGA1UEAwwbKi5zYW5kYm94LnZlcnlnb29kcHJveHkuY29tMSEwHwYD
VQQKDBhWZXJ5IEdvb2QgU2VjdXJpdHksIEluYy4xLjAsBgNVBAsMJVZlcnkgR29v
ZCBTZWN1cml0eSAtIEVuZ2luZWVyaW5nIFRlYW0wggEiMA0GCSqGSIb3DQEBAQUA
A4IBDwAwggEKAoIBAQDI3ukHpxIlDCvFjpqn4gAkrQVdWll/uI0Kv3wirwZ3Qrpg
BVeXjInJ+rV9r0ouBIoY8IgRLak5Hy/tSeV6nAVHv0t41B7VyoeTAsZYSWU11deR
DBSBXHWH9zKEvXkkPdy9tgHnvLIzui2H59OPljV7z3sCLguRIvIIw8djaV9z7FRm
KRsfmYHKOBlSO4TlpfXQg7jQ5ds65q8FFGvTB5qAgLXS8W8pvdk8jccmuzQXFUY+
ZtHgjThg7BHWWUn+7m6hQ6iHHCj34Qu69F8nLamd+KJ//14lukdyKs3AMrYsFaby
k+UGemM/s2q3B+39B6YKaHao0SRzSJC7qDwbWPy3AgMBAAGjZDBiMB0GA1UdDgQW
BBRWlIRrE2p2P018VTzTb6BaeOFhAzAPBgNVHRMBAf8EBTADAQH/MAsGA1UdDwQE
AwIBtjAjBgNVHSUEHDAaBggrBgEFBQcDAQYIKwYBBQUHAwIGBFUdJQAwDQYJKoZI
hvcNAQENBQADggEBAGWxLFlr0b9lWkOLcZtR9IDVxDL9z+UPFEk70D3NPaqXkoE/
TNNUkXgS6+VBA2G8nigq2Yj8qoIM+kTXPb8TzWv+lrcLm+i+4AShKVknpB15cC1C
/NJfyYGRW66s/w7HNS20RmrdN+bWS0PA4CVLXdGzUJn0PCsfsS+6Acn7RPAE+0A8
WB7JzXWi8x9mOJwiOhodp4j41mv+5eHM0reMh6ycuYbjquDNpiNnsLztk6MGsgAP
5C59drQWJU47738BcfbByuSTYFog6zNYCm7ACqbtiwvFTwjneNebOhsOlaEAHjup
d4QBqYVs7pzkhNNp9oUvv4wGf/KJcw5B9E6Tpfk=
-----END CERTIFICATE-----`
)

func TestNewClient(t *testing.T) {
	// Test data

	tests := []struct {
		name          string
		config        Config
		expectedError bool
	}{
		{
			name: "valid configuration",
			config: Config{
				Username:           VGS_USERNAME,
				Password:           VGS_PASSWORD,
				ProxyHost:          VGS_PROXY_HOST,
				CertPEM:            VGS_PEM_CERT,
				InsecureSkipVerify: true,
			},
			expectedError: false,
		},
		{
			name: "invalid certificate",
			config: Config{
				Username:           VGS_USERNAME,
				Password:           VGS_PASSWORD,
				ProxyHost:          VGS_PROXY_HOST,
				CertPEM:            "invalid cert",
				InsecureSkipVerify: true,
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, "https://"+VGS_USERNAME+":"+VGS_PASSWORD+"@"+VGS_PROXY_HOST, os.Getenv("HTTPS_PROXY"))
			}
		})
	}
}

func TestDoRequest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-value", r.Header.Get("Test-Header"))

		// Read and verify request body
		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "test-value", body["test_key"])

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer server.Close()

	// Create test client with minimal configuration
	client, err := NewClient(Config{
		Username:           VGS_USERNAME,
		Password:           VGS_PASSWORD,
		ProxyHost:          VGS_PROXY_HOST,
		CertPEM:            VGS_PEM_CERT,
		InsecureSkipVerify: true,
	})
	require.NoError(t, err)

	// Test request
	opts := RequestOptions{
		Method:  "POST",
		URL:     server.URL,
		Headers: map[string]string{"Test-Header": "test-value"},
		Body:    map[string]string{"test_key": "test-value"},
	}

	// Test DoRequest
	body, err := client.DoRequest(opts)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "success")

	// Test DoJSONRequest
	var result map[string]string
	err = client.DoJSONRequest(opts, &result)
	assert.NoError(t, err)
	assert.Equal(t, "success", result["status"])
}

func TestDoRequestError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
	}))
	defer server.Close()

	// Create test client
	client, err := NewClient(Config{
		Username:           VGS_USERNAME,
		Password:           VGS_PASSWORD,
		ProxyHost:          VGS_PROXY_HOST,
		CertPEM:            VGS_PEM_CERT,
		InsecureSkipVerify: true,
	})
	require.NoError(t, err)

	// Test request with error response
	opts := RequestOptions{
		Method: "GET",
		URL:    server.URL,
	}

	_, err = client.DoRequest(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "400")
	assert.Contains(t, err.Error(), "bad request")
}
