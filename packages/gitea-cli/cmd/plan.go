package cmd

import (
	"fmt"
	"strings"

	internalGitea "com.lisitede.backstage.gitea/internal/gitea"
	"com.lisitede.backstage.gitea/internal/plan"
	"github.com/spf13/cobra"
)

// planRouter plan 命令路由管理器
var planRouter Router

// planCreateFlags 用于创建 Plan 的 flags
var planCreateFlags struct {
	title string
	name  string
}

// planCmd 计划管理命令
var planCmd = &cobra.Command{
	Use:   "plan <method> <path>",
	Short: "A CLI tool to manage Plan / Phase / Task.",
	Long: `A CLI tool to manage Plan / Phase / Task by RESTful-style path.
Examples:
	backstage-gitea plan LIST /cms-mgr/plans
	backstage-gitea plan GET  /cms-mgr/PLAN-102
	backstage-gitea plan LIST /cms-mgr/PLAN-102/phases
	backstage-gitea plan LIST /cms-mgr/PLAN-102/PHASE-200/tasks

Create Plan Example:
	backstage-gitea plan POST /cms-mgr/plans --title "PLAN-102: Upload Image"

Create Phase Example:
	backstage-gitea plan POST /cms-mgr/PLAN-102/phases --name "Design & Develop"

Create Task Example:
	backstage-gitea plan POST /cms-mgr/PLAN-102/PHASE-200/tasks --name "Design Database"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		method := strings.ToUpper(args[0])
		path := args[1]

		// 构建 args 映射
		argsMap := map[string]string{}
		if planCreateFlags.name != "" {
			argsMap["name"] = planCreateFlags.name
		}
		if planCreateFlags.title != "" {
			argsMap["title"] = planCreateFlags.title
		}

		// 调用路由
		result, err := planRouter.Invoke(method, path, argsMap)
		if err != nil {
			return outputError(err)
		}

		// 输出结果
		printResult(result)
		return nil
	},
}

func init() {
	// 注册 plan 路由
	planRouter.Verb("LIST", "/:appId/plans", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return listPlanOfApp(params["appId"])
	})

	planRouter.Verb("POST", "/:appId/plans", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		if args["title"] == "" {
			return nil, fmt.Errorf("POST /:appId/plans requires --title flag")
		}
		return plan.CreatePlan(params["appId"], args["title"])
	})

	planRouter.Verb("GET", "/:appId/:planId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return showPlan(params["appId"], params["planId"])
	})

	// 同步 Plan 到 Gitea Wiki
	planRouter.Verb("PATCH", "/:appId/:planId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return plan.SyncPlanToWiki(params["appId"], params["planId"])
	})

	planRouter.Verb("LIST", "/:appId/:planId/phases", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return plan.ListPhaseOfPlan(params["appId"], params["planId"])
	})

	planRouter.Verb("POST", "/:appId/:planId/phases", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		if args["name"] == "" {
			return nil, fmt.Errorf("POST /:appId/:planId/phases requires --name flag")
		}
		return plan.CreatePhase(params["appId"], params["planId"], args["name"])
	})

	planRouter.Verb("LIST", "/:appId/:planId/:phase/tasks", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return plan.ListTaskOfPhase(params["appId"], params["planId"], params["phase"])
	})

	planRouter.Verb("POST", "/:appId/:planId/:phase/tasks", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		if args["name"] == "" {
			return nil, fmt.Errorf("POST /:appId/:planId/:phase/tasks requires --name flag")
		}
		return plan.CreateTask(params["appId"], params["planId"], params["phase"], args["name"])
	})

	planRouter.Verb("GET", "/:appId/:planId/:phaseId/:taskId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return plan.ShowTask(params["appId"], params["planId"], params["phaseId"], params["taskId"])
	})

	planRouter.Verb("MARKDOWN", "/:appId/:planId/current-task", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return markdownCurrentTask(params["appId"], params["planId"])
	})

	planCmd.Flags().SortFlags = false
	planCmd.Flags().StringVar(&planCreateFlags.name, "name", "", "when create / update, provide the name")
	planCmd.Flags().StringVar(&planCreateFlags.title, "title", "", "when create / update, provide the title; i.e. title = id + name")
	rootCmd.AddCommand(planCmd)
}

// listPlanOfApp 获取指定应用名称下的所有 plan repos
func listPlanOfApp(appId string) ([]plan.Plan, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}
	repos, err := adapter.ListRepoOfOwner(appId)
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

// showPlan 获取指定 planId 的单个 Plan 详情
// planId: Plan.Name 字段（完整的 Gitea Repo 名，如 "plan-102-HttpClient-Rules"）
func showPlan(appId, planId string) (interface{}, error) {
	adapter, err := internalGitea.NewAdapter()
	if err != nil {
		return nil, err
	}

	repo, err := adapter.GetRepo(appId, planId)
	if err != nil {
		return nil, err
	}

	milestones, err := adapter.ListMilestoneOfRepo(appId, planId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch milestones: %w", err)
	}

	translator := plan.NewPlanTranslator()
	phases := translator.TranslateMilestoneList2PhaseList(milestones)

	// 该接口不查询 Tasks，如果返回空数组会造成 Agent 误解
	for i := range phases {
		phases[i].Tasks = nil
	}

	p, err := translator.TranslateRepo2Plan(repo)
	if err != nil {
		return nil, err
	}

	p.Phases = phases

	return p, nil
}

// markdownCurrentTask 获取当前 Task 的 Markdown 展示
// 取第一个 Phase 的第一个 Task（目前没有状态管理，先取第一个）
func markdownCurrentTask(appId, planId string) (interface{}, error) {
	// 获取 Phases（已按 Phase.Id 升序排序）
	phases, err := plan.ListPhaseOfPlan(appId, planId)
	if err != nil {
		return nil, err
	}

	if len(phases) == 0 {
		return "No phases found", nil
	}

	// 取第一个 Phase
	currentPhase := &phases[0]

	// 获取该 Phase 下的所有 Tasks
	tasks, err := plan.ListTaskOfPhase(appId, planId, currentPhase.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}

	if len(tasks) == 0 {
		return "No tasks found", nil
	}

	// 筛选出状态为 TODO/FAIL/UNKNOWN 的 Task，取第一个
	var currentTask *plan.Task
	for _, t := range tasks {
		if t.Status == "TODO" || t.Status == "FAIL" || t.Status == "UNKNOWN" {
			currentTask = &t
			break
		}
	}

	// 如果没有符合条件的 Task，返回空
	if currentTask == nil {
		return "No pending tasks found", nil
	}

	// 转换为 Markdown
	return plan.TranslateTask2Markdown(*currentTask), nil
}
