//go:build ignore

package main

import (
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
	return "2.7.1.0"
}

func (p *WhiteboardPlugin) Init() error {
	client, err := newVikunjaHTTPClient()
	if err != nil {
		p.httpClient = nil
		return err
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
	restful := &whiteboardRESTful{httpClient: p.httpClient}
	restful.registerRoutes(g)
}

var singleton = &WhiteboardPlugin{}

func NewPlugin() plugins.Plugin {
	return singleton
}

func NewAuthenticatedRouterPlugin() plugins.AuthenticatedRouterPlugin {
	return singleton
}
