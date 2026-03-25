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

// ShowIssueOfRepoByPrefix 根据 Issue ID 前缀获取指定仓库的单个 Issue
func (a *Adapter) ShowIssueOfRepoByPrefix(owner, repoName, issueTitlePrefix string) (*gitea.Issue, error) {
	issues, err := a.SearchRepoIssues(owner, repoName, gitea.ListIssueOption{
		KeyWord: issueTitlePrefix,
	})

	if err != nil {
		return nil, err
	}

	for _, issue := range issues {
		if strings.HasPrefix(issue.Title, issueTitlePrefix) {
			return issue, nil
		}
	}

	return nil, fmt.Errorf("issue not found: %s", issueTitlePrefix)
}
