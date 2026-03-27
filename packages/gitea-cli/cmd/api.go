package cmd

import (
	"strings"

	internalGitea "com.lisitede.backstage.gitea/internal/gitea"
	"github.com/spf13/cobra"
)

// apiRouter API 命令路由管理器
var apiRouter Router

// apiCmd 通用 API 命令
var apiCmd = &cobra.Command{
	Use:   "api <method> <path>",
	Short: "Gitea Http API wrapper",
	Long: `Gitea Http API wrapper by RESTful path.
Examples:
	backstage-gitea api GET /repos
	backstage-gitea api GET /repos/:owner/:repoName
	backstage-gitea api GET /repos/:owner/:repoName/labels
	backstage-gitea api GET /repos/:owner/:repoName/milestones
	backstage-gitea api GET /repos/:owner/:repoName/issues
	backstage-gitea api GET /repos/:owner/:repoName/issues/:issueId
	backstage-gitea api GET /version
	backstage-gitea api GET /users/:username/repos
	backstage-gitea api POST /repos/:owner/:repoName/transfer-to/:newOwner`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		method := strings.ToUpper(args[0])
		path := args[1]

		// 构建 args 映射
		argsMap := map[string]string{}

		// 调用路由
		result, err := apiRouter.Invoke(method, path, argsMap)
		if err != nil {
			return outputError(err)
		}

		// 输出结果
		printResult(result)
		return nil
	},
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

	apiRouter.Verb("GET", "/repos/:owner/:repoName", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.GetRepo(params["owner"], params["repoName"])
	})

	apiRouter.Verb("GET", "/repos/:owner/:repoName/labels", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListLabelOfRepo(params["owner"], params["repoName"])
	})

	apiRouter.Verb("GET", "/repos/:owner/:repoName/milestones", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListMilestoneOfRepo(params["owner"], params["repoName"])
	})

	apiRouter.Verb("GET", "/repos/:owner/:repoName/issues", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListIssueOfRepo(params["owner"], params["repoName"])
	})

	apiRouter.Verb("GET", "/repos/:owner/:repoName/issues/:issueNo", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.GetIssueOfRepo(params["owner"], params["repoName"], params["issueNo"])
	})

	apiRouter.Verb("DELETE", "/repos/:owner/:repoName/issues/:issueNo", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.DeleteIssueOfRepo(params["owner"], params["repoName"], params["issueNo"])
	})

	apiRouter.Verb("GET", "/repos/:owner/:repoName/issues/:issueNo/comments", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListCommentOfIssue(params["owner"], params["repoName"], params["issueNo"])
	})

	apiRouter.Verb("GET", "/repos/:owner/:repoName/wikis", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListWikiOfRepo(params["owner"], params["repoName"])
	})

	apiRouter.Verb("GET", "/repos/:owner/:repoName/wikis/:wikiName", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.GetWikiOfRepo(params["owner"], params["repoName"], params["wikiName"])
	})

	apiRouter.Verb("PUT", "/repos/:owner/:repoName/wikis/:wikiName", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.UpdateWikiOfRepo(params["owner"], params["repoName"], params["wikiName"], "Hello Current Timestamp 2")
	})

	apiRouter.Verb("GET", "/version", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.GetGiteaVersion()
	})

	apiRouter.Verb("GET", "/users/:username/repos", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.ListRepoOfOwner(params["username"])
	})

	apiRouter.Verb("POST", "/repos/:owner/:repoName/transfer-to/:newOwner", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		adapter, err := internalGitea.NewAdapter()
		if err != nil {
			return nil, err
		}
		return adapter.TransferRepo(params["owner"], params["repoName"], params["newOwner"])
	})

	rootCmd.AddCommand(apiCmd)
}
