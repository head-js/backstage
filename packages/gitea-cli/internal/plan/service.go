package plan

import (
	"fmt"
	"strings"
	"time"

	"code.gitea.io/sdk/gitea"
	internalGitea "com.lisitede.backstage.gitea/internal/gitea"
)

// TranslatePhaseId2MilestoneId 根据 phaseId (如 PHASE-01) 查找对应的 Gitea Milestone numeric ID
// 流程：调用 api 获取 milestones → 遍历匹配 title 前缀（忽略大小写）
func TranslatePhaseId2MilestoneId(appName, planName, phaseId string) (string, error) {
	milestones, err := internalGitea.ListMilestoneOfRepo(appName, planName)
	if err != nil {
		return "", err
	}

	for _, m := range milestones {
		if strings.HasPrefix(strings.ToUpper(m.Title), phaseId) {
			return fmt.Sprintf("%d", m.ID), nil
		}
	}

	return "", fmt.Errorf("milestone not found for phaseId: %s", phaseId)
}

func SyncPlanToWiki(appId, planId string) (interface{}, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to create adapter: %w", err)
	}

	repo, err := adapter.GetRepo(appId, planId)
	if err != nil {
		return nil, err
	}

	translator := NewPlanTranslator()
	plan, err := translator.TranslateRepo2Plan(repo)
	if err != nil {
		return nil, err
	}

	milestones, err := internalGitea.ListMilestoneOfRepo(appId, planId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch milestones: %w", err)
	}

	phases := translator.TranslateMilestoneList2PhaseList(milestones)

	// TODO: 实现 Plan 同步到 Wiki 的逻辑
	// 2. 获取每个 Phase 的所有 Task
	// 3. 将 Plan、Phase、Task 的信息同步到 Gitea Wiki

	var phasesMarkdown strings.Builder
	for _, p := range phases {
		phasesMarkdown.WriteString(fmt.Sprintf("### %s\n\n", p.Name))
	}

	markdown := fmt.Sprintf(`# %s: %s

> updated_at: %s

## Phases

%s`, plan.Id, plan.Name, time.Now().Format("2006-01-02 15:04:05"), phasesMarkdown.String())

	// 约定保存在 Wiki/Index
	_, err = adapter.UpdateWikiOfRepo(appId, planId, "Index", markdown)
	if err != nil {
		return nil, fmt.Errorf("failed to update wiki: %w", err)
	}

	return map[string]interface{}{
		"code":    0,
		"message": "ok",
	}, nil
}

// CreatePlan
// appId: e.g. cms-mgr
// planTitle: e.g. "PLAN-102: UploadImage"
func CreatePlan(appId, planTitle string) (*gitea.Repository, error) {
	planId, err := ExtractPlanId(planTitle)
	if err != nil {
		return nil, fmt.Errorf("invalid plan name: %w", err)
	}

	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to create adapter: %w", err)
	}

	// 初始化 Repo
	repo, err := adapter.CreateRepo(appId, planId, planTitle)
	if err != nil {
		return nil, err
	}

	// 初始化 Phase Labels
	for _, status := range []Status{StatusPass, StatusFail, StatusTodo, StatusHold} {
		labelName := fmt.Sprintf("PHASE-%s", status)
		color := StatusColor[status]
		if _, err := adapter.CreateLabel(appId, planId, labelName, color); err != nil {
			return nil, fmt.Errorf("failed to create label %s: %w", labelName, err)
		}
	}

	// 初始化 Task Labels
	for _, status := range []Status{StatusPass, StatusFail, StatusTodo, StatusHold} {
		labelName := fmt.Sprintf("TASK-%s", status)
		color := StatusColor[status]
		if _, err := adapter.CreateLabel(appId, planId, labelName, color); err != nil {
			return nil, fmt.Errorf("failed to create label %s: %w", labelName, err)
		}
	}

	// 初始化 Issue for Plan
	if _, err := adapter.CreateIssue(appId, planId, planTitle, ""); err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	// 初始化 Wiki/Index
	_, err = adapter.UpdateWikiOfRepo(appId, planId, "Index", planTitle)
	if err != nil {
		return nil, fmt.Errorf("failed to create wiki: %w", err)
	}

	return repo, nil
}

// CreatePhase
// appId: e.g. cms-mgr
// planId: e.g. PLAN-102
// phaseName: e.g. "Upload-Image"
func CreatePhase(appId, planId, phaseName string) (*gitea.Milestone, error) {
	nextPhaseId, err := GenNextPhaseId(appId, planId)
	if err != nil {
		return nil, err
	}

	title := fmt.Sprintf("%s: %s", nextPhaseId, phaseName)

	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to create adapter: %w", err)
	}

	milestone, err := adapter.CreateMilestone(appId, planId, title)
	if err != nil {
		return nil, err
	}

	issue, err := adapter.CreateIssue(appId, planId, title, "")
	if err != nil {
		return nil, err
	}

	// 获取 Label 并添加到 Issue
	label, err := adapter.GetLabelByName(appId, planId, "PHASE-TODO")
	if err != nil {
		return nil, err
	}

	// 即 Issue.Number
	_, err = adapter.AddLabelToIssue(appId, planId, fmt.Sprintf("%d", issue.Index), fmt.Sprintf("%d", label.ID))
	if err != nil {
		return nil, err
	}

	return milestone, nil
}

// CreateTask 创建 Task
// appId: e.g. cms-mgr
// planId: e.g. PLAN-102
// phaseId: e.g. PHASE-200
// taskName: e.g. "Integrate @adobe/editor"
func CreateTask(appId, planId, phaseId string, taskName string) (*gitea.Issue, error) {
	milestoneId, err := TranslatePhaseId2MilestoneId(appId, planId, phaseId)
	if err != nil {
		return nil, err
	}

	nextTaskId, err := GenNextTaskId(appId, planId, phaseId)
	if err != nil {
		return nil, err
	}

	title := fmt.Sprintf("%s: %s", nextTaskId, taskName)

	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}
	issue, err := adapter.CreateIssue(appId, planId, title, milestoneId)
	if err != nil {
		return nil, err
	}

	label, err := adapter.GetLabelByName(appId, planId, "TASK-TODO")
	if err != nil {
		return nil, err
	}

	// 即 Issue.Number
	_, err = adapter.AddLabelToIssue(appId, planId, fmt.Sprintf("%d", issue.Index), fmt.Sprintf("%d", label.ID))
	if err != nil {
		return nil, err
	}

	return issue, nil
}
