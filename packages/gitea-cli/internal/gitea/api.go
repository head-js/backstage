package gitea

import (
	"code.gitea.io/sdk/gitea"
)

// CreateRepo 创建仓库
// owner: 目标 owner（仓库创建后转移到此用户）
// repoName: 仓库名称
func CreateRepo(owner string, repoName string) (*gitea.Repository, error) {
	adapter, err := NewAdapter()
	if err != nil {
		return nil, err
	}

	// 先用系统默认 owner 创建仓库
	_, err = adapter.CreateRepo(repoName)
	if err != nil {
		return nil, err
	}

	// 立刻把仓库转移给目标 owner
	return adapter.TransferRepo("backstage", repoName, owner)
}

// TransferRepo 转移仓库
// oldOwner: 原仓库所有者
// repoName: 仓库名称
// newOwner: 新仓库所有者
func TransferRepo(oldOwner, repoName, newOwner string) (*gitea.Repository, error) {
	adapter, err := NewAdapter()
	if err != nil {
		return nil, err
	}
	return adapter.TransferRepo(oldOwner, repoName, newOwner)
}

// CreateLabel 在指定仓库创建 Label
func CreateLabel(owner, repoName, name, color string) (*gitea.Label, error) {
	adapter, err := NewAdapter()
	if err != nil {
		return nil, err
	}
	return adapter.CreateLabel(owner, repoName, name, color)
}

// ListMilestoneOfRepo 获取指定仓库的 Milestones 列表
func ListMilestoneOfRepo(owner, repoName string) ([]*gitea.Milestone, error) {
	adapter, err := NewAdapter()
	if err != nil {
		return nil, err
	}
	milestones, err := adapter.ListRepoMilestones(owner, repoName)
	if err != nil {
		return nil, err
	}
	if milestones == nil {
		return []*gitea.Milestone{}, nil
	}
	return milestones, nil
}

// ListIssueOfMilestone 获取指定 Milestone 下的 Issues 列表
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
