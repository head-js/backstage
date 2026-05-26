package cmd

import (
	"fmt"
	"strings"

	internalClarify "com.lisitede.backstage.gitea/internal/clarify"
	"github.com/spf13/cobra"
)

var clarifyRouter Router

var clarifyFlags struct {
	context string
}

var clarifyCmd = &cobra.Command{
	Use:   "clarify <method> <path>",
	Short: "Clarify issue operations",
	Long: `A CLI tool to read and comment on Clarify Issues (CLARIFY-xxx) by RESTful-style path.
Examples:
	backstage-gitea clarify GET  /cms-mgr/PLAN-102/CLARIFY-003
	backstage-gitea clarify POST /cms-mgr/PLAN-102/CLARIFY-003 --context "./path/to/context.md"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		method := strings.ToUpper(args[0])
		path := args[1]

		argsMap := map[string]string{}
		if clarifyFlags.context != "" {
			argsMap["context"] = clarifyFlags.context
		}

		result, err := clarifyRouter.Invoke(method, path, argsMap)
		if err != nil {
			return outputError(err)
		}

		printResult(result)
		return nil
	},
}

func init() {
	clarifyRouter.Verb("GET", "/:appId/:planId/:clarifyId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalClarify.GetClarify(params["appId"], params["planId"], params["clarifyId"])
	})

	clarifyRouter.Verb("POST", "/:appId/:planId/:clarifyId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		if args["context"] == "" {
			return nil, fmt.Errorf("POST /:appId/:planId/:clarifyId requires --context flag")
		}
		return internalClarify.AddClarify(params["appId"], params["planId"], params["clarifyId"], args["context"])
	})

	clarifyCmd.Flags().SortFlags = false
	clarifyCmd.Flags().StringVar(&clarifyFlags.context, "context", "", "path to context markdown file")
	rootCmd.AddCommand(clarifyCmd)
}
