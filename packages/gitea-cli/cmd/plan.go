package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"code.gitea.io/sdk/gitea"
	internalGitea "com.lisitede.backstage.gitea/internal/gitea"
	"com.lisitede.backstage.gitea/internal/plan"
	"github.com/spf13/cobra"
	"github.com/ucarion/urlpath"
)

// planRoutes plan 命令路由列表
var planRoutes = []Route{
	{Method: "LIST", Pattern: "/apps/:app/plans", Matcher: urlpath.New("/apps/:app/plans")},
	{Method: "GET", Pattern: "/apps/:app/plans/:planName", Matcher: urlpath.New("/apps/:app/plans/:planName")},
	{Method: "LIST", Pattern: "/apps/:app/plans/:planName/phases/:phaseId/tasks", Matcher: urlpath.New("/apps/:app/plans/:planName/phases/:phaseId/tasks")},
}

// planCmd 计划管理命令
var planCmd = &cobra.Command{
	Use:   "plan <method> <path>",
	Short: "Manage plans",
	Long: `Manage plans in Gitea.
Examples:
	backstage-gitea plan LIST /apps/mall-view-platform/plans
	backstage-gitea plan GET /apps/mall-view-platform/plans/plan-102-HttpClient-Rules`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		method := args[0]
		path := args[1]

		// 路由匹配
		matches, err := match(planRoutes, method, path)
		if err != nil {
			return outputError(err)
		}

		// 执行 handler
		var result interface{}
		switch matches.Pattern {
		case "/apps/:app/plans":
			result, err = listPlanOfApp(matches.Params["app"])
		case "/apps/:app/plans/:planName":
			result, err = showPlan(matches.Params["app"], matches.Params["planName"])
		case "/apps/:app/plans/:planName/phases/:phaseId/tasks":
			result, err = listTaskOfPhase(matches.Params["app"], matches.Params["planName"], matches.Params["phaseId"])
		default:
			return outputError(fmt.Errorf("unsupported path: %s", path))
		}
		if err != nil {
			return outputError(err)
		}

		// 输出结果
		if jsonFlag {
			return json.NewEncoder(os.Stdout).Encode(result)
		}

		// 简单 JSON 输出
		printResult(result)
		return nil
	},
}

// listPlanOfApp 获取指定应用名称下的所有 plan repos
func listPlanOfApp(appName string) (interface{}, error) {
	repos, err := listRepoOfUser(appName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repos: %w", err)
	}

	translator := plan.NewPlanTranslator()
	plans, err := translator.TranslateRepoList2PlanList(repos)
	if err != nil {
		return nil, fmt.Errorf("failed to translate repos to plans: %w", err)
	}

	return plans, nil
}

// showPlan 获取指定 planName 的单个 Plan 详情
// planName: Plan.Name 字段（完整的 Gitea Repo 名，如 "plan-102-HttpClient-Rules"）
func showPlan(appName, planName string) (*plan.Plan, error) {
	// 直接获取单个 repo
	repo, err := showRepo(appName, planName)
	if err != nil {
		return nil, err
	}

	// 获取该 repo 的 projects
	milestones, err := internalGitea.ListMilestoneOfRepo(appName, planName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch milestones: %w", err)
	}

	// 将 projects 转换为 phases（只包含符合 PHASE-xx 格式的）
	translator := plan.NewPlanTranslator()
	phases := translator.TranslateMilestoneList2PhaseList(milestones)

	// 转换为 Plan 对象
	p, err := translator.TranslateRepo2Plan(repo)
	if err != nil {
		return nil, err
	}

	// 填充 Phases 字段
	p.Phases = phases

	return p, nil
}

// listTaskOfPhase 获取指定 Phase 下的所有 Tasks
// 返回 []plan.Task 数组
func listTaskOfPhase(appName, planName, phaseId string) (interface{}, error) {
	// 调用 domain service（内部会调用 api 的 listMilestoneOfRepo）
	milestoneId, err := plan.TranslatePhaseId2MilestoneId(appName, planName, phaseId)
	if err != nil {
		return nil, err
	}

	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}

	issues, err := adapter.ListRepoIssues(appName, planName)
	if err != nil {
		return nil, err
	}

	// 过滤属于指定 milestone 的 issues
	var filteredIssues []*gitea.Issue
	for _, issue := range issues {
		if issue.Milestone != nil && fmt.Sprintf("%d", issue.Milestone.ID) == milestoneId {
			filteredIssues = append(filteredIssues, issue)
		}
	}

	translator := plan.NewPlanTranslator()
	tasks := translator.TranslateIssueList2TaskList(filteredIssues)

	return tasks, nil
}

func init() {
	rootCmd.AddCommand(planCmd)
}
