package cmd

import (
	"strings"

	internalVikunja "com.lisitede.backstage.vikunja/internal/vikunja"
	"github.com/spf13/cobra"
)

var apiRouter Router

var apiCmd = &cobra.Command{
	Use:   "api <method> <path>",
	Short: "Vikunja Http API wrapper",
	Long: `Vikunja Http API wrapper by RESTful path.
Examples:
	backstage-vikunja api GET /version`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		method := strings.ToUpper(args[0])
		path := args[1]

		argsMap := map[string]string{}

		result, err := apiRouter.Invoke(method, path, argsMap)
		if err != nil {
			return outputError(err)
		}

		printResult(result)
		return nil
	},
}

func init() {
	apiRouter.Verb("GET", "/version", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalVikunja.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.GetInfo()
	})

	apiCmd.Flags().SortFlags = false
	rootCmd.AddCommand(apiCmd)
}
