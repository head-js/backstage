package cmd

import (
	"fmt"
	"strings"

	"com.lisitede.backstage.linear/framework"
	internalLinear "com.lisitede.backstage.linear/internal/linear"
	"github.com/spf13/cobra"
)

// apiRouter API 命令路由管理器
var apiRouter Router

// apiCreateFlags 用于创建 Issue 的 flags
var apiCreateFlags struct {
	title   string
	context string
}

// apiCmd 通用 API 命令
var apiCmd = &cobra.Command{
	Use:   "api <method> <path>",
	Short: "Linear GraphQL API wrapper",
	Long: `Linear GraphQL API wrapper by RESTful-style path.
Examples:
	backstage-linear api GET /viewer
	backstage-linear api GET /teams
	backstage-linear api GET /teams/:teamKey/projects
	backstage-linear api GET /teams/:teamKey/projects/:projectName
	backstage-linear api GET /teams/:teamKey/projects/:projectName/issues
	backstage-linear api GET /issues/:issueId

Create Issue Example:
	backstage-linear api POST /teams/LSTD/projects/Blame/issues --title "BLAME-001 ..." --context ".context/BLAME-001.md"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		method := strings.ToUpper(args[0])
		path := args[1]

		argsMap := map[string]string{}
		if apiCreateFlags.title != "" {
			argsMap["title"] = apiCreateFlags.title
		}
		if apiCreateFlags.context != "" {
			argsMap["context"] = apiCreateFlags.context
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
	apiRouter.Verb("GET", "/viewer", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalLinear.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.GetViewer()
	})

	apiRouter.Verb("GET", "/teams", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalLinear.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListTeam()
	})

	apiRouter.Verb("GET", "/teams/:teamKey/projects", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalLinear.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListProjectOfTeam(params["teamKey"])
	})

	apiRouter.Verb("GET", "/teams/:teamKey/projects/:projectName", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalLinear.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ShowProject(params["teamKey"], params["projectName"])
	})

	apiRouter.Verb("GET", "/teams/:teamKey/projects/:projectName/issues", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalLinear.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListIssueOfProject(params["teamKey"], params["projectName"])
	})

	apiRouter.Verb("POST", "/teams/:teamKey/projects/:projectName/issues", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		if args["title"] == "" {
			return nil, fmt.Errorf("POST /teams/:teamKey/projects/:projectName/issues requires --title flag")
		}
		if args["context"] == "" {
			return nil, fmt.Errorf("POST /teams/:teamKey/projects/:projectName/issues requires --context flag")
		}

		description, err := framework.LoadContext(args["context"])
		if err != nil {
			return nil, err
		}

		adapter, err := internalLinear.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.CreateIssueInProject(params["teamKey"], params["projectName"], args["title"], description)
	})

	apiRouter.Verb("GET", "/issues/:issueId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalLinear.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.GetIssue(params["issueId"])
	})

	apiCmd.Flags().SortFlags = false
	apiCmd.Flags().StringVar(&apiCreateFlags.title, "title", "", "when create issue, provide the title")
	apiCmd.Flags().StringVar(&apiCreateFlags.context, "context", "", "when create issue, provide the context or path to context markdown file")
	rootCmd.AddCommand(apiCmd)
}
