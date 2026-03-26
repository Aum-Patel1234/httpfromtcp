package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 25, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// check if valid token
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Multiple headers with duplicate keys (CSV merge)
	headers = NewHeaders()
	data = []byte(
		"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"HOST: example.com\r\n" + // duplicate
			"\r\n",
	)
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	// CSV merge happens
	assert.Equal(t, "localhost:42069,example.com", headers.Get("Host"))
	assert.Equal(t, "curl/7.81.0", headers.Get("User-Agent"))
	assert.Equal(t, len(data), n)
	assert.True(t, done)
}
