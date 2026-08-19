package cmd

import (
	"strings"

	"com.lisitede.backstage.vikunja/framework"
	"com.lisitede.backstage.vikunja/internal/scrum"
	"github.com/spf13/cobra"
)

var scrumRouter Router

var scrumFlags struct {
	name   string
	filter string
	status string
}

var scrumCmd = &cobra.Command{
	Use:   "scrum <method> <path>",
	Short: "Scrum workflow commands",
	Long: `Scrum workflow commands.
Examples:
	backstage-vikunja scrum POST /workspaces --name "简单短标题"
	backstage-vikunja scrum GET /workspaces/42/tasks --filter "TO DO"
	backstage-vikunja scrum PATCH /workspaces/42/tasks/123 --status "IN DO"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		method := strings.ToUpper(args[0])
		path := args[1]

		argsMap := map[string]string{
			"name":   scrumFlags.name,
			"filter": scrumFlags.filter,
			"status": scrumFlags.status,
		}

		result, err := scrumRouter.Invoke(method, path, argsMap)
		if err != nil {
			return outputError(err)
		}

		printResult(result)
		return nil
	},
}

func init() {
	scrumRouter.Verb("POST", "/init", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return scrum.Initialize()
	})

	scrumRouter.Verb("GET", "/workspaces", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return scrum.ListWorkspaces()
	})

	scrumRouter.Verb("POST", "/workspaces", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		name := strings.TrimSpace(args["name"])
		if name == "" {
			return nil, framework.InvalidFormatException("POST /workspaces requires --name")
		}
		return scrum.CreateWorkspace(name)
	})

	scrumRouter.Verb("GET", "/workspaces/:workspaceId/tasks", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		filter := strings.TrimSpace(args["filter"])
		if filter == "" {
			return nil, framework.InvalidFormatException("GET /workspaces/:workspaceId/tasks requires --filter")
		}
		return scrum.ListTasks(params["workspaceId"], filter)
	})

	scrumRouter.Verb("PATCH", "/workspaces/:workspaceId/tasks/:taskId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		status := strings.TrimSpace(args["status"])
		if status == "" {
			return nil, framework.InvalidFormatException("PATCH /workspaces/:workspaceId/tasks/:taskId requires --status")
		}
		return scrum.UpdateTaskStatus(params["workspaceId"], params["taskId"], status)
	})

	scrumCmd.Flags().StringVar(&scrumFlags.name, "name", "", "workspace title")
	scrumCmd.Flags().StringVar(&scrumFlags.filter, "filter", "", "kanban bucket title, e.g. \"TO DO\"")
	scrumCmd.Flags().StringVar(&scrumFlags.status, "status", "", "task status, e.g. \"IN DO\"")
	scrumCmd.Flags().SortFlags = false
	rootCmd.AddCommand(scrumCmd)
}
