package vikunja

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	urlEnvironmentVariable   = "BACKSTAGE_VIKUNJA_URL"
	tokenEnvironmentVariable = "BACKSTAGE_VIKUNJA_TOKEN"
	apiVersionPath           = "/api/v2"
)

// Adapter provides business-neutral access to the Vikunja HTTP API.
type Adapter struct {
	baseURL string
	token   string
	client  *http.Client
}

// GetInfo returns the complete Vikunja instance information document.
func (a *Adapter) GetInfo() (any, error) {
	return a.Do(context.Background(), http.MethodGet, "/info", nil)
}

// NewAdapter creates an adapter from environment configuration.
func NewAdapter() (*Adapter, error) {
	baseURL := strings.TrimSpace(os.Getenv(urlEnvironmentVariable))
	if baseURL == "" {
		return nil, fmt.Errorf("%s is required", urlEnvironmentVariable)
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("%s must be an absolute URL", urlEnvironmentVariable)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("%s must use http or https", urlEnvironmentVariable)
	}
	if parsedURL.RawQuery != "" || parsedURL.Fragment != "" {
		return nil, fmt.Errorf("%s must not contain a query or fragment", urlEnvironmentVariable)
	}

	return &Adapter{
		baseURL: strings.TrimRight(baseURL, "/") + apiVersionPath,
		token:   strings.TrimSpace(os.Getenv(tokenEnvironmentVariable)),
		client:  &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// Do sends a JSON request to a path relative to the Vikunja v2 API root.
func (a *Adapter) Do(ctx context.Context, method, path string, body any) (any, error) {
	if strings.TrimSpace(method) == "" {
		return nil, fmt.Errorf("HTTP method is required")
	}
	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("API path must start with '/': %s", path)
	}

	var requestBody io.Reader
	if body != nil {
		encodedBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to encode request body: %w", err)
		}
		requestBody = bytes.NewReader(encodedBody)
	}

	request, err := http.NewRequestWithContext(ctx, method, a.baseURL+path, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if a.token != "" {
		request.Header.Set("Authorization", "Bearer "+a.token)
	}

	response, err := a.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Vikunja API request failed: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Vikunja API response: %w", err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		message := strings.TrimSpace(string(responseBody))
		if message == "" {
			message = response.Status
		}
		return nil, fmt.Errorf("Vikunja API returned %s: %s", response.Status, message)
	}

	if len(bytes.TrimSpace(responseBody)) == 0 {
		return nil, nil
	}

	var result any
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return string(responseBody), nil
	}
	return result, nil
}
