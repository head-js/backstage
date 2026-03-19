package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"code.gitea.io/sdk/gitea"
	internalGitea "com.lisitede.backstage.gitea/internal/gitea"
	"github.com/spf13/cobra"
)

// apiRouter API 命令路由管理器
var apiRouter Router

// createRepoFlags 用于 POST /repos 的 flags
var createRepoFlags struct {
	name string
}

// transferRepoFlags 用于 POST /repos/:owner/:repoName/transfer 的 flags
var transferRepoFlags struct {
	newOwner string
}

// createLabelFlags 用于 POST /repos/:owner/:repoName/labels 的 flags
var createLabelFlags struct {
	color string
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
		method := strings.ToUpper(args[0])
		path := args[1]

		// 构建 args 映射
		argsMap := map[string]string{}
		if createRepoFlags.name != "" {
			argsMap["name"] = createRepoFlags.name
		}
		if transferRepoFlags.newOwner != "" {
			argsMap["newOwner"] = transferRepoFlags.newOwner
		}
		if createLabelFlags.color != "" {
			argsMap["color"] = createLabelFlags.color
		}

		// 调用路由
		result, err := apiRouter.Invoke(method, path, argsMap)
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
	// 注册 api 路由
	apiRouter.Verb("GET", "/repos", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListRepos()
	})

	apiRouter.Verb("POST", "/repos", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		if args["name"] == "" {
			return nil, fmt.Errorf("POST /repos requires --name flag")
		}
		// args["name"] 是 repo 名称，需要创建仓库
		// 这里需要确定 owner，暂时使用默认值
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.CreateRepo(args["name"])
	})

	apiRouter.Verb("GET", "/repos/:username/:repoName", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.GetRepo(params["username"], params["repoName"])
	})

	apiRouter.Verb("GET", "/repos/:owner/:repoName/labels", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListLabelOfRepo(params["owner"], params["repoName"])
	})

	apiRouter.Verb("POST", "/repos/:owner/:repoName/labels", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.CreateLabel(params["owner"], params["repoName"], args["name"], args["color"])
	})

	apiRouter.Verb("GET", "/repos/:username/:repoName/milestones", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return internalGitea.ListMilestoneOfRepo(params["username"], params["repoName"])
	})

	apiRouter.Verb("GET", "/version", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.GetGiteaVersion()
	})

	apiRouter.Verb("GET", "/users/:username/repos", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return listRepoOfUser(params["username"])
	})

	apiRouter.Verb("POST", "/repos/:owner/:repoName/transfer", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		if args["newOwner"] == "" {
			return nil, fmt.Errorf("POST /repos/:owner/:repoName/transfer requires --new-owner flag")
		}
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.TransferRepo(params["owner"], params["repoName"], args["newOwner"])
	})

	rootCmd.AddCommand(apiCmd)

	// 为 POST /repos 添加 flags
	apiCmd.Flags().StringVar(&createRepoFlags.name, "name", "", "Repository name for POST /repos")

	// 为 POST /repos/:owner/:repoName/transfer 添加 flags
	apiCmd.Flags().StringVar(&transferRepoFlags.newOwner, "new-owner", "", "New owner for repository transfer")

	// 为 POST /repos/:owner/:repoName/labels 添加 flags
	apiCmd.Flags().StringVar(&createLabelFlags.color, "color", "", "Label color (hex) for POST /repos/:owner/:repoName/labels")
}

/* ==== 下列方法被 api 和 plan 命令共享使用 ==== */

// listRepoOfUser 获取指定用户下的所有仓库列表
func listRepoOfUser(owner string) ([]*gitea.Repository, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}
	return adapter.ListUserRepos(owner)
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
