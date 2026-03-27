package agent

import (
	"fmt"
	"os"
	"path/filepath"

	"cmp"
	"slices"

	"code.gitea.io/sdk/gitea"

	"com.lisitede.backstage.gitea/framework"
	internalGitea "com.lisitede.backstage.gitea/internal/gitea"
	"com.lisitede.backstage.gitea/internal/plan"
)

func GetPlan(appId string, planId string) (string, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return "", fmt.Errorf("failed to create adapter: %w", err)
	}

	// Plan 创建时 Issue title 以 planId 开头，如 "PLAN-102: xxx"
	issue, err := adapter.ShowIssueById(appId, planId, planId)
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

// getCurrentPhase 获取当前 Phase
// 筛选出状态为 TODO/UNKNOWN 的 Phase，取第一个
func getCurrentPhaseId(appId, planId string) (string, error) {
	// plan.ListPhaseOfPlan 已按 Phase.Id 升序排序
	phases, err := plan.ListPhaseOfPlan(appId, planId)
	if err != nil {
		return "", err
	}

	if len(phases) == 0 {
		return "No phases found", nil
	}

	// 筛选出状态为 TODO/UNKNOWN 的 Phase，取第一个
	var currentPhaseId string
	for _, p := range phases {
		if p.Status == "TODO" || p.Status == "UNKNOWN" {
			currentPhaseId = p.Id
			break
		}
	}

	// 如果没有符合条件的 Phase，返回空
	if currentPhaseId == "" {
		return "No pending phases found", nil
	}

	return currentPhaseId, nil
}

func HeadCurrentPhase(appId, planId string) (string, error) {
	currentPhaseId, err := getCurrentPhaseId(appId, planId)
	if err != nil {
		return "", err
	}

	return HeadPhase(appId, planId, currentPhaseId)
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

func GetCurrentPhase(appId, planId string) (string, error) {
	currentPhaseId, err := getCurrentPhaseId(appId, planId)
	if err != nil {
		return "", err
	}

	return GetPhase(appId, planId, currentPhaseId)
}

// GetPhase 获取 Phase 详情
func GetPhase(appId, planId, phaseId string) (string, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return "", fmt.Errorf("failed to create adapter: %w", err)
	}

	// Phase 创建时 Issue title 以 phaseId 开头，如 "PHASE-001: xxx"
	issue, err := adapter.ShowIssueById(appId, planId, phaseId)
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

	// 如果没有评论，翻译 Issue 为 Phase 并返回标题
	translator := plan.NewPlanTranslator()
	phase, err := translator.TranslateIssue2Phase(issue)
	if err != nil {
		return "", fmt.Errorf("failed to translate issue to phase: %w", err)
	}

	return "## " + phase.Title, nil
}

// getCurrentTaskId 获取当前 Task
// 筛选出状态为 TODO/UNKNOWN 的 Task
func getCurrentTaskId(appId, planId, phaseId string) (string, error) {
	tasks, err := plan.ListTaskOfPhase(appId, planId, phaseId)
	if err != nil {
		return "", err
	}

	if len(tasks) == 0 {
		return "No tasks found", nil
	}

	// 筛选出状态为 TODO/UNKNOWN 的 Task，取第一个
	var currentTaskId string
	for _, t := range tasks {
		if t.Status == "TODO" || t.Status == "UNKNOWN" {
			currentTaskId = t.Id
			break
		}
	}

	// 如果没有符合条件的 Task，返回空
	if currentTaskId == "" {
		return "No pending tasks found", nil
	}

	return currentTaskId, nil
}

func HeadCurrentTask(appId, planId, phaseId string) (string, error) {
	currentTaskId, err := getCurrentTaskId(appId, planId, phaseId)
	if err != nil {
		return "", err
	}

	return HeadTask(appId, planId, phaseId, currentTaskId)
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

func GetCurrentTask(appId, planId, phaseId string) (string, error) {
	currentTaskId, err := getCurrentTaskId(appId, planId, phaseId)
	if err != nil {
		return "", err
	}

	return GetTask(appId, planId, phaseId, currentTaskId)
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

// UpdateTask 更新任务状态
// appId: 应用 ID，如 "cms-mgr"
// planId: 计划 ID，如 "PLAN-102"
// phaseId: 阶段 ID，如 "PHASE-200"
// taskId: 任务 ID，如 "TASK-101"
// status: 任务状态，如 "PASS" | "FAIL" | "TODO" | "HOLD"
// context: context 文件路径（参数是文件位置，需要读取文件内容）
func UpdateTask(appId, planId, phaseId, taskId, status, context string) (map[string]interface{}, error) {
	// 校验状态是否合法
	taskStatus := plan.Status(status)
	if !taskStatus.IsValid() {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	// 获取任务详情
	task, err := plan.ShowTask(appId, planId, phaseId, taskId)
	if err != nil {
		return nil, err
	}

	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to create adapter: %w", err)
	}

	issueNo := fmt.Sprintf("%d", task.Gitea.No)

	// 状态转换为 Label 名称
	labelName := taskStatus.Translate2Label("TASK")

	// 获取目标 Label
	label, err := adapter.GetLabelByName(appId, planId, labelName)
	if err != nil {
		return nil, fmt.Errorf("failed to get label %s: %w", labelName, err)
	}

	// 替换 Issue 的标签为新状态标签（使用 Labels 标识状态，不修改 Issue state）
	_, err = adapter.ReplaceIssueLabels(appId, planId, issueNo, []int64{label.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to replace labels: %w", err)
	}

	// 如果提供了 context 文件，读取并作为评论添加
	if context != "" {
		// 获取当前工作目录并拼接相对路径
		pwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
		contextPath := filepath.Join(pwd, context)
		// fmt.Printf("[DEBUG] context file path: %s\n", contextPath)

		content, err := os.ReadFile(contextPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read context file: %w", err)
		}

		// 添加评论
		_, err = adapter.CreateComment(appId, planId, issueNo, string(content))
		if err != nil {
			return nil, fmt.Errorf("failed to add comment: %w", err)
		}
	}

	return framework.RestOK, nil
}
