package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"

	"backstage.lisitede.com/hub/internal/config"
	"backstage.lisitede.com/hub/internal/server"
)

//go:embed all:frontend/dist
var assets embed.FS

type Status struct {
	Version string `json:"version"`
	Addr    string `json:"addr"`
	BaseURL string `json:"baseUrl"`
	Running bool   `json:"running"`
}

type HubService struct {
	server *server.Server
}

func NewHubService(srv *server.Server) *HubService {
	return &HubService{server: srv}
}

func (h *HubService) GetStatus() Status {
	addr := h.server.Addr()
	return Status{
		Version: config.Version,
		Addr:    addr,
		BaseURL: "http://" + addr,
		Running: h.server.Running(),
	}
}

func main() {
	if err := run(); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Load()

	srv, err := server.Start(cfg)
	if err != nil {
		return fmt.Errorf("start http server: %w", err)
	}
	defer srv.Stop()

	dist, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		return fmt.Errorf("load frontend assets: %w", err)
	}

	app := application.New(application.Options{
		Name:        "Backstage Hub",
		Description: "Backstage hub",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(dist),
		},
		Services: []application.Service{
			application.NewService(NewHubService(srv)),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Windows: application.WindowsOptions{
			DisableQuitOnLastWindowClosed: false,
		},
		OnShutdown: srv.Stop,
	})

	go func() {
		if err := <-srv.Err(); err != nil {
			log.Printf("http server terminated unexpectedly: %v", err)
			app.Quit()
		}
	}()

	select {
	case err := <-srv.Err():
		return fmt.Errorf("http server terminated during startup: %w", err)
	default:
	}

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Backstage Hub",
		Width:  720,
		Height: 480,
	})

	runErr := app.Run()
	if err := srv.Failure(); err != nil {
		return fmt.Errorf("http server failed: %w", err)
	}
	return runErr
}
