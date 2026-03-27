// Package plan 提供 Plan 领域模型的翻译器
//
// # 可用的方法如下，严禁 Agentic Workers 自行添加方法：
//
//	ExtractPlanId(planTitle)        从 Plan 标题提取 ID 和名称
//	ExtractPhaseId(milestoneTitle)  从 Milestone 标题提取 ID 和名称
//	ExtractTaskId(issueTitle)       从 Issue 标题提取 ID 和名称
package plan

import (
	"regexp"

	"com.lisitede.backstage.gitea/framework"
)

// ExtractPlanId 从 repo 名称提取 Plan ID 和名称
// Gitea.Repository.Name -> "PLAN-102: My Plan"
// Plan.Id               -> "PLAN-102"
// Plan.Name             -> "My Plan"
// 返回: planId, planName, error
func ExtractPlanId(planTitle string) (string, string, error) {
	pattern := `(?i)^PLAN-(\d{3}):\s*(.+)$`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(planTitle)

	if len(matches) < 3 {
		return "", "", framework.InvalidFormatException("Plan.Title must follow format 'PLAN-{id}: {name}': " + planTitle)
	}

	planId := "PLAN-" + matches[1] // 完整 ID：PLAN-102
	planName := matches[2]         // 名称部分：My Plan

	return planId, planName, nil
}

// ExtractPhaseId 从 milestone 名称提取 Phase ID 和名称
// Gitea.Milestone.Title -> "PHASE-200: 设计阶段"
// Phase.Id              -> "PHASE-200"
// Phase.Name            -> "设计阶段"
// numId                 -> "200"
// 返回: phaseId, phaseName, numId, error
func ExtractPhaseId(milestoneTitle string) (string, string, string, error) {
	pattern := `^PHASE-(\d{3}):\s*(.+)$`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(milestoneTitle)

	if len(matches) < 3 {
		return "", "", "", framework.InvalidFormatException("milestone title must follow format 'PHASE-{id}: {name}': " + milestoneTitle)
	}

	phaseId := "PHASE-" + matches[1] // 完整 ID：PHASE-200
	phaseName := matches[2]          // 名称部分：设计阶段
	numId := matches[1]              // 纯数字 ID：200
	return phaseId, phaseName, numId, nil
}

// ExtractTaskId 从 issue 标题提取 Task ID 和名称
// Gitea.Issue.Title -> "TASK-101: 修复权限模块 (perm/order)"
// Task.Id           -> "TASK-101"
// Task.Name         -> "修复权限模块 (perm/order)"
// numId             -> "101"
// 返回: taskId, taskName, numId, error
func ExtractTaskId(issueTitle string) (string, string, string, error) {
	pattern := `^TASK-(\d{3}): (.+)$`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(issueTitle)

	if len(matches) < 3 {
		return "", "", "", framework.InvalidFormatException("issue title must follow format 'TASK-{id}: {name}': " + issueTitle)
	}

	taskId := "TASK-" + matches[1] // 完整 ID：TASK-101
	taskName := matches[2]         // 名称部分：修复权限模块 (perm/order)
	numId := matches[1]           // 纯数字 ID：101
	return taskId, taskName, numId, nil
}
