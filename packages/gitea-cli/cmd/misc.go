package cmd

import (
	"fmt"
	"strings"

	internalMisc "com.lisitede.backstage.gitea/internal/misc"
	"github.com/spf13/cobra"
)

var miscRouter Router

// miscUpdateFlags 用于更新 Task 的 flags
var miscUpdateFlags struct {
	context string
	title   string
}

var miscCmd = &cobra.Command{
	Use:   "misc <method> <path>",
	Short: "Misc API for blame operations.",
	Long: `Misc API for blame operations.

Examples:
  # Create Blame
  backstage-gitea misc POST /cms-mgr/PLAN-102/blames --title "..." --context "..."

  # List Blames
  backstage-gitea misc GET /cms-mgr/PLAN-102/blames`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		method := strings.ToUpper(args[0])
		path := args[1]

		// 构建 args 映射
		argsMap := map[string]string{}
		if miscUpdateFlags.context != "" {
			argsMap["context"] = miscUpdateFlags.context
		}
		if miscUpdateFlags.title != "" {
			argsMap["title"] = miscUpdateFlags.title
		}

		// 调用路由
		result, err := miscRouter.Invoke(method, path, argsMap)
		if err != nil {
			return outputError(err)
		}

		printResult(result)
		return nil
	},
}

func init() {
	miscRouter.Verb("POST", "/:appId/:planId/blames", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		if args["title"] == "" {
			return nil, fmt.Errorf("POST /:appId/:planId/blames requires --title flag")
		}
		if args["context"] == "" {
			return nil, fmt.Errorf("POST /:appId/:planId/blames requires --context flag")
		}
		return internalMisc.CreateBlame(params["appId"], params["planId"], args["title"], args["context"])
	})

	miscRouter.Verb("GET", "/:appId/:planId/blames", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalMisc.ListBlame(params["appId"], params["planId"])
	})

	miscCmd.Flags().SortFlags = false
	miscCmd.Flags().StringVar(&miscUpdateFlags.title, "title", "", "")
	_ = miscCmd.Flags().MarkHidden("title")
	miscCmd.Flags().StringVar(&miscUpdateFlags.context, "context", "", "when create / update, provide the context or path to context markdown file")

	rootCmd.AddCommand(miscCmd)
}
