package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "backstage-linear",
	Short: "A CLI tool to manage Plan / Phase / Task.",
	Long: `A CLI tool to manage Plan / Phase / Task. The backend uses Linear API.
The root command provides a brief healthcheck and has no real functions.`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func init() {
	defaultHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd.Parent() == nil {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n", cmd.Long)
			fmt.Fprintf(cmd.OutOrStdout(), "Flags:\n%s\n", cmd.LocalFlags().FlagUsages())
		} else {
			defaultHelp(cmd, args)
		}
	})
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// outputError outputs error in human or JSON format
func outputError(err error) error {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	return nil
}

// printResult 输出结果，支持 string 直接打印和 struct JSON 输出
func printResult(data interface{}) {
	if data == nil {
		fmt.Println("OK")
		return
	}

	if s, ok := data.(string); ok {
		fmt.Println(s)
		return
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("%v\n", data)
		return
	}
	fmt.Println(string(jsonBytes))
}
