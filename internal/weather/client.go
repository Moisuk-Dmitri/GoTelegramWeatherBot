package weather

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type clientI interface {
	get(context.Context, string, ...string) ([]byte, error)
}

type client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) clientI {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return &client{
		httpClient: httpClient,
	}
}

func (c *client) get(ctx context.Context, baseURL string, args ...string) ([]byte, error) {
	jsonData := make([]byte, 0)

	var url strings.Builder
	url.WriteString(baseURL)
	for _, arg := range args {
		url.WriteString(arg)
	}

	req, err := http.NewRequestWithContext(ctx, "GET",
		url.String(),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("create weather request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send weather request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read weather response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("weather api returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
