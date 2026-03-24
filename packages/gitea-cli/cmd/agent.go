package cmd

import (
	"fmt"
	"strings"

	internalGitea "com.lisitede.backstage.gitea/internal/gitea"
	"com.lisitede.backstage.gitea/internal/plan"
	"github.com/spf13/cobra"
)

var agentRouter Router

var agentCmd = &cobra.Command{
	Use:   "agent <method> <path>",
	Short: "Agent-optimized API for Plan / Phase / Task operations.",
	Long: `Agent-optimized API for Plan / Phase / Task operations.
This command provides UX-optimized interfaces tailored for Agent usage patterns.`,
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
	agentRouter.Verb("GET-MARKDOWN", "/:appName/:planName", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return getMarkdownPlan(params["appName"], params["planName"])
	})

	rootCmd.AddCommand(agentCmd)
}

func getMarkdownPlan(appName string, planName string) (string, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return "", fmt.Errorf("failed to create adapter: %w", err)
	}

	repo, err := adapter.GetRepo(appName, planName)
	if err != nil {
		return "", err
	}

	planId, err := plan.ExtractPlanId(repo.Name)
	if err != nil {
		return "", err
	}

	issue, err := adapter.ShowIssueOfRepoByPrefix(appName, planName, planId)
	if err != nil {
		return "", err
	}

	return issue.Body, nil
}
