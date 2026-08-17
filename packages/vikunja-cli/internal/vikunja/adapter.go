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
	"strconv"
	"strings"
	"time"

	"com.lisitede.backstage.vikunja/framework"
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

// CreateProject creates a new project with the given title and returns the created project.
func (a *Adapter) CreateProject(title string) (any, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, framework.InvalidFormatException("project title is required")
	}

	body := map[string]string{"title": title}
	return a.Do(context.Background(), http.MethodPost, "/projects", body)
}

// ListProjects returns the project collection visible to the authenticated user.
func (a *Adapter) ListProjects() (any, error) {
	return a.Do(context.Background(), http.MethodGet, "/projects", nil)
}

// ListProjectTasks returns the tasks belonging to a project.
func (a *Adapter) ListProjectTasks(projectID string) (any, error) {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return nil, framework.InvalidFormatException("project ID is required")
	}

	return a.Do(context.Background(), http.MethodGet, "/projects/"+url.PathEscape(projectID)+"/tasks", nil)
}

// ListProjectViews returns the views configured for a project.
func (a *Adapter) ListProjectViews(projectID string) (any, error) {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return nil, framework.InvalidFormatException("project ID is required")
	}

	return a.Do(context.Background(), http.MethodGet, "/projects/"+url.PathEscape(projectID)+"/views", nil)
}

// ListProjectViewBuckets returns the buckets configured for a project view.
func (a *Adapter) ListProjectViewBuckets(projectID, viewID string) (any, error) {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return nil, framework.InvalidFormatException("project ID is required")
	}

	viewID = strings.TrimSpace(viewID)
	if viewID == "" {
		return nil, framework.InvalidFormatException("view ID is required")
	}

	path := "/projects/" + url.PathEscape(projectID) + "/views/" + url.PathEscape(viewID) + "/buckets"
	return a.Do(context.Background(), http.MethodGet, path, nil)
}

// CreateBucket creates a kanban bucket in a project view.
func (a *Adapter) CreateBucket(projectID, viewID, title string, position float64) (any, error) {
	projectID, _, err := parseID("project", projectID)
	if err != nil {
		return nil, err
	}
	viewID, _, err = parseID("view", viewID)
	if err != nil {
		return nil, err
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, framework.InvalidFormatException("bucket title is required")
	}

	path := "/projects/" + projectID + "/views/" + viewID + "/buckets"
	body := map[string]any{"title": title, "position": position}
	return a.Do(context.Background(), http.MethodPost, path, body)
}

// UpdateBucket updates a kanban bucket's title and position.
func (a *Adapter) UpdateBucket(projectID, viewID, bucketID, title string, position float64) (any, error) {
	projectID, _, err := parseID("project", projectID)
	if err != nil {
		return nil, err
	}
	viewID, _, err = parseID("view", viewID)
	if err != nil {
		return nil, err
	}
	bucketID, _, err = parseID("bucket", bucketID)
	if err != nil {
		return nil, err
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, framework.InvalidFormatException("bucket title is required")
	}

	path := "/projects/" + projectID + "/views/" + viewID + "/buckets/" + bucketID
	body := map[string]any{"title": title, "position": position}
	return a.Do(context.Background(), http.MethodPut, path, body)
}

// MoveTaskToBucket places a task in a bucket of a project view.
func (a *Adapter) MoveTaskToBucket(projectID, viewID, bucketID, taskID string) (any, error) {
	projectID, _, err := parseID("project", projectID)
	if err != nil {
		return nil, err
	}
	viewID, _, err = parseID("view", viewID)
	if err != nil {
		return nil, err
	}
	bucketID, _, err = parseID("bucket", bucketID)
	if err != nil {
		return nil, err
	}
	_, taskIDValue, err := parseID("task", taskID)
	if err != nil {
		return nil, err
	}

	path := "/projects/" + projectID + "/views/" + viewID + "/buckets/" + bucketID + "/tasks"
	body := map[string]int64{"task_id": taskIDValue}
	return a.Do(context.Background(), http.MethodPut, path, body)
}

func parseID(name, value string) (string, int64, error) {
	value = strings.TrimSpace(value)
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return "", 0, framework.InvalidFormatException(name + " ID must be a positive integer")
	}
	return value, parsed, nil
}

// GetTask returns a task by its globally unique ID.
func (a *Adapter) GetTask(taskID string) (any, error) {
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, framework.InvalidFormatException("task ID is required")
	}

	return a.Do(context.Background(), http.MethodGet, "/tasks/"+url.PathEscape(taskID), nil)
}

// ListTaskComments returns the comments on a task.
func (a *Adapter) ListTaskComments(taskID string) (any, error) {
	taskID, _, err := parseID("task", taskID)
	if err != nil {
		return nil, err
	}

	return a.Do(context.Background(), http.MethodGet, "/tasks/"+taskID+"/comments", nil)
}

// CreateTaskComment adds a comment to a task.
func (a *Adapter) CreateTaskComment(taskID, comment string) (any, error) {
	taskID, _, err := parseID("task", taskID)
	if err != nil {
		return nil, err
	}

	comment = strings.TrimSpace(comment)
	if comment == "" {
		return nil, framework.InvalidFormatException("comment is required")
	}

	body := map[string]string{"comment": comment}
	return a.Do(context.Background(), http.MethodPost, "/tasks/"+taskID+"/comments", body)
}

// NewAdapter creates an adapter from environment configuration.
func NewAdapter() (*Adapter, error) {
	baseURL := strings.TrimSpace(os.Getenv(urlEnvironmentVariable))
	if baseURL == "" {
		return nil, framework.InvalidFormatException(urlEnvironmentVariable + " is required")
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, framework.InvalidFormatException(urlEnvironmentVariable + " must be an absolute URL")
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, framework.InvalidFormatException(urlEnvironmentVariable + " must use http or https")
	}
	if parsedURL.RawQuery != "" || parsedURL.Fragment != "" {
		return nil, framework.InvalidFormatException(urlEnvironmentVariable + " must not contain a query or fragment")
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
		return nil, framework.InvalidFormatException("HTTP method is required")
	}
	if !strings.HasPrefix(path, "/") {
		return nil, framework.InvalidFormatException("API path must start with '/': " + path)
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
