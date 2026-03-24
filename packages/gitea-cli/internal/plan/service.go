package plan

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

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

// GenNextTaskId 生成下一个 Task 编号
// appName: 应用名称
// planName: Plan 名称
// phaseId: Phase ID（如 PHASE-01）
// 返回 Task 编号字符串（如 "01"）
func GenNextTaskId(appName, planName, phaseId string) (string, error) {
	milestoneId, err := TranslatePhaseId2MilestoneId(appName, planName, phaseId)
	if err != nil {
		return "", err
	}

	issues, err := internalGitea.ListIssueOfMilestone(appName, planName, milestoneId)
	if err != nil {
		return "", err
	}

	// 调用 translator 把 issue.title 转成我们业务定义的 taskId
	var ids []string
	for _, issue := range issues {
		_, id, err := ExtractTaskId(issue.Title)
		if err != nil {
			fmt.Printf("Warning: invalid task naming for issue %s: %v\n", issue.Title, err)
			continue
		}
		// id 已经是数字部分
		ids = append(ids, id)
	}

	// 排序
	sort.Strings(ids)

	// 打印整个数组
	fmt.Printf("Task IDs: %v\n", ids)

	nextId, err := genNextId(ids)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("TASK-%s", nextId), nil
}

// GenNextPhaseId 生成下一个 Phase 编号
// appName: 应用名称
// planName: Plan 名称
// 返回 Phase 编号字符串（如 "PHASE-200"）
func GenNextPhaseId(appName, planName string) (string, error) {
	milestones, err := internalGitea.ListMilestoneOfRepo(appName, planName)
	if err != nil {
		return "", err
	}

	// 使用 translator 提取每个 milestone 的 phase ID
	translator := NewPlanTranslator()
	var ids []string
	for _, m := range milestones {
		_, id, err := translator.ExtractPhaseId(m.Title)
		if err != nil {
			fmt.Printf("Warning: invalid phase naming for milestone %s: %v\n", m.Title, err)
			continue
		}
		// id 已经是数字部分
		ids = append(ids, id)
	}

	// 排序
	sort.Strings(ids)

	// 打印整个数组
	fmt.Printf("Phase IDs: %v\n", ids)

	nextId, err := genNextId(ids)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("PHASE-%s", nextId), nil
}

// genNextId 私有方法
func genNextId(existIds []string) (string, error) {
	if len(existIds) == 0 {
		return "100", nil
	}

	lastNum := existIds[len(existIds)-1]
	lastInt, err := strconv.Atoi(lastNum)
	if err != nil {
		return "", fmt.Errorf("invalid task id: %s", lastNum)
	}

	next := lastInt + 1
	if next >= 900 {
		return "", fmt.Errorf("no available task id: reached 900")
	}

	return fmt.Sprintf("%03d", next), nil
}

// genReserveId 私有方法
// 规则：取最大 ID 的百位，直接跳到下一个百位（如 110 → 200，225 → 300）
func genReserveId(existIds []string) (string, error) {
	if len(existIds) == 0 {
		return "100", nil
	}

	lastNum := existIds[len(existIds)-1]
	a := lastNum[0] - '0'

	if a >= 9 {
		return "", fmt.Errorf("no available task id: reached 900")
	}

	return fmt.Sprintf("%d00", a+1), nil
}

// SyncPlanToWiki 同步 Plan 到 Gitea Wiki
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
