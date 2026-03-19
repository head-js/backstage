package gitea

import (
	"fmt"
	"os"

	"code.gitea.io/sdk/gitea"
)

// Adapter Gitea API 适配器
type Adapter struct {
	client *gitea.Client
}

// NewAdapter 创建新的 Gitea 适配器
func NewAdapter() (*Adapter, error) {
	url := os.Getenv("BACKSTAGE_GITEA_URL")
	token := os.Getenv("BACKSTAGE_GITEA_TOKEN")

	if url == "" {
		return nil, fmt.Errorf("BACKSTAGE_GITEA_URL is required")
	}

	client, err := gitea.NewClient(url, gitea.SetToken(token))
	if err != nil {
		return nil, err
	}

	return &Adapter{client: client}, nil
}

// GetGiteaVersion 获取 Gitea 版本
func (a *Adapter) GetGiteaVersion() (string, error) {
	version, _, err := a.client.ServerVersion()
	return version, err
}

// ListRepos 列出所有仓库 (GET /repos)
func (a *Adapter) ListRepos() ([]*gitea.Repository, error) {
	repos, _, err := a.client.SearchRepos(gitea.SearchRepoOptions{})
	return repos, err
}

// ListUserRepos 列出用户的所有仓库 (GET /users/{username}/repos)
func (a *Adapter) ListUserRepos(username string) ([]*gitea.Repository, error) {
	repos, _, err := a.client.ListUserRepos(username, gitea.ListReposOptions{})
	return repos, err
}

// GetRepo 获取仓库详情 (GET /repos/:username/:repoName)
func (a *Adapter) GetRepo(owner, repo string) (*gitea.Repository, error) {
	result, _, err := a.client.GetRepo(owner, repo)
	return result, err
}

// ListRepoIssues 列出仓库 Issue 列表 (GET /repos/:username/:repoName/issues)
func (a *Adapter) ListRepoIssues(owner, repo string) ([]*gitea.Issue, error) {
	issues, _, err := a.client.ListRepoIssues(owner, repo, gitea.ListIssueOption{})
	return issues, err
}

func (a *Adapter) SearchRepoIssues(owner, repo string, searchOptions gitea.ListIssueOption) ([]*gitea.Issue, error) {
	issues, _, err := a.client.ListRepoIssues(owner, repo, searchOptions)
	return issues, err
}

// ListRepoMilestones 列出仓库 Milestone 列表 (GET /repos/:username/:repoName/milestones)
func (a *Adapter) ListRepoMilestones(owner, repo string) ([]*gitea.Milestone, error) {
	milestones, _, err := a.client.ListRepoMilestones(owner, repo, gitea.ListMilestoneOption{})
	return milestones, err
}

// CreateRepo 创建仓库
// 以当前持有 Token 的 owner 为准
func (a *Adapter) CreateRepo(name string) (*gitea.Repository, error) {
	repo, _, err := a.client.CreateRepo(gitea.CreateRepoOption{
		Name:        name,
		Description: "",
		Private:     false,
	})
	return repo, err
}

// TransferRepo 转移仓库
func (a *Adapter) TransferRepo(oldOwner, repoName, newOwner string) (*gitea.Repository, error) {
	repo, _, err := a.client.TransferRepo(oldOwner, repoName, gitea.TransferRepoOption{
		NewOwner: newOwner,
	})
	return repo, err
}

// CreateMilestone 创建里程碑
func (a *Adapter) CreateMilestone(owner, repo, title string) (*gitea.Milestone, error) {
	milestone, _, err := a.client.CreateMilestone(owner, repo, gitea.CreateMilestoneOption{
		Title: title,
	})
	return milestone, err
}

// CreateIssue 创建 Issue
// milestoneId: Milestone 的数字 ID（可选，传入空字符串表示不关联）
func (a *Adapter) CreateIssue(owner, repo, title string, milestoneId string) (*gitea.Issue, error) {
	opts := gitea.CreateIssueOption{
		Title: title,
	}
	if milestoneId != "" {
		var id int64
		fmt.Sscanf(milestoneId, "%d", &id)
		if id > 0 {
			opts.Milestone = id
		}
	}
	issue, _, err := a.client.CreateIssue(owner, repo, opts)
	return issue, err
}
