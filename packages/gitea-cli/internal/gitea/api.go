package gitea

import (
	"code.gitea.io/sdk/gitea"
)

// listMilestoneOfRepo 获取指定仓库的 Milestones 列表
// username: Gitea 用户名（owner）
// repoName: 仓库名称
func ListMilestoneOfRepo(username, repoName string) ([]*gitea.Milestone, error) {
	adapter, err := NewAdapter()
	if err != nil {
		return nil, err
	}
	return adapter.ListRepoMilestones(username, repoName)
}
