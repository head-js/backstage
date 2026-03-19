package gitea

import (
	"code.gitea.io/sdk/gitea"
)

// ListMilestoneOfRepo 获取指定仓库的 Milestones 列表
// username: Gitea 用户名（owner）
// repoName: 仓库名称
func ListMilestoneOfRepo(username, repoName string) ([]*gitea.Milestone, error) {
	adapter, err := NewAdapter()
	if err != nil {
		return nil, err
	}
	return adapter.ListRepoMilestones(username, repoName)
}

// ListIssueOfMilestone 获取指定 Milestone 下的 Issues 列表
// owner: 仓库所有者
// repo: 仓库名称
// milestoneId: Milestone 的数字 ID
func ListIssueOfMilestone(owner, repo, milestoneId string) ([]*gitea.Issue, error) {
	adapter, err := NewAdapter()
	if err != nil {
		return nil, err
	}

	issues, err := adapter.SearchRepoIssues(owner, repo, gitea.ListIssueOption{
		Milestones: []string{milestoneId},
	})
	return issues, err
}
