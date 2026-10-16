package cmd

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"backstage/hub/internal/config"
	"backstage/hub/internal/server"
)

var (
	serveHost string
	servePort int
)

var rootCmd = &cobra.Command{
	Use:   "backstage-hub",
	Short: "Backstage hub service",
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP service in the foreground",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Default()
		cfg.Host = serveHost
		cfg.Port = servePort

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		return server.Run(ctx, cfg)
	},
}

func init() {
	serveCmd.Flags().StringVar(&serveHost, "host", "127.0.0.1", "address to listen on")
	serveCmd.Flags().IntVar(&servePort, "port", 20080, "port to listen on (0 = auto-assign)")
	rootCmd.AddCommand(serveCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
