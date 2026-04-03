package cmd

import (
	"strings"

	"com.lisitede.backstage.gitea/framework"
	internalAgent "com.lisitede.backstage.gitea/internal/agent"
	"github.com/spf13/cobra"
)

var agentRouter Router

// agentUpdateFlags 用于更新 Task 的 flags
var agentUpdateFlags struct {
	status  string
	context string
}

var agentCmd = &cobra.Command{
	Use:   "agent <method> <path>",
	Short: "Agent-optimized API for Plan / Phase / Task operations.",
	Long: `Agent-optimized API for Plan / Phase / Task operations.
This command provides UX-optimized interfaces tailored for Agent usage patterns.

Examples:
  # Get Plan
  backstage-gitea agent GET /cms-mgr/PLAN-102

  # Get Current Phase
  backstage-gitea agent HEAD /cms-mgr/PLAN-102/current-phase

  # Get Current Task
  backstage-gitea agent HEAD /cms-mgr/PLAN-102/PHASE-200/current-task

  # Get Phase Metadata
  backstage-gitea agent HEAD /cms-mgr/PLAN-102/PHASE-200

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

		// 构建 args 映射
		argsMap := map[string]string{}
		if agentUpdateFlags.status != "" {
			argsMap["status"] = agentUpdateFlags.status
		}
		if agentUpdateFlags.context != "" {
			argsMap["context"] = agentUpdateFlags.context
		}

		// 调用路由
		result, err := agentRouter.Invoke(method, path, argsMap)
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

	agentRouter.Verb("HEAD", "/:appId/:planId/current-phase", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.HeadCurrentPhase(params["appId"], params["planId"])
	})

	agentRouter.Verb("GET", "/:appId/:planId/current-phase", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.GetCurrentPhase(params["appId"], params["planId"])
	})

	agentRouter.Verb("HEAD", "/:appId/:planId/:phaseId/current-task", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.HeadCurrentTask(params["appId"], params["planId"], params["phaseId"])
	})

	agentRouter.Verb("GET", "/:appId/:planId/:phaseId/current-task", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.GetCurrentTask(params["appId"], params["planId"], params["phaseId"])
	})

	// current 要放在上面，从而优先匹配

	// Task
	agentRouter.Verb("HEAD", "/:appId/:planId/:phaseId/:taskId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.HeadTask(params["appId"], params["planId"], params["phaseId"], params["taskId"])
	})

	agentRouter.Verb("GET", "/:appId/:planId/:phaseId/:taskId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.GetTask(params["appId"], params["planId"], params["phaseId"], params["taskId"])
	})

	agentRouter.Verb("PUT", "/:appId/:planId/:phaseId/:taskId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.UpdateTask(params["appId"], params["planId"], params["phaseId"], params["taskId"], args["status"], args["context"])
	})

	// Phase
	agentRouter.Verb("HEAD", "/:appId/:planId/:phaseId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.HeadPhase(params["appId"], params["planId"], params["phaseId"])
	})

	agentRouter.Verb("GET", "/:appId/:planId/:phaseId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.GetPhase(params["appId"], params["planId"], params["phaseId"])
	})

	agentRouter.Verb("PUT", "/:appId/:planId/:phaseId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.UpdatePhase(params["appId"], params["planId"], params["phaseId"], args["status"], args["context"])
	})

	// Plan
	agentRouter.Verb("HEAD", "/:appId/:planId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return nil, framework.NotImplementedException("")
	})

	agentRouter.Verb("GET", "/:appId/:planId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return nil, framework.NotImplementedException("")
	})

	agentRouter.Verb("PUT", "/:appId/:planId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalAgent.UpdatePlan(params["appId"], params["planId"], args["status"], args["context"])
	})

	agentCmd.Flags().SortFlags = false
	agentCmd.Flags().StringVar(&agentUpdateFlags.context, "context", "", "when create / update, provide the context or path to context markdown file")
	agentCmd.Flags().StringVar(&agentUpdateFlags.status, "status", "", "when create / update, provide the status")

	rootCmd.AddCommand(agentCmd)
}
