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
