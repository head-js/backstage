# Design

## internal/gitea

围绕 Gitea SDK 构建两层代理，各司其职。

### adapter.go - 命名适配

Gitea SDK 接口的直接包装，用项目规范的命名方式重新暴露，屏蔽原库命名差异。

### extra.go - 功能组合

将 adapter 提供的原子接口组合成复合功能。纯工具层，不承载业务逻辑。

### 使用示例

底层使用 adapter 的规范接口：

```go
import internalGitea "com.lisitede.backstage.gitea/internal/gitea"

adapter, err := internalGitea.NewAdapter()
issues, err := adapter.ListIssueOfRepo(owner, repoName)
```

上层调用 extra 的复合功能：

```go
issue, err := internalGitea.ShowIssueById(owner, repoName, issueId)
```

## internal/plan

围绕 Plan、Phase、Task 的业务逻辑构建，使用 gitea adapter 和 extra 提供的能力。

### service.go - 业务服务

轻业务逻辑（与 Gitea 原能力相近）：暂时放在 cmd/plan.go 快速迭代。

重业务逻辑（明显区别于 Gitea 原能力，高度耦合）：写在 service.go。

## cmd/agent

为 Agent 提供的能力接口，针对 Agent 的使用习惯进行了 UX 优化。
