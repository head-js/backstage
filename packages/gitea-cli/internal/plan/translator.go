// Package plan 提供 Plan 领域模型的翻译器
//
// # 可用的方法如下，严禁 Agentic Workers 自行添加方法：
//
//	NewPlanTranslator()                              创建翻译器实例
//	TranslateRepoList2PlanList(repos)                Repo 列表 → Plan 列表
//	TranslateRepo2Plan(repo)                         单个 Repo → Plan
//	TranslateIssue2Plan(issue)                       单个 Issue → Plan
//	TranslateMilestoneList2PhaseList(milestones)     Milestone 列表 → Phase 列表
//	TranslateMilestone2Phase(milestone)              单个 Milestone → Phase
//	TranslateIssue2Phase(issue)                      单个 Issue → Phase
//	TranslateIssueList2Phase(issues)                 Issue 列表 → Phase 列表
//	TranslateIssueList2TaskList(issues)             Issue 列表 → Task 列表
//	TranslateIssue2Task(issue)                       单个 Issue → Task
package plan

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"code.gitea.io/sdk/gitea"
)

// PlanTranslator 用于将 Gitea 对象转换为 Plan 领域模型
type PlanTranslator struct{}

// NewPlanTranslator 创建 PlanTranslator 实例
func NewPlanTranslator() *PlanTranslator {
	return &PlanTranslator{}
}

// TranslateRepoList2PlanList 将 Gitea repos 列表转换为 Plan 列表
// 只转换符合命名规范的 repos（plan-xxx 格式）
func (pt *PlanTranslator) TranslateRepoList2PlanList(repos []*gitea.Repository) ([]Plan, error) {
	var plans []Plan

	for _, repo := range repos {
		if repo == nil {
			continue
		}

		plan, err := pt.TranslateRepo2Plan(repo)
		if err != nil {
			continue
		}

		plans = append(plans, *plan)
	}

	return plans, nil
}

// TranslateIssue2Plan 将单个 Gitea Issue 转换为 Plan
// Issue 代表 Plan 中的顶层 Issue（如 "PLAN-102: UploadImage"）
func (pt *PlanTranslator) TranslateIssue2Plan(issue *gitea.Issue) (*Plan, error) {
	planId, planName, err := ExtractPlanId(issue.Title)
	if err != nil {
		return nil, fmt.Errorf("invalid plan naming: %w", err)
	}

	updatedAt := ""
	if !issue.Updated.IsZero() {
		updatedAt = issue.Updated.Format(time.RFC3339)
	}

	plan := &Plan{
		Id:    planId,
		Name:  planName,
		Title: issue.Title,
		Gitea: GiteaExtra{
			Type:      "ISSUE",
			Id:        issue.ID,
			No:        issue.Index,
			Title:     issue.Title,
			Body:      issue.Body,
			State:     string(issue.State),
			CreatedAt: issue.Created.Format(time.RFC3339),
			UpdatedAt: updatedAt,
		},
		Phases: []Phase{},
	}

	return plan, nil
}

// TranslateRepo2Plan 将单个 Gitea Repo 转换为 Plan
func (pt *PlanTranslator) TranslateRepo2Plan(repo *gitea.Repository) (*Plan, error) {
	planId, planName, err := ExtractPlanId(repo.Description)
	if err != nil {
		return nil, err
	}

	plan := &Plan{
		Id:    planId,
		Name:  planName,
		Title: repo.Description,
		Gitea: GiteaExtra{
			Type:        "REPO",
			Id:          repo.ID,
			Description: repo.Description,
			CreatedAt:   repo.Created.Format(time.RFC3339),
			UpdatedAt:   repo.Updated.Format(time.RFC3339),
		},
		Phases: []Phase{},
	}

	return plan, nil
}

// TranslateIssue2Phase 将单个 Gitea Issue 转换为 Phase
func (pt *PlanTranslator) TranslateIssue2Phase(issue *gitea.Issue) (*Phase, error) {
	phaseId, phaseName, _, err := ExtractPhaseId(issue.Title)
	if err != nil {
		return nil, err
	}

	updatedAt := ""
	if !issue.Updated.IsZero() {
		updatedAt = issue.Updated.Format(time.RFC3339)
	}

	// 提取 Label 名称列表，并根据 Label 设置 Status
	var labelNames []string
	status := "UNKNOWN"
	for _, label := range issue.Labels {
		if label == nil {
			continue
		}
		labelNames = append(labelNames, label.Name)
		// 根据 Label 设置 Status
		switch label.Name {
		case StatusTodo.Translate2Label("PHASE"):
			status = "TODO"
		case StatusHold.Translate2Label("PHASE"):
			status = "HOLD"
		case StatusPass.Translate2Label("PHASE"):
			status = "PASS"
		case StatusFail.Translate2Label("PHASE"):
			status = "FAIL"
		}
	}

	return &Phase{
		Id:     phaseId,
		Name:   phaseName,
		Title:  issue.Title,
		Status: status,
		Tasks:  []Task{},
		Gitea: GiteaExtra{
			Type:      "ISSUE",
			Id:        issue.ID,
			No:        issue.Index,
			Title:     issue.Title,
			State:     string(issue.State),
			Labels:    labelNames,
			CreatedAt: issue.Created.Format(time.RFC3339),
			UpdatedAt: updatedAt,
		},
	}, nil
}

