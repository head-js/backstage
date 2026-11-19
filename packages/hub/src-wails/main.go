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
	"backstage.lisitede.com/hub/internal/vault"
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

type VaultService struct {
	store *vault.Store
}

func NewVaultService(store *vault.Store) *VaultService {
	return &VaultService{store: store}
}

func (v *VaultService) AddSecret() string {
	return v.store.AddSecret()
}

func (v *VaultService) RemoveKey(secret string) {
	v.store.RemoveKey(secret)
}

func (v *VaultService) ListRecords() []vault.Record {
	return v.store.All()
}

func main() {
	if err := run(); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	vaultStore := vault.New()
	srv, err := server.Start(cfg, vaultStore)
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
			application.NewService(NewVaultService(vaultStore)),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		Windows: application.WindowsOptions{
			DisableQuitOnLastWindowClosed: true,
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

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Backstage Hub",
		Width:  720,
		Height: 480,
	})

	setupTray(app, window)

	runErr := app.Run()
	if err := srv.Failure(); err != nil {
		return fmt.Errorf("http server failed: %w", err)
	}
	return runErr
}
