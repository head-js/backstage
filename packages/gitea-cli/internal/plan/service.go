package plan

import (
	"fmt"
	"os"
	"path/filepath"
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
		return nil, err
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

	slices.SortFunc(phases, func(a, b Phase) int {
		return cmp.Compare(a.Id, b.Id)
	})

	// 获取每个 Phase 的所有 Task
	for i := range phases {
		tasks, err := ListTaskOfPhase(appId, planId, phases[i].Id)
		if err != nil {
			return nil, err
		}
		phases[i].Tasks = tasks
	}

	var phasesMarkdown strings.Builder
	for _, p := range phases {
		phasesMarkdown.WriteString(fmt.Sprintf("### %s: %s\n\n", p.Id, p.Name))
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
	planId, _, _, err := ExtractPlanId(planTitle)
	if err != nil {
		return nil, err
	}

	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}

	// 初始化 Repo
	_, err = adapter.CreateRepo(appId, planId, planTitle)
	if err != nil {
		return nil, err
	}

	// 初始化 Phase Labels
	for _, status := range []Status{StatusPass, StatusFail, StatusTodo, StatusHold} {
		labelName := status.Translate2Label("PHASE")
		color := StatusColor[status]
		if _, err := adapter.CreateLabel(appId, planId, labelName, color); err != nil {
			return nil, fmt.Errorf("failed to create label %s: %w", labelName, err)
		}
	}

	// 初始化 Task Labels
	for _, status := range []Status{StatusPass, StatusFail, StatusTodo, StatusHold} {
		labelName := status.Translate2Label("TASK")
		color := StatusColor[status]
		if _, err := adapter.CreateLabel(appId, planId, labelName, color); err != nil {
			return nil, fmt.Errorf("failed to create label %s: %w", labelName, err)
		}
	}

	// 初始化 Issue for Plan
	issue, err := adapter.CreateIssue(appId, planId, planTitle, "", "")
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
		return nil, err
	}

	_, err = adapter.CreateMilestone(appId, planId, title)
	if err != nil {
		return nil, err
	}

	issue, err := adapter.CreateIssue(appId, planId, title, "", "")
	if err != nil {
		return nil, err
	}

	// 获取 Label 并添加到 Issue
	label, err := adapter.GetLabelByName(appId, planId, StatusTodo.Translate2Label("PHASE"))
	if err != nil {
		return nil, err
	}

	// 即 Issue.Number
	_, err = adapter.AddLabelToIssue(appId, planId, fmt.Sprintf("%d", issue.Index), fmt.Sprintf("%d", label.ID))
	if err != nil {
		return nil, err
	}

	// 使用 ShowPhase 统一处理
	phase, err := ShowPhase(appId, planId, nextPhaseId)
	if err != nil {
		return nil, err
	}

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
	issue, err := adapter.CreateIssue(appId, planId, title, "", milestoneId)
	if err != nil {
		return nil, err
	}

	label, err := adapter.GetLabelByName(appId, planId, StatusTodo.Translate2Label("TASK"))
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

	issues, err := adapter.SearchIssueOfRepo(appId, planId, gitea.ListIssueOption{
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

	issue, err := adapter.ShowIssueByPrefixId(appId, planId, phaseId)
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
	return nil, framework.NotFoundException("task not found: " + taskId)
}

// ListAllPlans 获取全量 Plans（参考 api.go 的 GET /repos）
func ListAllPlans() ([]Plan, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}

	repos, err := adapter.ListRepos()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repos: %w", err)
	}

	translator := NewPlanTranslator()
	plans, err := translator.TranslateRepoList2PlanList(repos)
	if err != nil {
		return nil, fmt.Errorf("failed to translate repos to plans: %w", err)
	}

	return plans, nil
}

// ListPlanOfApp 获取指定应用名称下的所有 plan repos
func ListPlanOfApp(appId string) ([]Plan, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}
	repos, err := adapter.ListRepoOfOwner(appId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repos: %w", err)
	}

	translator := NewPlanTranslator()
	plans, err := translator.TranslateRepoList2PlanList(repos)
	if err != nil {
		return nil, fmt.Errorf("failed to translate repos to plans: %w", err)
	}

	return plans, nil
}

// ListPhaseOfPlan 获取指定 Plan 下的所有 Phases
func ListPhaseOfPlan(appId, planId string) ([]Phase, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}

	// 通过搜索 PHASE-* 前缀的 Issues 获取 Phases
	issues, err := adapter.SearchIssueByPrefix(appId, planId, "PHASE-")
	if err != nil {
		return nil, fmt.Errorf("failed to search phases: %w", err)
	}

	translator := NewPlanTranslator()
	phases, err := translator.TranslateIssueList2Phase(issues)
	if err != nil {
		return nil, err
	}

	// 按 Phase.Id 升序排序
	slices.SortFunc(phases, func(a, b Phase) int {
		return cmp.Compare(a.Id, b.Id)
	})

	return phases, nil
}

