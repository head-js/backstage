package gitea

import (
	"encoding/base64"
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

// ListRepoOfOwner 列出用户的所有仓库
func (a *Adapter) ListRepoOfOwner(owner string) ([]*gitea.Repository, error) {
	repos, _, err := a.client.ListUserRepos(owner, gitea.ListReposOptions{})
	return repos, err
}

// GetRepo 获取仓库详情 (GET /repos/:username/:repoName)
func (a *Adapter) GetRepo(owner, repo string) (*gitea.Repository, error) {
	result, _, err := a.client.GetRepo(owner, repo)
	return result, err
}

// ListIssueOfRepo 列出仓库 Issue 列表
func (a *Adapter) ListIssueOfRepo(owner, repoName string) ([]*gitea.Issue, error) {
	issues, _, err := a.client.ListRepoIssues(owner, repoName, gitea.ListIssueOption{})
	return issues, err
}

func (a *Adapter) SearchRepoIssues(owner, repoName string, searchOptions gitea.ListIssueOption) ([]*gitea.Issue, error) {
	issues, _, err := a.client.ListRepoIssues(owner, repoName, searchOptions)
	return issues, err
}

// DeleteIssueOfRepo 删除仓库指定 Issue
func (a *Adapter) DeleteIssueOfRepo(owner, repo, issueId string) (interface{}, error) {
	var issueIdInt int64
	fmt.Sscanf(issueId, "%d", &issueIdInt)

	_, err := a.client.DeleteIssue(owner, repo, issueIdInt)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"code": 0, "message": "ok"}, nil
}

// ListMilestoneOfRepo 列出仓库 Milestone 列表
func (a *Adapter) ListMilestoneOfRepo(owner, repo string) ([]*gitea.Milestone, error) {
	milestones, _, err := a.client.ListRepoMilestones(owner, repo, gitea.ListMilestoneOption{})
	return milestones, err
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

// ListLabelOfRepo 列出仓库 Label 列表 (GET /repos/:owner/:repoName/labels)
func (a *Adapter) ListLabelOfRepo(owner, repo string) ([]*gitea.Label, error) {
	labels, _, err := a.client.ListRepoLabels(owner, repo, gitea.ListLabelsOptions{})
	return labels, err
}

// CreateLabel 创建 Label (POST /repos/:owner/:repoName/labels)
func (a *Adapter) CreateLabel(owner, repo, name, color string) (*gitea.Label, error) {
	label, _, err := a.client.CreateLabel(owner, repo, gitea.CreateLabelOption{
		Name:  name,
		Color: color,
	})
	return label, err
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

// AddLabelToIssue 为 Issue 添加 Label
// issueId: Issue 的数字 ID（字符串形式）
// labelId: Label 的数字 ID（字符串形式）
func (a *Adapter) AddLabelToIssue(owner, repo, issueId, labelId string) ([]*gitea.Label, error) {
	var issueIdInt, labelIdInt int64
	fmt.Sscanf(issueId, "%d", &issueIdInt)
	fmt.Sscanf(labelId, "%d", &labelIdInt)

	labels, _, err := a.client.AddIssueLabels(owner, repo, issueIdInt, gitea.IssueLabelsOption{
		Labels: []int64{labelIdInt},
	})
	return labels, err
}

// ClearLabelFromIssue 清除 Issue 的 Label
// issueId: Issue 的数字 ID（字符串形式）
// labelId: Label 的数字 ID（字符串形式）
func (a *Adapter) ClearLabelFromIssue(owner, repo, issueId, labelId string) error {
	var issueIdInt, labelIdInt int64
	fmt.Sscanf(issueId, "%d", &issueIdInt)
	fmt.Sscanf(labelId, "%d", &labelIdInt)

	_, err := a.client.ClearIssueLabels(owner, repo, issueIdInt)
	return err
}

// GetLabelByName 根据名称获取 Label
// owner: 仓库所有者
// repo: 仓库名称
// labelName: Label 名称
func (a *Adapter) GetLabelByName(owner, repo, labelName string) (*gitea.Label, error) {
	labels, _, err := a.client.ListRepoLabels(owner, repo, gitea.ListLabelsOptions{})
	if err != nil {
		return nil, err
	}

	for _, label := range labels {
		if label.Name == labelName {
			return label, nil
		}
	}

	return nil, fmt.Errorf("label not found: %s", labelName)
}

// ListWikiOfRepo 列出仓库 Wiki 列表
func (a *Adapter) ListWikiOfRepo(owner, repo string) ([]*gitea.WikiPageMetaData, error) {
	wikis, _, err := a.client.ListWikiPages(owner, repo, gitea.ListWikiPagesOptions{})
	return wikis, err
}

// GetWikiOfRepo 获取仓库指定 Wiki 页面
func (a *Adapter) GetWikiOfRepo(owner, repo, wikiName string) (*gitea.WikiPage, error) {
	wiki, _, err := a.client.GetWikiPage(owner, repo, wikiName)
	return wiki, err
}

// UpdateWikiOfRepo 更新仓库指定 Wiki 页面
func (a *Adapter) UpdateWikiOfRepo(owner, repo, wikiName, content string) (*gitea.WikiPage, error) {
	encodedContent := base64.StdEncoding.EncodeToString([]byte(content))
	wiki, _, err := a.client.EditWikiPage(owner, repo, wikiName, gitea.CreateWikiPageOptions{
		Title:         wikiName,
		ContentBase64: encodedContent,
	})
	return wiki, err
}
