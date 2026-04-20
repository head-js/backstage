package cmd

import (
	"fmt"

	internalAgent "com.lisitede.backstage/internal/agent"
	"github.com/spf13/cobra"
)

var (
	logLevelDebug     bool
	logLevelDangerous bool
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
	Short: "Identify the calling shell or environment.",
	Long: `Identify the calling shell or environment.

Hasshin (発進) prints the detected caller to stdout and dumps diagnostic signals (process tree, stdio types, env vars) to stderr.

So the agent can adapt its behavior for interaction.

Examples:
  backstage agent hasshin`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var logLevel internalAgent.LogLevel
		if logLevelDangerous {
			logLevel = internalAgent.LogLevelDangerous
		} else if logLevelDebug {
			logLevel = internalAgent.LogLevelDebug
		} else {
			logLevel = internalAgent.LogLevelInfo
		}
		internalAgent.Hasshin(logLevel)
		return nil
	},
}

func init() {
	hasshinCmd.Flags().BoolVar(&logLevelDebug, "debug", false, "enable debug output")
	hasshinCmd.Flags().BoolVar(&logLevelDangerous, "dangerous", false, "enable dangerous operations")
	hasshinCmd.Flags().MarkHidden("debug")
	hasshinCmd.Flags().MarkHidden("dangerous")
	hasshinCmd.MarkFlagsMutuallyExclusive("debug", "dangerous")

	hasshinCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n", hasshinCmd.Long)
	})

	agentCmd.AddCommand(hasshinCmd)
	rootCmd.AddCommand(agentCmd)
}
