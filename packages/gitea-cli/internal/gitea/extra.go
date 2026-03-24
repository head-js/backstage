package gitea

import (
	"fmt"
	"strings"

	"code.gitea.io/sdk/gitea"
)

// extra.go 提供 Gitea SDK 的功能扩展
// 通过组合 adapter 的原子接口，构建项目所需的复合功能
// 保持纯工具层属性，不承载业务逻辑

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
