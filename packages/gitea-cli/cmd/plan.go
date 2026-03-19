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
	{Method: "LIST", Pattern: "/apps/:appName/plans", Matcher: urlpath.New("/apps/:appName/plans")},
	{Method: "POST", Pattern: "/apps/:appName/plans", Matcher: urlpath.New("/apps/:appName/plans")},
	{Method: "GET", Pattern: "/apps/:appName/plans/:planName", Matcher: urlpath.New("/apps/:appName/plans/:planName")},
	{Method: "LIST", Pattern: "/apps/:appName/plans/:planName/phases/:phase/tasks", Matcher: urlpath.New("/apps/:appName/plans/:planName/phases/:phase/tasks")},
	{Method: "POST", Pattern: "/apps/:appName/plans/:planName/phases", Matcher: urlpath.New("/apps/:appName/plans/:planName/phases")},
	{Method: "POST", Pattern: "/apps/:appName/plans/:planName/phases/:phase/tasks", Matcher: urlpath.New("/apps/:appName/plans/:planName/phases/:phase/tasks")},
}

// planCreateFlags 用于创建 Plan 的 flags
var planCreateFlags struct {
	name string
}

// planCmd 计划管理命令
var planCmd = &cobra.Command{
	Use:   "plan <method> <path>",
	Short: "Manage plans",
	Long: `Manage plans in Gitea.
Examples:
	backstage-gitea plan LIST /apps/mall-view-platform/plans
	backstage-gitea plan POST /apps/mall-view-platform/plans --name "plan-102-HttpClient-Rules"
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
		case "/apps/:appName/plans":
			switch matches.Method {
			case "POST":
				result, err = createRepo(matches.Params["appName"], planCreateFlags.name)
			case "LIST":
				result, err = listPlanOfApp(matches.Params["appName"])
			default:
				return outputError(fmt.Errorf("unsupported method: %s", matches.Method))
			}
		case "/apps/:appName/plans/:planName":
			result, err = showPlan(matches.Params["appName"], matches.Params["planName"])
		case "/apps/:appName/plans/:planName/phases":
			if matches.Method == "POST" {
				if planCreateFlags.name == "" {
					return outputError(fmt.Errorf("POST /apps/:appName/plans/:planName/phases requires --name flag"))
				}
				result, err = createPhase(matches.Params["appName"], matches.Params["planName"], planCreateFlags.name)
			} else {
				return outputError(fmt.Errorf("unsupported method: %s", matches.Method))
			}
		case "/apps/:appName/plans/:planName/phases/:phase/tasks":
			switch matches.Method {
			case "LIST":
				result, err = listTaskOfPhase(matches.Params["appName"], matches.Params["planName"], matches.Params["phase"])
			case "POST":
				if planCreateFlags.name == "" {
					return outputError(fmt.Errorf("POST /apps/:appName/plans/:planName/phases/:phase/tasks requires --name flag"))
				}
				result, err = createTask(matches.Params["appName"], matches.Params["planName"], planCreateFlags.name, matches.Params["phase"])
			default:
				return outputError(fmt.Errorf("unsupported method: %s", matches.Method))
			}
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
	repo, err := showRepo(appName, planName)
	if err != nil {
		return nil, err
	}

	milestones, err := internalGitea.ListMilestoneOfRepo(appName, planName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch milestones: %w", err)
	}

	translator := plan.NewPlanTranslator()
	phases := translator.TranslateMilestoneList2PhaseList(milestones)

	p, err := translator.TranslateRepo2Plan(repo)
	if err != nil {
		return nil, err
	}

	p.Phases = phases

	return p, nil
}

// listTaskOfPhase 获取指定 Phase 下的所有 Tasks
func listTaskOfPhase(appName, planName, phaseId string) (interface{}, error) {
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

// createPhase 创建 Phase（Milestone）
// appName: 应用名称
// planName: Plan 名称
// name: Phase 名称
func createPhase(appName, planName, name string) (*gitea.Milestone, error) {
	nextPhaseId, err := plan.GenNextPhaseId(appName, planName)
	if err != nil {
		return nil, err
	}

	title := fmt.Sprintf("%s: %s", nextPhaseId, name)

	return createMilestone(appName, planName, title)
}

// createTask 创建 Task
// appName: 应用名称
// planName: Plan 名称
// name: Task 标题
// phase: Phase ID
func createTask(appName, planName, name, phase string) (*gitea.Issue, error) {
	milestoneId, err := plan.TranslatePhaseId2MilestoneId(appName, planName, phase)
	if err != nil {
		return nil, fmt.Errorf("failed to translate phaseId to milestoneId: %w", err)
	}

	nextTaskId, err := plan.GenNextTaskId(appName, planName, phase)
	if err != nil {
		return nil, err
	}

	title := fmt.Sprintf("%s: %s", nextTaskId, name)

	return createIssue(appName, planName, title, milestoneId)
}

func init() {
	rootCmd.AddCommand(planCmd)

	planCmd.Flags().StringVar(&planCreateFlags.name, "name", "", "Name")
}
