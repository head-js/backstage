package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"backstage/hub/internal/config"
)

const openAPISpec = `{
  "openapi": "3.1.0",
  "info": {
    "title": "backstage-hub",
    "version": "` + config.Version + `"
  },
  "paths": {
    "/api/ping": {
      "get": {
        "operationId": "ping",
        "summary": "Ping the hub",
        "tags": ["system"],
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {"message": {"type": "string"}},
                  "required": ["message"]
                }
              }
            }
          }
        }
      }
    },
    "/api/v1/version": {
      "get": {
        "operationId": "version",
        "summary": "Get hub version info",
        "tags": ["system"],
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "version": {"type": "string"},
                    "api_prefix": {"type": "string"}
                  },
                  "required": ["version", "api_prefix"]
                }
              }
            }
          }
        }
      }
    }
  }
}`

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/ping", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"message": "pong"})
	})

	mux.HandleFunc("GET /api/v1/version", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"version":    config.Version,
			"api_prefix": config.APIPrefix,
		})
	})

	mux.HandleFunc("GET /api/openapi", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, openAPISpec)
	})

	return mux
}

func Run(ctx context.Context, cfg config.Config) error {
	srv := &http.Server{Handler: newMux()}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	addr := listener.Addr().(*net.TCPAddr)
	log.Printf("backstage-hub listening on http://%s", addr.String())

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
		log.Printf("backstage-hub stopped")
		return nil
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return fmt.Errorf("serve: %w", err)
	}
}
