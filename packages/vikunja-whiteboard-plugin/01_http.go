//go:build ignore

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"code.vikunja.io/api/pkg/config"
)

const (
	vikunjaAPIV2Path = "/api/v2"

	httpStatusOK                 = http.StatusOK
	httpStatusMultipleChoices    = http.StatusMultipleChoices
	httpStatusBadRequest         = http.StatusBadRequest
	httpStatusUnauthorized       = http.StatusUnauthorized
	httpStatusBadGateway         = http.StatusBadGateway
	httpStatusServiceUnavailable = http.StatusServiceUnavailable
	httpStatusGatewayTimeout     = http.StatusGatewayTimeout
)

type vikunjaHTTPClient struct {
	baseURL string
	client  *http.Client
}

func newVikunjaHTTPClient() (*vikunjaHTTPClient, error) {
	configuredURL := strings.TrimSpace(config.ServicePublicURL.GetString())
	if configuredURL == "" {
		return nil, errors.New("service.publicurl is required")
	}

	parsedURL, err := url.Parse(configuredURL)
	if err != nil {
		return nil, fmt.Errorf("service.publicurl must be a valid absolute URL: %w", err)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("service.publicurl must use http or https")
	}
	if parsedURL.Host == "" {
		return nil, fmt.Errorf("service.publicurl must be an absolute URL with a host")
	}
	if parsedURL.User != nil {
		return nil, fmt.Errorf("service.publicurl must not contain userinfo")
	}
	if parsedURL.RawQuery != "" || parsedURL.ForceQuery || parsedURL.Fragment != "" {
		return nil, fmt.Errorf("service.publicurl must not contain a query or fragment")
	}

	rootPath := strings.TrimRight(parsedURL.Path, "/")
	if strings.HasSuffix(rootPath, "/api/v1") || strings.HasSuffix(rootPath, "/api/v2") ||
		strings.HasSuffix(rootPath, "/api/v1/plugins") {
		return nil, fmt.Errorf("service.publicurl must be the Vikunja service root URL")
	}
	parsedURL.Path = rootPath + vikunjaAPIV2Path
	parsedURL.RawPath = ""

	return &vikunjaHTTPClient{
		baseURL: parsedURL.String(),
		client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}, nil
}

func (a *vikunjaHTTPClient) get(
	ctx context.Context,
	path string,
	authorization string,
) (statusCode int, contentType string, body []byte, err error) {
	if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") {
		return 0, "", nil, fmt.Errorf("Vikunja API path must be an internal absolute path")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, a.baseURL+path, nil)
	if err != nil {
		return 0, "", nil, fmt.Errorf("failed to create Vikunja API request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", authorization)

	response, err := a.client.Do(request)
	if err != nil {
		return 0, "", nil, fmt.Errorf("Vikunja API request failed: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, "", nil, fmt.Errorf("failed to read Vikunja API response: %w", err)
	}

	return response.StatusCode, response.Header.Get("Content-Type"), responseBody, nil
}

func isVikunjaHTTPTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var networkError net.Error
	return errors.As(err, &networkError) && networkError.Timeout()
}
