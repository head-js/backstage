package cmd

import (
	internalAgent "com.lisitede.backstage/internal/agent"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Report CLI invocation context for AI agents.",
	Long: `Report CLI invocation context for AI agents.

Describes the caller — shell, IDE, session, process tree, environment
variables — so agents can adapt to how and where they were invoked.

Examples:
  backstage agent hasshin`,
}

var hasshinCmd = &cobra.Command{
	Use:   "hasshin",
	Short: "Identify the calling shell or environment (macOS).",
	Long: `Identify the calling shell or environment.

Hasshin (発進) prints the detected caller to stdout and dumps diagnostic
signals (process tree, stdio types, env vars) to stderr. Currently macOS
only; output format may still change.

Examples:
  backstage agent hasshin`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		shell, err := internalAgent.Hasshin()
		if err != nil {
			return outputError(err)
		}
		printResult(shell)
		return nil
	},
}

func init() {
	agentCmd.AddCommand(hasshinCmd)
	rootCmd.AddCommand(agentCmd)
}
