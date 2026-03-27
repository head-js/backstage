package gitea

import (
	"fmt"
	"strings"

	"code.gitea.io/sdk/gitea"
)

// extra.go 提供 Gitea SDK 的功能扩展
// 通过组合 adapter 的原子接口，构建项目所需的复合功能
// 保持纯工具层属性，不承载业务逻辑

// CreateRepo 创建仓库
// owner: 目标 owner（仓库创建后转移到此用户）
// repoName: 仓库名称
// description: 仓库描述
func (a *Adapter) CreateRepo(owner string, repoName string, description string) (*gitea.Repository, error) {
	// 先用系统默认 owner 创建仓库
	_, _, err := a.client.CreateRepo(gitea.CreateRepoOption{
		Name:        repoName,
		Description: description,
		Private:     false,
	})
	if err != nil {
		return nil, err
	}

	// 立刻把仓库转移给目标 owner
	return a.TransferRepo("backstage", repoName, owner)
}

// ShowIssueById 根据 Issue ID 前缀获取指定仓库的单个 Issue
func (a *Adapter) ShowIssueById(owner, repoName, issueId string) (*gitea.Issue, error) {
	issues, err := a.SearchIssueByPrefix(owner, repoName, issueId)
	if err != nil {
		return nil, err
	}

	if len(issues) == 0 {
		return nil, fmt.Errorf("issue not found: %s", issueId)
	}

	return issues[0], nil
}

// SearchIssueByPrefix 根据 Issue Title 前缀搜索指定仓库的所有 Issues
func (a *Adapter) SearchIssueByPrefix(owner, repoName, prefix string) ([]*gitea.Issue, error) {
	issues, err := a.SearchRepoIssues(owner, repoName, gitea.ListIssueOption{
		KeyWord: prefix,
	})

	if err != nil {
		return nil, err
	}

	var matchedIssues []*gitea.Issue
	for _, issue := range issues {
		if strings.HasPrefix(issue.Title, prefix) {
			matchedIssues = append(matchedIssues, issue)
		}
	}

	return matchedIssues, nil
}
