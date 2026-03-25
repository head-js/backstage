package gitea

import (
	"code.gitea.io/sdk/gitea"
)

// ListMilestoneOfRepo 获取指定仓库的 Milestones 列表
func ListMilestoneOfRepo(owner, repoName string) ([]*gitea.Milestone, error) {
	adapter, err := NewAdapter()
	if err != nil {
		return nil, err
	}
	milestones, err := adapter.ListMilestoneOfRepo(owner, repoName)
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
