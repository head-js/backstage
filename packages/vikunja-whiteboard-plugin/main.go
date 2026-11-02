//go:build ignore

package main

import (
	"encoding/json"
	"strconv"

	"code.vikunja.io/api/pkg/log"
	"code.vikunja.io/api/pkg/plugins"
	"github.com/labstack/echo/v5"
)

// WhiteboardPlugin implements the lifecycle and authenticated router capabilities.
type WhiteboardPlugin struct {
	httpClient *vikunjaHTTPClient
}

func (p *WhiteboardPlugin) Name() string {
	return "whiteboard"
}

func (p *WhiteboardPlugin) Version() string {
	return "0.1.0"
}

type whiteboardInitializationError struct {
	cause error
}

func (e *whiteboardInitializationError) Error() string {
	return "whiteboard initialization failed: " + e.cause.Error()
}

func (e *whiteboardInitializationError) Unwrap() error {
	return e.cause
}

func (p *WhiteboardPlugin) Init() error {
	client, err := newVikunjaHTTPClient()
	if err != nil {
		p.httpClient = nil
		return &whiteboardInitializationError{cause: err}
	}

	p.httpClient = client
	log.Infof("whiteboard initialized")
	return nil
}

func (p *WhiteboardPlugin) Shutdown() error {
	log.Infof("whiteboard shut down")
	return nil
}

func (p *WhiteboardPlugin) RegisterAuthenticatedRoutes(g *echo.Group) {
	g.GET("/whiteboard/ping", handlePing)
	g.GET("/whiteboard/tasks/:taskId/description", handleTaskDescription)
}

func handlePing(c *echo.Context) error {
	return c.JSON(httpStatusOK, map[string]interface{}{
		"code":    0,
		"message": "ok",
	})
}

func handleTaskDescription(c *echo.Context) error {
	taskID, err := strconv.ParseInt(c.Param("taskId"), 10, 64)
	if err != nil || taskID < 1 {
		return echo.NewHTTPError(httpStatusBadRequest, "Invalid task ID")
	}
	if c.Request().URL.RawQuery != "" {
		return echo.NewHTTPError(httpStatusBadRequest, "Query parameters are not supported")
	}

	authorization := c.Request().Header.Get("Authorization")
	if authorization == "" {
		return echo.NewHTTPError(httpStatusUnauthorized, "Authorization is required")
	}
	if singleton.httpClient == nil {
		return echo.NewHTTPError(httpStatusServiceUnavailable, "Whiteboard HTTP client is not initialized")
	}

	statusCode, contentType, body, err := singleton.httpClient.get(
		c.Request().Context(),
		"/tasks/"+strconv.FormatInt(taskID, 10)+"?format=markdown",
		authorization,
	)
	if err != nil {
		if isVikunjaHTTPTimeout(err) {
			return echo.NewHTTPError(httpStatusGatewayTimeout, "Vikunja API request timed out")
		}
		return echo.NewHTTPError(httpStatusBadGateway, "Vikunja API request failed")
	}
	if statusCode >= httpStatusMultipleChoices && statusCode < httpStatusBadRequest {
		return echo.NewHTTPError(httpStatusBadGateway, "Vikunja API returned a redirect")
	}
	if statusCode >= httpStatusBadRequest {
		return c.Blob(statusCode, contentType, body)
	}
	if statusCode != httpStatusOK {
		return echo.NewHTTPError(httpStatusBadGateway, "Vikunja API returned an unexpected status")
	}

	var task struct {
		ID          *int64  `json:"id"`
		Description *string `json:"description"`
	}
	if err := json.Unmarshal(body, &task); err != nil {
		return echo.NewHTTPError(httpStatusBadGateway, "Vikunja API returned invalid task JSON")
	}
	if task.ID == nil || *task.ID != taskID {
		return echo.NewHTTPError(httpStatusBadGateway, "Vikunja API returned a mismatched task ID")
	}
	if task.Description == nil {
		return echo.NewHTTPError(httpStatusBadGateway, "Vikunja API response is missing a valid description")
	}

	return c.JSON(httpStatusOK, map[string]interface{}{
		"task_id":     *task.ID,
		"description": *task.Description,
	})
}

var singleton = &WhiteboardPlugin{}

func NewPlugin() plugins.Plugin {
	return singleton
}

func NewAuthenticatedRouterPlugin() plugins.AuthenticatedRouterPlugin {
	return singleton
}
