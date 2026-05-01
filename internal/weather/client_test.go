package weather

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_NilHTTPClient(t *testing.T) {
	c := NewClient(nil)

	require.NotNil(t, c)
	require.NotNil(t, c.(*client).httpClient)

	assert.Equal(t, 10*time.Second, c.(*client).httpClient.Timeout)
}

func TestClient_Get_Succes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"current": {
				"time": "2026-04-30T12:00",
				"temperature_2m": 14.3
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.Client())

	resp, err := client.get(context.Background(), server.URL)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestClient_Get_SuccesWithArgs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"current": {
				"time": "2026-04-30T12:00",
				"temperature_2m": 14.3
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.Client())

	s := make([]string, 0)
	s = append(s, server.URL[:len(server.URL)/2])
	s = append(s, server.URL[len(server.URL)/2:])
	resp, err := client.get(context.Background(), s[0], s[1])

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestClient_Get_RequestCreationFail(t *testing.T) {
	client := NewClient(nil)

	_, err := client.get(context.Background(), "http://[::1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "create weather request")
}

type errorRoundTripper struct{}

func (e *errorRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("network error")
}

func TestClient_Get_DoError(t *testing.T) {
	httpClient := &http.Client{
		Transport: &errorRoundTripper{},
	}

	c := NewClient(httpClient)

	_, err := c.get(context.Background(), "http://example.com")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "send weather request")
}

type errorReadCloser struct{}

func (e *errorReadCloser) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

func (e *errorReadCloser) Close() error {
	return nil
}

type readErrorRoundTripper struct{}

func (r *readErrorRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       &errorReadCloser{},
		Header:     make(http.Header),
	}, nil
}

func TestClient_Get_ReadBodyError(t *testing.T) {
	httpClient := &http.Client{
		Transport: &readErrorRoundTripper{},
	}

	c := NewClient(httpClient)

	_, err := c.get(context.Background(), "http://example.com")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "read weather response")
}

func TestClient_Get_StatusNotOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`bad request`))
	}))
	defer server.Close()

	c := NewClient(server.Client())

	_, err := c.get(context.Background(), server.URL)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "weather api returned status 400")
	assert.Contains(t, err.Error(), "bad request")
}
