//go:build ignore

package main

// Yaegi evaluates files in name order. Keep this file independent of later files.

import (
	vikunjabytes "bytes"
	vikunjacontext "context"
	vikunjajson "encoding/json"
	vikunjaerrors "errors"
	vikunjafmt "fmt"
	vikunjaio "io"
	http "net/http"
	vikunjaurl "net/url"
	vikunjastrings "strings"

	"code.vikunja.io/api/pkg/config"
)

type vikunjaHTTPClient struct {
	baseURL string
	client  *http.Client
}

func newVikunjaHTTPClient() (*vikunjaHTTPClient, error) {
	configuredURL := vikunjastrings.TrimSpace(config.ServicePublicURL.GetString())
	if configuredURL == "" {
		return nil, vikunjaerrors.New("service.publicurl is required")
	}

	parsedURL, err := vikunjaurl.Parse(configuredURL)
	if err != nil {
		return nil, err
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, vikunjafmt.Errorf("service.publicurl must use http or https")
	}
	if parsedURL.Host == "" {
		return nil, vikunjafmt.Errorf("service.publicurl must be an absolute URL with a host")
	}
	if parsedURL.User != nil {
		return nil, vikunjafmt.Errorf("service.publicurl must not contain userinfo")
	}
	if parsedURL.RawQuery != "" || parsedURL.ForceQuery || parsedURL.Fragment != "" {
		return nil, vikunjafmt.Errorf("service.publicurl must not contain a query or fragment")
	}

	rootPath := vikunjastrings.TrimRight(parsedURL.Path, "/")
	if vikunjastrings.HasSuffix(rootPath, "/api/v1") || vikunjastrings.HasSuffix(rootPath, "/api/v2") ||
		vikunjastrings.HasSuffix(rootPath, "/api/v1/plugins") {
		return nil, vikunjafmt.Errorf("service.publicurl must be the Vikunja service root URL")
	}
	parsedURL.Path = rootPath
	parsedURL.RawPath = ""

	return &vikunjaHTTPClient{
		baseURL: parsedURL.String(),
		client: &http.Client{
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}, nil
}

func (a *vikunjaHTTPClient) do(
	ctx vikunjacontext.Context,
	method, path string,
	body vikunjaio.Reader,
	authorization string,
) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, method, a.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", authorization)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := a.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := vikunjaio.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, vikunjafmt.Errorf("Vikunja API returned %s: %s", response.Status, responseBody)
	}
	return responseBody, nil
}

// listProjects does not depend on later-loaded types.
func (a *vikunjaHTTPClient) listProjects(ctx vikunjacontext.Context, authorization string) ([]byte, error) {
	return a.do(ctx, http.MethodGet, "/api/v2/projects?per_page=50", nil, authorization)
}

func (a *vikunjaHTTPClient) getTaskByIndex(ctx vikunjacontext.Context, projectID, index int64, authorization string) ([]byte, error) {
	path := vikunjafmt.Sprintf("/api/v2/projects/%d/tasks/by-index/%d?format=markdown", projectID, index)
	return a.do(ctx, http.MethodGet, path, nil, authorization)
}

func (a *vikunjaHTTPClient) listProjectTasks(ctx vikunjacontext.Context, projectID int64, authorization string) ([]byte, error) {
	path := vikunjafmt.Sprintf("/api/v2/projects/%d/tasks?per_page=50&format=markdown", projectID)
	return a.do(ctx, http.MethodGet, path, nil, authorization)
}

func (a *vikunjaHTTPClient) createTask(ctx vikunjacontext.Context, projectID int64, title, description, authorization string) ([]byte, error) {
	body, err := vikunjajson.Marshal(map[string]string{
		"title":       title,
		"description": description,
	})
	if err != nil {
		return nil, err
	}
	path := vikunjafmt.Sprintf("/api/v2/projects/%d/tasks?format=markdown", projectID)
	return a.do(ctx, http.MethodPost, path, vikunjabytes.NewReader(body), authorization)
}

func (a *vikunjaHTTPClient) deleteTask(ctx vikunjacontext.Context, taskID int64, authorization string) error {
	path := vikunjafmt.Sprintf("/api/v2/tasks/%d", taskID)
	_, err := a.do(ctx, http.MethodDelete, path, nil, authorization)
	return err
}

func (a *vikunjaHTTPClient) createTaskRelation(ctx vikunjacontext.Context, sourceID, targetID int64, authorization string) error {
	body, err := vikunjajson.Marshal(map[string]interface{}{
		"other_task_id": targetID,
		"relation_kind": "blocking",
	})
	if err != nil {
		return err
	}
	path := vikunjafmt.Sprintf("/api/v2/tasks/%d/relations", sourceID)
	_, err = a.do(ctx, http.MethodPost, path, vikunjabytes.NewReader(body), authorization)
	return err
}

func (a *vikunjaHTTPClient) deleteTaskRelation(ctx vikunjacontext.Context, sourceID, targetID int64, authorization string) error {
	path := vikunjafmt.Sprintf("/api/v2/tasks/%d/relations/blocking/%d", sourceID, targetID)
	_, err := a.do(ctx, http.MethodDelete, path, nil, authorization)
	return err
}
