package plan

import (
	"fmt"
	"sort"
	"strconv"

	"code.gitea.io/sdk/gitea"

	internalGitea "com.lisitede.backstage.gitea/internal/gitea"
)

// GenNextPhaseId 生成下一个 Phase 编号
// appName: 应用名称
// planName: Plan 名称
// 返回 Phase 编号字符串（如 "PHASE-200"）
func GenNextPhaseId(appName, planName string) (string, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return "", err
	}
	milestones, err := adapter.ListMilestoneOfRepo(appName, planName)
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

	// CreatePhase 默认使用 NextTierId 模式
	nextId, err := genNextTierId(ids)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("PHASE-%s", nextId), nil
}

// GenNextTaskId 生成下一个 Task 编号
// appName: 应用名称
// planName: Plan 名称
// phaseId: Phase ID（如 PHASE-01）
// 返回 Task 编号字符串（如 "TASK-120"）
func GenNextTaskId(appName, planName, phaseId string) (string, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return "", err
	}

	milestoneId, err := TranslatePhaseId2MilestoneId(appName, planName, phaseId)
	if err != nil {
		return "", err
	}

	issues, err := adapter.SearchRepoIssues(appName, planName, gitea.ListIssueOption{
		Milestones: []string{milestoneId},
	})
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

	// 默认使用 NextBlockId 模式
	nextId, err := genNextBlockId(ids)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("TASK-%s", nextId), nil
}

// genNextId 私有方法
// 规则：连续递增（如 100 → 101，225 → 226）
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

// genNextTierId 私有方法
// 规则：取最大 ID 的百位，直接跳到下一个百位（如 110 → 200，225 → 300）
func genNextTierId(existIds []string) (string, error) {
	if len(existIds) == 0 {
		return "100", nil
	}

	lastNum := existIds[len(existIds)-1]
	a := lastNum[0] - '0'

	if a >= 9 {
		return "", fmt.Errorf("no available id: reached 900")
	}

	return fmt.Sprintf("%d00", a+1), nil
}

// genNextBlockId 私有方法
// 规则：取最大 ID 的十位，直接跳到下一个十位（如 110 → 120，225 → 230）
func genNextBlockId(existIds []string) (string, error) {
	if len(existIds) == 0 {
		return "100", nil
	}

	lastNum := existIds[len(existIds)-1]

	// 取百位
	hundred := lastNum[0] - '0'
	// 取十位
	ten := lastNum[1] - '0'

	if ten >= 9 {
		// 十位到 9 时，跳到下一个百位
		if hundred >= 9 {
			return "", fmt.Errorf("no available id: reached 900")
		}
		return fmt.Sprintf("%d00", hundred+1), nil
	}

	return fmt.Sprintf("%d%d0", hundred, ten+1), nil
}
