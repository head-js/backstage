package plan

import (
	"fmt"
	"sort"
	"strings"

	giteaApi "com.lisitede.backstage.gitea/internal/gitea"
)

// TranslatePhaseId2MilestoneId 根据 phaseId (如 PHASE-01) 查找对应的 Gitea Milestone numeric ID
// 流程：调用 api 获取 milestones → 遍历匹配 title 前缀（忽略大小写）
func TranslatePhaseId2MilestoneId(appName, planName, phaseId string) (string, error) {
	milestones, err := giteaApi.ListMilestoneOfRepo(appName, planName)
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

	issues, err := giteaApi.ListIssueOfMilestone(appName, planName, milestoneId)
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
	milestones, err := giteaApi.ListMilestoneOfRepo(appName, planName)
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

// genNextId 私有方法，算法重用
// 规则：取最大 ID 的百位，直接跳到下一个百位（如 110 → 200，225 → 300）
func genNextId(existIds []string) (string, error) {
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
