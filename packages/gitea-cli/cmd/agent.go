package cmd

import (
	"strings"

	internalAgent "com.lisitede.backstage.gitea/internal/agent"
	"github.com/spf13/cobra"
)

var agentRouter Router

var agentCmd = &cobra.Command{
	Use:   "agent <method> <path>",
	Short: "Agent-optimized API for Plan / Phase / Task operations.",
	Long: `Agent-optimized API for Plan / Phase / Task operations.
This command provides UX-optimized interfaces tailored for Agent usage patterns.

Examples:
  # Get Plan
  backstage-gitea agent GET /cms-mgr/PLAN-102

  # Get Phase
  backstage-gitea agent GET /cms-mgr/PLAN-102/PHASE-200

	# Get Task Metadata
	backstage-gitea agent HEAD /cms-mgr/PLAN-102/PHASE-200/TASK-101
  # Get Task
  backstage-gitea agent GET /cms-mgr/PLAN-102/PHASE-200/TASK-101`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		method := strings.ToUpper(args[0])
		path := args[1]

		result, err := agentRouter.Invoke(method, path, map[string]string{})
		if err != nil {
			return outputError(err)
		}

		printResult(result)
		return nil
	},
}

func init() {
	agentRouter.Verb("GET", "/:appId/:planId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.GetPlan(params["appId"], params["planId"])
	})

	agentRouter.Verb("HEAD", "/:appId/:planId/:phaseId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.HeadPhase(params["appId"], params["planId"], params["phaseId"])
	})

	agentRouter.Verb("GET", "/:appId/:planId/:phaseId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.GetPhase(params["appId"], params["planId"], params["phaseId"])
	})

	agentRouter.Verb("HEAD", "/:appId/:planId/:phaseId/:taskId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.HeadTask(params["appId"], params["planId"], params["phaseId"], params["taskId"])
	})

	agentRouter.Verb("GET", "/:appId/:planId/:phaseId/:taskId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.GetTask(params["appId"], params["planId"], params["phaseId"], params["taskId"])
	})

	rootCmd.AddCommand(agentCmd)
}