// DownloadPlan 下载 Plan 数据（含所有 Phase 及 Task）
// appId: e.g. cms-mgr
// planId: e.g. PLAN-102
// 流程：
// 1. 先调用 SyncPlanToWiki 刷新 Wiki
// 2. 下载 Wiki/Home 到 {planId}/README.md
// 3. 下载 {planId}/{planId}.md（Plan 最新 Comment）
// 4. foreach Phase：{planId}/{phaseId}/{phaseId}.md
// 5. foreach Task：{planId}/{phaseId}/{taskId}.md
func DownloadPlan(appId, planId string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	fmt.Printf("\033[31mDownload plan to: %s\033[0m\n", pwd)

	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return "", err
	}

	// 0. 先调用 SyncPlanToWiki 刷新 Wiki，保证最新
	if _, err := SyncPlanToWiki(appId, planId); err != nil {
		return "", fmt.Errorf("failed to sync plan to wiki: %w", err)
	}

	// 创建 {planId} 目录
	planDir := filepath.Join(pwd, planId)
	if err := os.MkdirAll(planDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// 1. 下载 Wiki/Home 到 {planId}/README.md
	wikiContent, err := adapter.GetWikiContent(appId, planId, "Home")
	if err != nil {
		return "", fmt.Errorf("failed to get wiki: %w", err)
	}
	readmeFile := filepath.Join(planDir, "README.md")
	if err := os.WriteFile(readmeFile, []byte(wikiContent), 0644); err != nil {

		return "", fmt.Errorf("failed to write README.md: %w", err)
	}

	// 2. 获取 Plan 的 Issue
	issue, err := adapter.ShowIssueById(appId, planId, planId)
	if err != nil {
		return "", fmt.Errorf("failed to get plan issue: %w", err)
	}

	// 3. 获取 Plan Issue 的最新一条 Comment
	issueNo := fmt.Sprintf("%d", issue.Index)
	lastComment, err := adapter.ShowLastCommentOfIssue(appId, planId, issueNo)
	if err != nil {
		return "", err
	}

	// 4. 写入 {planId}/{planId}.md
	planFile := filepath.Join(planDir, fmt.Sprintf("%s.md", planId))
	if err := os.WriteFile(planFile, []byte(lastComment.Body), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// 5. 获取所有 Phases，foreach 下载
	phases, err := ListPhaseOfPlan(appId, planId)
	if err != nil {
		return "", fmt.Errorf("failed to list phases: %w", err)
	}
	for _, phase := range phases {
		// 获取 Phase 的 Issue
		phaseIssue, err := adapter.ShowIssueById(appId, planId, phase.Id)
		if err != nil {
			return "", fmt.Errorf("failed to get phase issue %s: %w", phase.Id, err)
		}
		// 获取 Phase Issue 的最新一条 Comment
		phaseIssueNo := fmt.Sprintf("%d", phaseIssue.Index)
		phaseComment, err := adapter.ShowLastCommentOfIssue(appId, planId, phaseIssueNo)
		if err != nil {
			return "", fmt.Errorf("failed to get comment of phase %s: %w", phase.Id, err)
		}
		// 创建 {planId}/{phaseId} 目录并写入 {planId}/{phaseId}/{phaseId}.md
		phaseDir := filepath.Join(planDir, phase.Id)
		if err := os.MkdirAll(phaseDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
		phaseFile := filepath.Join(phaseDir, fmt.Sprintf("%s.md", phase.Id))
		if err := os.WriteFile(phaseFile, []byte(phaseComment.Body), 0644); err != nil {
			return "", fmt.Errorf("failed to write file: %w", err)
		}

		// 6. 下载该 Phase 下的所有 Tasks
		tasks, err := ListTaskOfPhase(appId, planId, phase.Id)
		if err != nil {
			return "", fmt.Errorf("failed to list tasks of phase %s: %w", phase.Id, err)
		}
		for _, task := range tasks {
			// 获取 Task 的 Issue
			taskIssue, err := adapter.ShowIssueById(appId, planId, task.Id)
			if err != nil {
				return "", fmt.Errorf("failed to get task issue %s: %w", task.Id, err)
			}
			// 获取 Task Issue 的最新一条 Comment
			taskIssueNo := fmt.Sprintf("%d", taskIssue.Index)
			taskComment, err := adapter.ShowLastCommentOfIssue(appId, planId, taskIssueNo)
			if err != nil {
				return "", fmt.Errorf("failed to get comment of task %s: %w", task.Id, err)
			}
			// 写入 {planId}/{phaseId}/{taskId}.md
			taskFile := filepath.Join(phaseDir, fmt.Sprintf("%s.md", task.Id))
			if err := os.WriteFile(taskFile, []byte(taskComment.Body), 0644); err != nil {
				return "", fmt.Errorf("failed to write file: %w", err)
			}
		}
	}

	return fmt.Sprintf("Downloaded to: %s", planDir), nil
}
