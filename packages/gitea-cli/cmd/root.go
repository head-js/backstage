package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "backstage-gitea",
	Short: "A CLI tool to manage Gitea resources",
	Long: `A CLI tool to manage Gitea resources including repos, milestones, and issues.
Supports JSON output for easy parsing by agents.`,
}

var jsonFlag bool

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().BoolVarP(&jsonFlag, "json", "j", false, "Output in JSON format")
}

// outputError outputs error in human or JSON format
func outputError(err error) error {
	if jsonFlag {
		json.NewEncoder(os.Stderr).Encode(map[string]string{
			"error": err.Error(),
		})
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	// 返回 nil 避免 cobra 再打印一次相同的错误
	return nil
}

// printResult 输出结果，支持 string 直接打印和 struct JSON 输出
func printResult(data interface{}) {
	if data == nil {
		fmt.Println("OK")
		return
	}

	// 如果是 string 类型，直接输出
	if s, ok := data.(string); ok {
		fmt.Println(s)
		return
	}

	// 其他类型走 JSON 序列化
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("%v\n", data)
		return
	}
	fmt.Println(string(jsonBytes))
}
