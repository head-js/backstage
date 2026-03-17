package gitea

import (
	"fmt"
	"os"
	"strings"

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

// ListRepoMilestones 列出仓库 Milestone 列表 (GET /repos/:username/:repoName/milestones)
func (a *Adapter) ListRepoMilestones(owner, repo string) ([]*gitea.Milestone, error) {
	milestones, _, err := a.client.ListRepoMilestones(owner, repo, gitea.ListMilestoneOption{})
	return milestones, err
}

// ListIssuesByMilestonePrefix 按 milestone 名称前缀查询 Issues
// 流程：ListRepoMilestones → 客户端前缀过滤 → 逐个 ListRepoIssues 合并
func (a *Adapter) ListIssuesByMilestonePrefix(owner, repo, prefix string) ([]*gitea.Issue, error) {
	milestones, _, err := a.client.ListRepoMilestones(owner, repo, gitea.ListMilestoneOption{})
	if err != nil {
		return nil, err
	}

	var result []*gitea.Issue
	for _, m := range milestones {
		if !strings.HasPrefix(m.Title, prefix) {
			continue
		}
		issues, _, err := a.client.ListRepoIssues(owner, repo, gitea.ListIssueOption{
			Milestones: []string{m.Title},
		})
		if err != nil {
			return nil, err
		}
		result = append(result, issues...)
	}
	return result, nil
}
