// Package plan 提供 Plan 领域模型的翻译器
// 将 Gitea API 对象（Repository、Milestone、Issue）转换为 Plan 领域模型（Plan、Phase、Task）
//
// # 可用的方法如下，严禁 Agentic Workers 自行添加方法：
//
//	NewPlanTranslator()                              创建翻译器实例
//	TranslateRepoList2PlanList(repos)                Repo 列表 → Plan 列表
//	TranslateRepo2Plan(repo)                         单个 Repo → Plan
//	TranslateMilestoneList2PhaseList(milestones)     Milestone 列表 → Phase 列表
//	TranslateMilestone2Phase(milestone)              单个 Milestone → Phase
//	ExtractPhaseId(milestoneTitle)                   从 Milestone 标题提取 Phase ID
//	TranslateIssueList2TaskList(issues)             Issue 列表 → Task 列表
//	TranslateIssue2Task(issue)                       单个 Issue → Task
//	ExtractTaskId(issueTitle)                        从 Issue 标题提取 Task ID
package plan

import (
	"fmt"
	"regexp"
	"strings"
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

// TranslateRepo2Plan 将单个 Gitea Repo 转换为 Plan
func (pt *PlanTranslator) TranslateRepo2Plan(repo *gitea.Repository) (*Plan, error) {
	id, err := extractPlanId(repo.Name)
	if err != nil {
		fmt.Printf("Warning: invalid plan naming for repo %s: %v\n", repo.Name, err)
		return nil, fmt.Errorf("invalid plan naming: %w", err)
	}

	plan := &Plan{
		Id:   id,
		Name: repo.Name, // 原样保留 Gitea Repo 名称
		Gitea: GiteaExtra{
			Type:        "REPO",
			Id:          repo.ID,
			Name:        repo.Name,
			Description: repo.Description,
			CreatedAt:   repo.Created.Format(time.RFC3339),
			UpdatedAt:   repo.Updated.Format(time.RFC3339),
		},
		Phases: []Phase{},
	}

	return plan, nil
}

// extractPlanId 从 repo 名称提取计划 ID 并大写化
// Gitea.Repo.Name -> plan-102-HttpClient-Rules
// Plan.Id         -> PLAN-102（注意这个大写化）
func extractPlanId(repoName string) (string, error) {
	pattern := `^plan-(\d+)-(.+)$`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(repoName)

	if len(matches) < 2 {
		return "", fmt.Errorf("repo name must follow format 'plan-{id}-{name}': %s", repoName)
	}

	id := matches[1]                        // 数字部分：102
	planId := "PLAN-" + strings.ToUpper(id) // 大写化：PLAN-102
	return planId, nil
}

// ExtractPhaseId 从 milestone 名称提取 Phase ID
// Gitea.Milestone.Title -> Phase-02
// Phase.Id              -> PHASE-02
// 返回: match[1]=PHASE-xxx, match[2]=xxx数字部分, err
func extractPhaseId(milestoneTitle string) (string, string, error) {
	pattern := `(?i)^phase-(\d{3})`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(milestoneTitle)

	if len(matches) < 2 {
		return "", "", fmt.Errorf("milestone title must follow format 'Phase-{id}': %s", milestoneTitle)
	}

	id := matches[1]                          // 数字部分：02
	phaseId := "PHASE-" + strings.ToUpper(id) // 大写化：PHASE-02
	return phaseId, id, nil
}

// ExtractPhaseId 公开版本，供外部调用
func (pt *PlanTranslator) ExtractPhaseId(milestoneTitle string) (string, string, error) {
	return extractPhaseId(milestoneTitle)
}

// TranslateMilestone2Phase 将单个 Gitea Milestone 转换为 Phase（只含 id 和 name）
func (pt *PlanTranslator) TranslateMilestone2Phase(m *gitea.Milestone) Phase {
	phaseId, _, err := extractPhaseId(m.Title)
	if err != nil {
		fmt.Printf("Warning: invalid phase naming for milestone %s: %v\n", m.Title, err)
		phaseId = m.Title
	}

	updatedAt := ""
	if m.Updated != nil {
		updatedAt = m.Updated.Format(time.RFC3339)
	}

	return Phase{
		Id:    phaseId,
		Name:  m.Title,
		Tasks: []Task{},
		Gitea: GiteaExtra{
			Type:        "MILESTONE",
			Id:          m.ID,
			Name:        m.Title,
			Description: m.Description,
			CreatedAt:   m.Created.Format(time.RFC3339),
			UpdatedAt:   updatedAt,
		},
	}
}

// TranslateMilestoneList2PhaseList 将 Gitea Milestones 列表转换为 Phases 列表
// 每个 Phase 只包含 id 和 title
func (pt *PlanTranslator) TranslateMilestoneList2PhaseList(milestones []*gitea.Milestone) []Phase {
	var phases []Phase

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
	taskId, _, err := ExtractTaskId(issue.Title)
	if err != nil {
		fmt.Printf("Warning: invalid task naming for issue %s: %v\n", issue.Title, err)
		taskId = fmt.Sprintf("TASK-%d", issue.Index)
	}

	updatedAt := ""
	if !issue.Updated.IsZero() {
		updatedAt = issue.Updated.Format(time.RFC3339)
	}

	return Task{
		Id:   taskId,
		Name: issue.Title,
		Gitea: GiteaExtra{
			Type:      "ISSUE",
			Id:        issue.ID,
			Name:      issue.Title,
			Body:      issue.Body,
			State:     string(issue.State),
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

// ExtractTaskId 从 issue 标题提取 Task ID
// Gitea.Issue.Title -> "TASK-101: 修复权限模块 (perm/order)"
// Task.Id            -> "TASK-101"
// 返回: match[1]=TASK-xxx, match[2]=xxx数字部分, err
func ExtractTaskId(issueTitle string) (string, string, error) {
	pattern := `(?i)^task-(\d{3})`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(issueTitle)

	if len(matches) < 3 {
		return "", "", fmt.Errorf("issue title must start with 'TASK-{id}': %s", issueTitle)
	}

	return matches[1], matches[2], nil
}
