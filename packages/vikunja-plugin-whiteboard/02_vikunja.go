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
	vikunjasort "sort"
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

func (a *vikunjaHTTPClient) getTask(ctx vikunjacontext.Context, taskId int64, authorization string) ([]byte, error) {
	path := vikunjafmt.Sprintf("/api/v2/tasks/%d?format=markdown", taskId)
	return a.do(ctx, http.MethodGet, path, nil, authorization)
}

func (a *vikunjaHTTPClient) listProjectTasks(ctx vikunjacontext.Context, projectId int64, authorization string) ([]byte, error) {
	path := vikunjafmt.Sprintf("/api/v2/projects/%d/tasks?per_page=50&format=markdown", projectId)
	return a.do(ctx, http.MethodGet, path, nil, authorization)
}

func (a *vikunjaHTTPClient) createTask(ctx vikunjacontext.Context, projectId int64, title, description, authorization string) ([]byte, error) {
	body, err := vikunjajson.Marshal(map[string]string{
		"title":       title,
		"description": description,
	})
	if err != nil {
		return nil, err
	}
	path := vikunjafmt.Sprintf("/api/v2/projects/%d/tasks?format=markdown", projectId)
	return a.do(ctx, http.MethodPost, path, vikunjabytes.NewReader(body), authorization)
}

// updateTaskFields writes only the supplied fields through the v2 bulk endpoint.
// The single-task update (POST /tasks/{id}) rewrites every column, so it would
// blank fields this plugin does not manage; the bulk endpoint with an explicit
// fields list is the true partial update.
func (a *vikunjaHTTPClient) updateTaskFields(ctx vikunjacontext.Context, taskId int64, values map[string]interface{}, authorization string) error {
	fields := make([]string, 0, len(values))
	for field := range values {
		fields = append(fields, field)
	}
	vikunjasort.Strings(fields)
	body, err := vikunjajson.Marshal(map[string]interface{}{
		"task_ids": []int64{taskId},
		"fields":   fields,
		"values":   values,
	})
	if err != nil {
		return err
	}
	_, err = a.do(ctx, http.MethodPut, "/api/v2/tasks/bulk", vikunjabytes.NewReader(body), authorization)
	return err
}

func (a *vikunjaHTTPClient) deleteTask(ctx vikunjacontext.Context, taskId int64, authorization string) error {
	path := vikunjafmt.Sprintf("/api/v2/tasks/%d", taskId)
	_, err := a.do(ctx, http.MethodDelete, path, nil, authorization)
	return err
}

func (a *vikunjaHTTPClient) createTaskRelation(ctx vikunjacontext.Context, sourceId, targetId int64, authorization string) error {
	body, err := vikunjajson.Marshal(map[string]interface{}{
		"other_task_id": targetId,
		"relation_kind": "blocking",
	})
	if err != nil {
		return err
	}
	path := vikunjafmt.Sprintf("/api/v2/tasks/%d/relations", sourceId)
	_, err = a.do(ctx, http.MethodPost, path, vikunjabytes.NewReader(body), authorization)
	return err
}

func (a *vikunjaHTTPClient) deleteTaskRelation(ctx vikunjacontext.Context, sourceId, targetId int64, authorization string) error {
	path := vikunjafmt.Sprintf("/api/v2/tasks/%d/relations/blocking/%d", sourceId, targetId)
	_, err := a.do(ctx, http.MethodDelete, path, nil, authorization)
	return err
}
