package cmd

import (
	"strings"

	"com.lisitede.backstage.vikunja/framework"
	internalVikunja "com.lisitede.backstage.vikunja/internal/vikunja"
	"github.com/spf13/cobra"
)

var apiRouter Router

var apiFlags struct {
	context string
	name    string
	color   string
}

var apiCmd = &cobra.Command{
	Use:   "api <method> <path>",
	Short: "Vikunja Http API wrapper",
	Long: `Vikunja Http API wrapper by RESTful path.
Examples:
	backstage-vikunja api GET /version
	backstage-vikunja api GET /projects
	backstage-vikunja api GET /labels
	backstage-vikunja api POST /labels --name "important" --color "ff0000"
	backstage-vikunja api GET /projects/42/tasks
	backstage-vikunja api GET /projects/42/views
	backstage-vikunja api GET /projects/42/views/7/buckets
	backstage-vikunja api PUT /projects/42/views/7/buckets/9/tasks/123
	backstage-vikunja api GET /tasks/123
	backstage-vikunja api GET /tasks/123/comments
	backstage-vikunja api POST /tasks/123/comments --context "comment text"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		method := strings.ToUpper(args[0])
		path := args[1]

		argsMap := map[string]string{
			"context": apiFlags.context,
			"name":    apiFlags.name,
			"color":   apiFlags.color,
		}

		result, err := apiRouter.Invoke(method, path, argsMap)
		if err != nil {
			return outputError(err)
		}

		printResult(result)
		return nil
	},
}

func init() {
	apiRouter.Verb("GET", "/projects", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalVikunja.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListProjects()
	})

	apiRouter.Verb("POST", "/projects", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		name := strings.TrimSpace(args["name"])
		if name == "" {
			return nil, framework.InvalidFormatException("POST /projects requires --name")
		}

		adapter, err := internalVikunja.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.CreateProject(name)
	})

	apiRouter.Verb("GET", "/labels", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalVikunja.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListLabels()
	})

	apiRouter.Verb("POST", "/labels", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		name := strings.TrimSpace(args["name"])
		if name == "" {
			return nil, framework.InvalidFormatException("POST /labels requires --name")
		}
		color := strings.TrimSpace(args["color"])
		if color == "" {
			return nil, framework.InvalidFormatException("POST /labels requires --color")
		}

		adapter, err := internalVikunja.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.CreateLabel(name, color)
	})

	apiRouter.Verb("GET", "/projects/:projectId/tasks", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalVikunja.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListProjectTasks(params["projectId"])
	})

	apiRouter.Verb("GET", "/projects/:projectId/views", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalVikunja.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListProjectViews(params["projectId"])
	})

	apiRouter.Verb("GET", "/projects/:projectId/views/:viewId/buckets", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalVikunja.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListProjectViewBuckets(params["projectId"], params["viewId"])
	})

	apiRouter.Verb("PUT", "/projects/:projectId/views/:viewId/buckets/:bucketId/tasks/:taskId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalVikunja.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.MoveTaskToBucket(params["projectId"], params["viewId"], params["bucketId"], params["taskId"])
	})

	apiRouter.Verb("GET", "/tasks/:taskId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalVikunja.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.GetTask(params["taskId"])
	})

	apiRouter.Verb("GET", "/tasks/:taskId/comments", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalVikunja.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListTaskComments(params["taskId"])
	})

	apiRouter.Verb("POST", "/tasks/:taskId/comments", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		context := strings.TrimSpace(args["context"])
		if context == "" {
			return nil, framework.InvalidFormatException("POST /tasks/:taskId/comments requires --context")
		}

		adapter, err := internalVikunja.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.CreateTaskComment(params["taskId"], context)
	})

	apiRouter.Verb("GET", "/version", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalVikunja.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.GetInfo()
	})

	apiCmd.Flags().StringVar(&apiFlags.context, "context", "", "comment text")
	apiCmd.Flags().StringVar(&apiFlags.name, "name", "", "project or label title")
	apiCmd.Flags().StringVar(&apiFlags.color, "color", "", "label hex color")
	apiCmd.Flags().SortFlags = false
	rootCmd.AddCommand(apiCmd)
}
