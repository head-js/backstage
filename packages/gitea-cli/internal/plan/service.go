package plan

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"cmp"

	"code.gitea.io/sdk/gitea"
	"com.lisitede.backstage.gitea/framework"
	internalGitea "com.lisitede.backstage.gitea/internal/gitea"
)

// TranslatePhaseId2MilestoneId 根据 phaseId (如 PHASE-01) 查找对应的 Gitea Milestone numeric ID
// 流程：调用 api 获取 milestones → 遍历匹配 title 前缀（忽略大小写）
func TranslatePhaseId2MilestoneId(appName, planName, phaseId string) (string, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return "", err
	}
	milestones, err := adapter.ListMilestoneOfRepo(appName, planName)
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

	// 通过搜索 PHASE-* 前缀的 Issues 获取 Phases
	issues, err := adapter.SearchIssueByPrefix(appId, planId, "PHASE-")
	if err != nil {
		return nil, fmt.Errorf("failed to search phases: %w", err)
	}

	var phases []Phase
	for _, issue := range issues {
		phase, err := translator.TranslateIssue2Phase(issue)
		if err != nil {
			return nil, fmt.Errorf("failed to translate issue to phase: %w", err)
		}
		phases = append(phases, *phase)
	}

	// 获取每个 Phase 的所有 Task
	for i := range phases {
		tasks, err := ListTaskOfPhase(appId, planId, phases[i].Id)
		if err != nil {
			return nil, fmt.Errorf("failed to list tasks for phase %s: %w", phases[i].Id, err)
		}
		phases[i].Tasks = tasks
	}

	var phasesMarkdown strings.Builder
	for _, p := range phases {
		phasesMarkdown.WriteString(fmt.Sprintf("### %s\n\n", p.Name))
		// 输出该 Phase 下的所有 Tasks
		for _, t := range p.Tasks {
			checkbox := "[ ]"
			if t.Status == "PASS" {
				checkbox = "[x]"
			}
			phasesMarkdown.WriteString(fmt.Sprintf("- %s %s: %s\n\n", checkbox, t.Id, t.Name))
		}
	}

	markdown := fmt.Sprintf(`# %s: %s

> updated_at: %s

## Phases

%s`, plan.Id, plan.Name, time.Now().Format("2006-01-02 15:04:05"), phasesMarkdown.String())

	// 约定保存在 Wiki/Home
	_, err = adapter.UpdateWikiOfRepo(appId, planId, "Home", markdown)
	if err != nil {
		return nil, fmt.Errorf("failed to update wiki: %w", err)
	}

	return framework.RestOK, nil
}

// CreatePlan
// appId: e.g. cms-mgr
// planTitle: e.g. "PLAN-102: UploadImage"
func CreatePlan(appId, planTitle string) (*Plan, error) {
	planId, _, err := ExtractPlanId(planTitle)
	if err != nil {
		return nil, fmt.Errorf("invalid plan name: %w", err)
	}

	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to create adapter: %w", err)
	}

	// 初始化 Repo
	_, err = adapter.CreateRepo(appId, planId, planTitle)
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
	issue, err := adapter.CreateIssue(appId, planId, planTitle, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	// 初始化 Wiki/Home
	_, err = adapter.UpdateWikiOfRepo(appId, planId, "Home", planTitle)
	if err != nil {
		return nil, fmt.Errorf("failed to create wiki: %w", err)
	}

	translator := NewPlanTranslator()
	var newPlan *Plan
	newPlan, err = translator.TranslateIssue2Plan(issue)
	if err != nil {
		return nil, fmt.Errorf("failed to translate issue to plan: %w", err)
	}

	return newPlan, nil
}

// CreatePhase 创建 Phase
// appId: e.g. cms-mgr
// planId: e.g. PLAN-102
// phaseName: e.g. "Upload-Image"
func CreatePhase(appId, planId, phaseName string) (*Phase, error) {
	nextPhaseId, err := GenNextPhaseId(appId, planId)
	if err != nil {
		return nil, err
	}

	title := fmt.Sprintf("%s: %s", nextPhaseId, phaseName)

	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to create adapter: %w", err)
	}

	_, err = adapter.CreateMilestone(appId, planId, title)
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

	// 使用 translator 创建 Phase 对象
	translator := NewPlanTranslator()
	phase, _ := translator.TranslateIssue2Phase(issue)

	return phase, nil
}

// CreateTask 创建 Task
// appId: e.g. cms-mgr
// planId: e.g. PLAN-102
// phaseId: e.g. PHASE-200
// taskName: e.g. "Integrate @adobe/editor"
func CreateTask(appId, planId, phaseId string, taskName string) (*Task, error) {
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

	translator := NewPlanTranslator()
	task := translator.TranslateIssue2Task(issue)

	return &task, nil
}

// ListTaskOfPhase 获取指定 Phase 下的所有 Tasks
func ListTaskOfPhase(appId, planId, phaseId string) ([]Task, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}

	milestoneId, err := TranslatePhaseId2MilestoneId(appId, planId, phaseId)
	if err != nil {
		return nil, err
	}

	issues, err := adapter.SearchRepoIssues(appId, planId, gitea.ListIssueOption{
		Milestones: []string{milestoneId},
	})
	if err != nil {
		return nil, err
	}

	translator := NewPlanTranslator()
	tasks := translator.TranslateIssueList2TaskList(issues)

	// 按 Task.Id 升序排序
	slices.SortFunc(tasks, func(a, b Task) int {
		return cmp.Compare(a.Id, b.Id)
	})

	return tasks, nil
}

// ShowPhase 获取指定 Phase 详情
func ShowPhase(appId, planId, phaseId string) (*Phase, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}

	issue, err := adapter.ShowIssueById(appId, planId, phaseId)
	if err != nil {
		return nil, err
	}

	translator := NewPlanTranslator()
	phase, err := translator.TranslateIssue2Phase(issue)
	if err != nil {
		return nil, err
	}

	return phase, nil
}

// ShowTask 获取指定 Task 详情
func ShowTask(appId, planId, phaseId, taskId string) (*Task, error) {
	tasks, err := ListTaskOfPhase(appId, planId, phaseId)
	if err != nil {
		return nil, err
	}
	for _, t := range tasks {
		if t.Id == taskId {
			task := t
			return &task, nil
		}
	}
	return nil, fmt.Errorf("task not found")
}

// ListPhaseOfPlan 获取指定 Plan 下的所有 Phases
func ListPhaseOfPlan(appId, planId string) ([]Phase, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to create adapter: %w", err)
	}

	// 通过搜索 PHASE-* 前缀的 Issues 获取 Phases
	issues, err := adapter.SearchIssueByPrefix(appId, planId, "PHASE-")
	if err != nil {
		return nil, fmt.Errorf("failed to search phases: %w", err)
	}

	translator := NewPlanTranslator()
	var phases []Phase
	for _, issue := range issues {
		phase, err := translator.TranslateIssue2Phase(issue)
		if err != nil {
			return nil, fmt.Errorf("failed to translate issue to phase: %w", err)
		}
		phases = append(phases, *phase)
	}

	// 按 Phase.Id 升序排序
	slices.SortFunc(phases, func(a, b Phase) int {
		return cmp.Compare(a.Id, b.Id)
	})

	return phases, nil
}
