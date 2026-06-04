package cmd

import (
	"fmt"
	"strings"

	internalWorkshop "com.lisitede.backstage.gitea/internal/workshop"
	"github.com/spf13/cobra"
)

var workshopRouter Router

var workshopFlags struct {
	context string
}

var workshopCmd = &cobra.Command{
	Use:   "workshop <method> <path>",
	Short: "Workshop issue operations",
	Long: `A CLI tool to read and comment on Workshop Issues (WORKSHOP-xxx) by RESTful-style path.
Examples:
	backstage-gitea workshop GET  /cms-mgr/PLAN-102/WORKSHOP-003
	backstage-gitea workshop POST /cms-mgr/PLAN-102/WORKSHOP-003 --context "./path/to/context.md"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		method := strings.ToUpper(args[0])
		path := args[1]

		argsMap := map[string]string{}
		if workshopFlags.context != "" {
			argsMap["context"] = workshopFlags.context
		}

		result, err := workshopRouter.Invoke(method, path, argsMap)
		if err != nil {
			return outputError(err)
		}

		printResult(result)
		return nil
	},
}

func init() {
	workshopRouter.Verb("GET", "/:appId/:planId/:workshopId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalWorkshop.GetWorkshop(params["appId"], params["planId"], params["workshopId"])
	})

	workshopRouter.Verb("POST", "/:appId/:planId/:workshopId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		if args["context"] == "" {
			return nil, fmt.Errorf("POST /:appId/:planId/:workshopId requires --context flag")
		}
		return internalWorkshop.AddWorkshop(params["appId"], params["planId"], params["workshopId"], args["context"])
	})

	workshopCmd.Flags().SortFlags = false
	workshopCmd.Flags().StringVar(&workshopFlags.context, "context", "", "path to context markdown file")
	rootCmd.AddCommand(workshopCmd)
}
