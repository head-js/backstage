package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "backstage-vikunja",
	Short: "A CLI tool to access Vikunja.",
	Long: `A CLI tool to access Vikunja. The backend uses Vikunja Http API.
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

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

// outputError outputs an error and prevents Cobra from printing it twice.
func outputError(err error) error {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	return nil
}

// printResult prints strings directly and other values as JSON.
func printResult(data interface{}) {
	if data == nil {
		fmt.Println("OK")
		return
	}

	if value, ok := data.(string); ok {
		fmt.Println(value)
		return
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(data); err != nil {
		fmt.Printf("%v\n", data)
		return
	}
	fmt.Print(buf.String())
}
