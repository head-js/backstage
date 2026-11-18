package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"backstage.lisitede.com/hub/internal/config"
	mcpwhiteboard "backstage.lisitede.com/hub/internal/mcp-whiteboard"
	"backstage.lisitede.com/hub/internal/mcpserver"
)

type PingResponse struct {
	Body struct {
		Message string `json:"message"`
	}
}

type VersionResponse struct {
	Body struct {
		Version   string `json:"version"`
		APIPrefix string `json:"api_prefix"`
	}
}

func newAPI(mux *http.ServeMux) huma.API {
	humaConfig := huma.DefaultConfig("backstage-hub", config.Version)
	humaConfig.OpenAPIPath = "/api/openapi"
	humaConfig.DocsPath = ""
	humaConfig.SchemasPath = ""
	humaConfig.CreateHooks = nil

	api := humago.New(mux, humaConfig)

	huma.Register(api, huma.Operation{
		OperationID: "ping",
		Summary:     "Ping the hub",
		Method:      http.MethodGet,
		Path:        "/api/ping",
		Tags:        []string{"system"},
	}, func(ctx context.Context, input *struct{}) (*PingResponse, error) {
		resp := PingResponse{}
		resp.Body.Message = "pong"
		return &resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "version",
		Summary:     "Get hub version info",
		Method:      http.MethodGet,
		Path:        "/api/v1/version",
		Tags:        []string{"system"},
	}, func(ctx context.Context, input *struct{}) (*VersionResponse, error) {
		resp := VersionResponse{}
		resp.Body.Version = config.Version
		resp.Body.APIPrefix = config.APIPrefix
		return &resp, nil
	})

	return api
}

type Server struct {
	httpServer      *http.Server
	mcpServer       *mcpserver.Server
	whiteboardMCP   *mcpwhiteboard.Server
	listener        net.Listener
	serveDone       chan struct{}
	serveErr        error
	errCh           chan error
	shutdownTimeout time.Duration
	stopOnce        sync.Once
	stopDone        chan struct{}
}

func Start(cfg config.Config) (*Server, error) {
	mux := http.NewServeMux()
	newAPI(mux)
	mcpServer := mcpserver.New(config.Version)
	whiteboardMCP := mcpwhiteboard.New(config.Version, cfg.VikunjaURL, cfg.VikunjaToken)
	mux.Handle("/mcp-other", mcpServer.Handler())
	mux.Handle("/mcp-whiteboard", whiteboardMCP.Handler())

	listener, err := net.Listen("tcp", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)))
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if shutdownErr := mcpServer.Shutdown(ctx); shutdownErr != nil {
			log.Printf("backstage-hub MCP cleanup after listen failure: %v", shutdownErr)
		}
		if shutdownErr := whiteboardMCP.Shutdown(ctx); shutdownErr != nil {
			log.Printf("backstage-hub Whiteboard MCP cleanup after listen failure: %v", shutdownErr)
		}
		return nil, fmt.Errorf("listen: %w", err)
	}

	s := &Server{
		httpServer: &http.Server{
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		mcpServer:       mcpServer,
		whiteboardMCP:   whiteboardMCP,
		listener:        listener,
		serveDone:       make(chan struct{}),
		errCh:           make(chan error, 1),
		shutdownTimeout: 10 * time.Second,
		stopDone:        make(chan struct{}),
	}
	go s.serve()
	log.Printf("backstage-hub listening on http://%s", s.Addr())
	return s, nil
}

func (s *Server) serve() {
	defer close(s.serveDone)
	err := s.httpServer.Serve(s.listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.serveErr = err
		s.errCh <- err
	}
}

func (s *Server) Addr() string {
	return s.listener.Addr().String()
}

func (s *Server) Err() <-chan error {
	return s.errCh
}

func (s *Server) Running() bool {
	select {
	case <-s.serveDone:
		return false
	default:
		return true
	}
}

func (s *Server) Failure() error {
	select {
	case <-s.serveDone:
		return s.serveErr
	default:
		return nil
	}
}

func (s *Server) Stop() {
	s.stopOnce.Do(func() {
		go s.shutdown()
	})
	<-s.stopDone
}

func (s *Server) shutdown() {
	defer close(s.stopDone)
	mcpCtx, mcpCancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	if err := s.mcpServer.Shutdown(mcpCtx); err != nil {
		log.Printf("backstage-hub MCP shutdown: %v", err)
	}
	if err := s.whiteboardMCP.Shutdown(mcpCtx); err != nil {
		log.Printf("backstage-hub Whiteboard MCP shutdown: %v", err)
	}
	mcpCancel()

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("backstage-hub shutdown: %v; forcing close", err)
		_ = s.httpServer.Close()
	}
	<-s.serveDone
}
