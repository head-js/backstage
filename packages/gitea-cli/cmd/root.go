package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "backstage-gitea",
	Short: "A CLI tool to manage Gitea resources",
	Long: `A CLI tool to manage Gitea resources including repos, milestones, and issues.
Supports JSON output for easy parsing by agents.`,
}

var jsonFlag bool

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().BoolVarP(&jsonFlag, "json", "j", false, "Output in JSON format")
}

// outputError outputs error in human or JSON format
func outputError(err error) error {
	if jsonFlag {
		json.NewEncoder(os.Stderr).Encode(map[string]string{
			"error": err.Error(),
		})
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	return err
}