// TranslateIssueList2Phase 将 Gitea Issues 列表转换为 Phase 列表
func (pt *PlanTranslator) TranslateIssueList2Phase(issues []*gitea.Issue) ([]Phase, error) {
	var phases []Phase

	for _, issue := range issues {
		if issue == nil {
			continue
		}

		phase, err := pt.TranslateIssue2Phase(issue)
		if err != nil {
			return nil, fmt.Errorf("failed to translate issue to phase: %w", err)
		}
		phases = append(phases, *phase)
	}

	return phases, nil
}

// TranslateMilestone2Phase 将单个 Gitea Milestone 转换为 Phase（只含 id 和 name）
func (pt *PlanTranslator) TranslateMilestone2Phase(m *gitea.Milestone) Phase {
	phaseId, phaseName, _, err := ExtractPhaseId(m.Title)
	if err != nil {
		fmt.Printf("Warning: invalid phase naming for milestone %s: %v\n", m.Title, err)
		phaseId = m.Title
	}

	updatedAt := ""
	if m.Updated != nil {
		updatedAt = m.Updated.Format(time.RFC3339)
	}

	return Phase{
		Id:     phaseId,
		Name:   phaseName,
		Title:  m.Title,
		Status: "UNKNOWN", // Phase 的 status 由同名 Issue 控制，和 Milestone 无关
		Tasks:  []Task{},
		Gitea: GiteaExtra{
			Type:        "MILESTONE",
			Id:          m.ID,
			Title:       m.Title,
			State:       string(m.State),
			Description: m.Description,
			CreatedAt:   m.Created.Format(time.RFC3339),
			UpdatedAt:   updatedAt,
		},
	}
}

// TranslateMilestoneList2PhaseList 将 Gitea Milestones 列表转换为 Phases 列表
// 每个 Phase 只包含 id 和 title
func (pt *PlanTranslator) TranslateMilestoneList2PhaseList(milestones []*gitea.Milestone) []Phase {
	if milestones == nil {
		return []Phase{}
	}

	phases := []Phase{}

	for _, m := range milestones {
		if m == nil {
			continue
		}

		phase := pt.TranslateMilestone2Phase(m)
		phases = append(phases, phase)
	}

	return phases
}

// TranslateIssue2Task 将单个 Gitea Issue 转换为 Task
func (pt *PlanTranslator) TranslateIssue2Task(issue *gitea.Issue) Task {
	taskId, taskName, _, err := ExtractTaskId(issue.Title)
	if err != nil {
		fmt.Printf("Warning: invalid task naming for issue %s: %v\n", issue.Title, err)
		taskId = fmt.Sprintf("TASK-%d", issue.Index)
		taskName = issue.Title
	}

	updatedAt := ""
	if !issue.Updated.IsZero() {
		updatedAt = issue.Updated.Format(time.RFC3339)
	}

	// 提取 Label 名称列表，并根据 Label 设置 Status
	var labelNames []string
	status := "UNKNOWN"
	for _, label := range issue.Labels {
		if label == nil {
			continue
		}
		labelNames = append(labelNames, label.Name)
		// 根据 Label 设置 Status
		switch label.Name {
		case StatusTodo.Translate2Label("TASK"):
			status = "TODO"
		case StatusHold.Translate2Label("TASK"):
			status = "HOLD"
		case StatusPass.Translate2Label("TASK"):
			status = "PASS"
		case StatusFail.Translate2Label("TASK"):
			status = "FAIL"
		}
	}

	return Task{
		Id:      taskId,
		Name:    taskName,
		Status:  status,
		Context: issue.Body,
		Gitea: GiteaExtra{
			Type:      "ISSUE",
			Id:        issue.ID,
			No:        issue.Index,
			Title:     issue.Title,
			Body:      issue.Body,
			State:     string(issue.State),
			Labels:    labelNames,
			CreatedAt: issue.Created.Format(time.RFC3339),
			UpdatedAt: updatedAt,
		},
	}
}

// TranslateIssueList2TaskList 将 Gitea Issues 列表转换为 Tasks 列表
func (pt *PlanTranslator) TranslateIssueList2TaskList(issues []*gitea.Issue) []Task {
	var tasks []Task

	for _, issue := range issues {
		if issue == nil {
			continue
		}

		task := pt.TranslateIssue2Task(issue)
		tasks = append(tasks, task)
	}

	return tasks
}

// renderMarkdown 渲染 Markdown 模板（私有）
func renderMarkdown(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("markdown").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// TranslatePlan2Markdown 将 Plan 转换为 Markdown 格式
func TranslatePlan2Markdown(plan Plan) (markdown string) {
	const tmplStr = `# Plan: {{.Name}}

{{range .Phases}}
## {{.Name}}

{{end}}`

	result, err := renderMarkdown(tmplStr, plan)
	if err != nil {
		return ""
	}

	return result
}

// TranslatePhase2Markdown 将 Phase 转换为 Markdown 格式
func TranslatePhase2Markdown(phase Phase) (markdown string) {
	const tmplStr = `# Current Phase - {{.Title}}

## Tasks Breakdown

{{range .Tasks}}
- [ ] {{.Name}}

{{end}}`

	result, err := renderMarkdown(tmplStr, phase)
	if err != nil {
		return ""
	}

	return result
}

// TranslateTask2Markdown 将 Task 转换为 Markdown 格式
func TranslateTask2Markdown(task Task) (markdown string) {
	const tmplStr = `# Current Task

## Metadata

- id: {{.Id}}
- name: {{.Name}}
- status: {{.Status}}

## Context

{{.Context}}
`

	result, err := renderMarkdown(tmplStr, task)
	if err != nil {
		return ""
	}

	return result
}
