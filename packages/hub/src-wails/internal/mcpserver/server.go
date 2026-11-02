package mcpserver

import (
	"context"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	mcpsdk "github.com/mark3labs/mcp-go/server"
)

type Server struct {
	httpServer *mcpsdk.StreamableHTTPServer
}

func New(version string) *Server {
	server := mcpsdk.NewMCPServer(
		"backstage-hub",
		version,
		mcpsdk.WithToolCapabilities(false),
	)
	server.AddTool(mcp.NewTool(
		"ping",
		mcp.WithDescription("Check whether Backstage MCP is callable."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
	), func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("pong"), nil
	})

	return &Server{
		httpServer: mcpsdk.NewStreamableHTTPServer(server, mcpsdk.WithStateful(true)),
	}
}

func (s *Server) Handler() http.Handler {
	return s.httpServer
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
