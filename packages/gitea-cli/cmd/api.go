package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"code.gitea.io/sdk/gitea"
	internalGitea "com.lisitede.backstage.gitea/internal/gitea"
	"github.com/spf13/cobra"
	"github.com/ucarion/urlpath"
)

var apiRoutes = []Route{
	{Method: "GET", Pattern: "/repos", Matcher: urlpath.New("/repos")},
	{Method: "POST", Pattern: "/repos", Matcher: urlpath.New("/repos")},
	{Method: "GET", Pattern: "/repos/:username/:repoName", Matcher: urlpath.New("/repos/:username/:repoName")},
	{Method: "GET", Pattern: "/repos/:username/:repoName/issues", Matcher: urlpath.New("/repos/:username/:repoName/issues")},
	{Method: "GET", Pattern: "/repos/:username/:repoName/milestones", Matcher: urlpath.New("/repos/:username/:repoName/milestones")},
	{Method: "GET", Pattern: "/repos/:username/:repoName/milestones/:milestonePrefix/issues", Matcher: urlpath.New("/repos/:username/:repoName/milestones/:milestonePrefix/issues")},
	{Method: "GET", Pattern: "/version", Matcher: urlpath.New("/version")},
	{Method: "GET", Pattern: "/users/:username/repos", Matcher: urlpath.New("/users/:username/repos")},
}

// apiCmd 通用 API 命令
var apiCmd = &cobra.Command{
	Use:   "api <method> <path>",
	Short: "Call Gitea API",
	Long: `Call Gitea API by RESTful path.
Examples:
	backstage-gitea api GET /repos
	backstage-gitea api GET /repos/123
  backstage-gitea api POST /repos`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return outputError(err)
		}

		method := strings.ToUpper(args[0])
		path := args[1]

		// 路由匹配
		matches, err := match(apiRoutes, method, path)
		if err != nil {
			return outputError(err)
		}

		// 执行 handler
		result, err := executeHandler(adapter, path, matches)
		if err != nil {
			return outputError(err)
		}

		// 输出结果
		if jsonFlag {
			return json.NewEncoder(os.Stdout).Encode(result)
		}

		// 简单 JSON 输出
		printResult(result)
		return nil
	},
}

// executeHandler 根据匹配到的路由模板执行对应的 adapter 方法
func executeHandler(adapter *internalGitea.Adapter, path string, matches MatchResult) (interface{}, error) {
	switch matches.Pattern {
	case "/repos":
		if matches.Method == "GET" {
			return adapter.ListRepos()
		}
		if matches.Method == "POST" {
			return nil, fmt.Errorf("POST /repos requires --name and --description flags")
		}
	case "/repos/:username/:repoName":
		if matches.Method == "GET" {
			return adapter.GetRepo(matches.Params["username"], matches.Params["repoName"])
		}
	case "/repos/:username/:repoName/issues":
		if matches.Method == "GET" {
			return adapter.ListRepoIssues(matches.Params["username"], matches.Params["repoName"])
		}
	case "/repos/:username/:repoName/milestones":
		if matches.Method == "GET" {
			return adapter.ListRepoMilestones(matches.Params["username"], matches.Params["repoName"])
		}
	case "/repos/:username/:repoName/milestones/:milestonePrefix/issues":
		if matches.Method == "GET" {
			return adapter.ListIssuesByMilestonePrefix(matches.Params["username"], matches.Params["repoName"], matches.Params["milestonePrefix"])
		}
	case "/version":
		if matches.Method == "GET" {
			return adapter.GetGiteaVersion()
		}
	case "/users/:username/repos":
		if matches.Method == "GET" {
			return listRepoOfUser(matches.Params["username"])
		}
	}
	return nil, fmt.Errorf("unsupported API: %s %s", matches.Method, path)
}

// printResult 简单 JSON 输出
func printResult(data interface{}) {
	if data == nil {
		fmt.Println("OK")
		return
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("%v\n", data)
		return
	}
	fmt.Println(string(jsonBytes))
}

func init() {
	rootCmd.AddCommand(apiCmd)
}

/* ==== 下列方法被 api 和 plan 命令共享使用 ==== */

// listRepoOfUser 获取指定用户下的所有仓库列表
func listRepoOfUser(username string) ([]*gitea.Repository, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}
	return adapter.ListUserRepos(username)
}

// showRepo 获取指定用户下的单个仓库
// username: Gitea 用户名（owner）
// repoName: 仓库名称
func showRepo(username, repoName string) (*gitea.Repository, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}
	return adapter.GetRepo(username, repoName)
}

