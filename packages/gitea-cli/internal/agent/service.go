package agent

import (
	"fmt"

	"cmp"
	"slices"

	"code.gitea.io/sdk/gitea"

	internalGitea "com.lisitede.backstage.gitea/internal/gitea"
	"com.lisitede.backstage.gitea/internal/plan"
)

func GetPlan(appId string, planId string) (string, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return "", fmt.Errorf("failed to create adapter: %w", err)
	}

	// Plan 创建时 Issue title 以 planId 开头，如 "PLAN-102: xxx"
	issue, err := adapter.ShowIssueOfRepoByPrefix(appId, planId, planId)
	if err != nil {
		return "", err
	}

	// 获取评论
	issueNo := fmt.Sprintf("%d", issue.Index)
	comments, err := adapter.ListCommentOfIssue(appId, planId, issueNo)
	if err != nil {
		return "", err
	}

	// 按 ID 降序排序
	slices.SortFunc(comments, func(a, b *gitea.Comment) int {
		return cmp.Compare(b.ID, a.ID)
	})

	// 返回第一个评论的 body
	if len(comments) > 0 {
		return comments[0].Body, nil
	}

	return "", nil
}

// HeadPhase 获取 Phase Metadata
func HeadPhase(appId, planId, phaseId string) (string, error) {
	phase, err := plan.ShowPhase(appId, planId, phaseId)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`# Phase Metadata
- appId: %s
- planId: %s
- phaseId: %s
- phaseName: %s
- phaseStatus: %s
- $source: "backstage-gitea agent HEAD /%s/%s/%s"
`, appId, planId, phaseId, phase.Name, phase.Status, appId, planId, phaseId), nil
}

// GetPhase 获取 Phase 详情
func GetPhase(appId, planId, phaseId string) (string, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return "", fmt.Errorf("failed to create adapter: %w", err)
	}

	// Phase 创建时 Issue title 以 phaseId 开头，如 "PHASE-001: xxx"
	issue, err := adapter.ShowIssueOfRepoByPrefix(appId, planId, phaseId)
	if err != nil {
		return "", err
	}

	// 获取评论
	issueNo := fmt.Sprintf("%d", issue.Index)
	comments, err := adapter.ListCommentOfIssue(appId, planId, issueNo)
	if err != nil {
		return "", err
	}

	// 按 ID 降序排序
	slices.SortFunc(comments, func(a, b *gitea.Comment) int {
		return cmp.Compare(b.ID, a.ID)
	})

	// 返回第一个评论的 body
	if len(comments) > 0 {
		return comments[0].Body, nil
	}

	return "", nil
}

// HeadTask 获取 Task Metadata
func HeadTask(appId, planId, phaseId, taskId string) (string, error) {
	task, err := plan.ShowTask(appId, planId, phaseId, taskId)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`# Task Metadata
- appId: %s
- planId: %s
- phaseId: %s
- taskId: %s
- taskName: %s
- taskStatus: %s
- $source: "backstage-gitea agent GET /%s/%s/%s/%s"
`, appId, planId, phaseId, taskId, task.Name, task.Status, appId, planId, phaseId, taskId), nil
}

// GetTask 获取 Task 详情
func GetTask(appId, planId, phaseId, taskId string) (string, error) {
	task, err := plan.ShowTask(appId, planId, phaseId, taskId)
	if err != nil {
		return "", err
	}

	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return "", fmt.Errorf("failed to create adapter: %w", err)
	}

	// 获取评论
	issueNo := fmt.Sprintf("%d", task.Gitea.No)
	comments, err := adapter.ListCommentOfIssue(appId, planId, issueNo)
	if err != nil {
		return "", err
	}

	// 按 ID 降序排序
	slices.SortFunc(comments, func(a, b *gitea.Comment) int {
		return cmp.Compare(b.ID, a.ID)
	})

	// 返回第一个评论的 body
	if len(comments) > 0 {
		return comments[0].Body, nil
	}

	return "", nil
}
