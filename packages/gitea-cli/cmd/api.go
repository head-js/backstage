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
	{Method: "POST", Pattern: "/repos/:owner/:repoName/transfer", Matcher: urlpath.New("/repos/:owner/:repoName/transfer")},
}

// createRepoFlags 用于 POST /repos 的 flags
var createRepoFlags struct {
	name string
}

// transferRepoFlags 用于 POST /repos/:owner/:repoName/transfer 的 flags
var transferRepoFlags struct {
	newOwner string
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
	case "/version":
		if matches.Method == "GET" {
			return adapter.GetGiteaVersion()
		}
	case "/users/:username/repos":
		if matches.Method == "GET" {
			return listRepoOfUser(matches.Params["username"])
		}
	case "/repos/:owner/:repoName/transfer":
		if matches.Method == "POST" {
			if transferRepoFlags.newOwner == "" {
				return nil, fmt.Errorf("POST /repos/:owner/:repoName/transfer requires --new-owner flag")
			}
			return adapter.TransferRepo(matches.Params["owner"], matches.Params["repoName"], transferRepoFlags.newOwner)
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

	// 为 POST /repos 添加 flags
	apiCmd.Flags().StringVar(&createRepoFlags.name, "name", "", "Repository name for POST /repos")

	// 为 POST /repos/:owner/:repoName/transfer 添加 flags
	apiCmd.Flags().StringVar(&transferRepoFlags.newOwner, "new-owner", "", "New owner for repository transfer")
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

// createRepo 创建仓库
// owner: 目标 owner（仓库创建后转移到此用户）
// repoName: 仓库名称
func createRepo(owner string, repoName string) (*gitea.Repository, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}

	// 先用系统默认 owner 创建仓库
	_, err = adapter.CreateRepo(repoName)
	if err != nil {
		return nil, err
	}

	// 立刻把仓库转移给目标 owner
	return transferRepo("backstage", repoName, owner)
}

// transferRepo 转移仓库
// oldOwner: 原仓库所有者
// repoName: 仓库名称
// newOwner: 新仓库所有者
func transferRepo(oldOwner, repoName, newOwner string) (*gitea.Repository, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}
	return adapter.TransferRepo(oldOwner, repoName, newOwner)
}

// createMilestone 创建里程碑
// owner: 仓库所有者
// repo: 仓库名称
// title: 里程碑标题
func createMilestone(owner, repo, title string) (*gitea.Milestone, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}
	return adapter.CreateMilestone(owner, repo, title)
}

// createIssue 创建 Issue
// owner: 仓库所有者
// repo: 仓库名称
// title: Issue 标题
// milestoneId: Milestone 的数字 ID（用于关联 Milestone）
func createIssue(owner, repo, title, milestoneId string) (*gitea.Issue, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}
	return adapter.CreateIssue(owner, repo, title, milestoneId)
}
