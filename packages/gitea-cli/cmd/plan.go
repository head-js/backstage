package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"cmp"

	"code.gitea.io/sdk/gitea"
	internalGitea "com.lisitede.backstage.gitea/internal/gitea"
	"com.lisitede.backstage.gitea/internal/plan"
	"github.com/spf13/cobra"
)

// planRouter plan 命令路由管理器
var planRouter Router

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
	backstage-gitea plan LIST /cms/plans
	backstage-gitea plan POST /cms/plans --name "plan-102-UploadImage"
	backstage-gitea plan GET  /cms/plan-102-UploadImage
	backstage-gitea plan LIST /cms/plan-102-UploadImage/phases
	backstage-gitea plan POST /cms/plan-102-UploadImage/phases --name "Design"
	backstage-gitea plan LIST /cms/plan-102-UploadImage/PHASE-100/tasks
	backstage-gitea plan POST /cms/plan-102-UploadImage/PHASE-100/tasks --name "Design-Database"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		method := args[0]
		path := args[1]

		// 构建 args 映射
		argsMap := map[string]string{}
		if planCreateFlags.name != "" {
			argsMap["name"] = planCreateFlags.name
		}

		// 调用路由
		result, err := planRouter.Invoke(method, path, argsMap)
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
func showPlan(appName, planName string) (interface{}, error) {
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

// listPhaseOfPlan 获取指定 Plan 下的所有 Phases
func listPhaseOfPlan(appName, planName string) ([]plan.Phase, error) {
	milestones, err := internalGitea.ListMilestoneOfRepo(appName, planName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch milestones: %w", err)
	}

	translator := plan.NewPlanTranslator()
	phases := translator.TranslateMilestoneList2PhaseList(milestones)

	// 按 Phase.Id 升序排序
	slices.SortFunc(phases, func(a, b plan.Phase) int {
		return cmp.Compare(a.Id, b.Id)
	})

	return phases, nil
}

// listTaskOfPhase 获取指定 Phase 下的所有 Tasks
func listTaskOfPhase(appName, planName, phaseId string) ([]plan.Task, error) {
	milestoneId, err := plan.TranslatePhaseId2MilestoneId(appName, planName, phaseId)
	if err != nil {
		return nil, err
	}

	issues, err := internalGitea.ListIssueOfMilestone(appName, planName, milestoneId)
	if err != nil {
		return nil, err
	}

	translator := plan.NewPlanTranslator()
	tasks := translator.TranslateIssueList2TaskList(issues)

	// 按 Task.Id 升序排序
	slices.SortFunc(tasks, func(a, b plan.Task) int {
		return cmp.Compare(a.Id, b.Id)
	})

	return tasks, nil
}

// createPlan 创建 Plan（Gitea Repo），并初始化三个默认 Label
// appName: 应用名称
// name: Plan 名称
func createPlan(appName, name string) (interface{}, error) {
	repo, err := internalGitea.CreateRepo(appName, name)
	if err != nil {
		return nil, err
	}

	labels := []struct {
		name  string
		color string
	}{
		{"TASK-SUCCESS", "#009800"},
		{"TASK-FAIL", "#e11d21"},
		{"TASK-TODO", "#fbca04"},
	}
	for _, l := range labels {
		if _, err := internalGitea.CreateLabel(appName, name, l.name, l.color); err != nil {
			return nil, fmt.Errorf("failed to create label %s: %w", l.name, err)
		}
	}

	return repo, nil
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

func showTask(appName, planName, phaseId, taskId string) (*plan.Task, error) {
	tasks, err := listTaskOfPhase(appName, planName, phaseId)
	if err != nil {
		return nil, err
	}
	for _, t := range tasks {
		if t.Id == taskId {
			task := t
			return &task, nil
		}
	}
	return nil, fmt.Errorf("task not found")
}

func init() {
	// 注册 plan 路由
	planRouter.Verb("LIST", "/:appName/plans", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return listPlanOfApp(params["appName"])
	})

	planRouter.Verb("POST", "/:appName/plans", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		if args["name"] == "" {
			return nil, fmt.Errorf("POST /:appName/plans requires --name flag")
		}
		return createPlan(params["appName"], args["name"])
	})

	planRouter.Verb("GET", "/:appName/:planName", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return showPlan(params["appName"], params["planName"])
	})

	planRouter.Verb("LIST", "/:appName/:planName/phases", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return listPhaseOfPlan(params["appName"], params["planName"])
	})

	planRouter.Verb("POST", "/:appName/:planName/phases", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		if args["name"] == "" {
			return nil, fmt.Errorf("POST /:appName/:planName/phases requires --name flag")
		}
		return createPhase(params["appName"], params["planName"], args["name"])
	})

	planRouter.Verb("LIST", "/:appName/:planName/:phase/tasks", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return listTaskOfPhase(params["appName"], params["planName"], params["phase"])
	})

	planRouter.Verb("POST", "/:appName/:planName/:phase/tasks", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		if args["name"] == "" {
			return nil, fmt.Errorf("POST /:appName/:planName/:phase/tasks requires --name flag")
		}
		return createTask(params["appName"], params["planName"], args["name"], params["phase"])
	})

	planRouter.Verb("GET", "/:appName/:planName/:phaseId/:taskId", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return showTask(params["appName"], params["planName"], params["phaseId"], params["taskId"])
	})

	// Backstage Plan - World Model
	planRouter.Verb("MARKDOWN", "/:appName/:planName/current-phase", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return markdownCurrentPhase(params["appName"], params["planName"])
	})

	planRouter.Verb("MARKDOWN", "/:appName/:planName/current-task", func(method, pattern, pathname string, params, args map[string]string) (interface{}, error) {
		return markdownCurrentTask(params["appName"], params["planName"])
	})

	rootCmd.AddCommand(planCmd)

	planCmd.Flags().StringVar(&planCreateFlags.name, "name", "", "Name")
}

// markdownCurrentPhase 获取当前 Phase 的 Markdown 展示
// 筛选出状态为 TODO/UNKNOWN 的 Phase，取第一个
func markdownCurrentPhase(appName, planName string) (interface{}, error) {
	// listPhaseOfPlan 已按 Phase.Id 升序排序
	phases, err := listPhaseOfPlan(appName, planName)
	if err != nil {
		return nil, err
	}

	if len(phases) == 0 {
		return "No phases found", nil
	}

	// 筛选出状态为 TODO/UNKNOWN 的 Phase，取第一个
	var currentPhase *plan.Phase
	for _, p := range phases {
		if p.Status == "TODO" || p.Status == "UNKNOWN" {
			currentPhase = &p
			break
		}
	}

	// 如果没有符合条件的 Phase，返回空
	if currentPhase == nil {
		return "No pending phases found", nil
	}

	// 获取该 Phase 下的所有 Tasks
	tasks, err := listTaskOfPhase(appName, planName, currentPhase.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}

	currentPhase.Tasks = tasks

	// 转换为 Markdown
	return plan.TranslatePhase2Markdown(*currentPhase), nil
}

// markdownCurrentTask 获取当前 Task 的 Markdown 展示
// 取第一个 Phase 的第一个 Task（目前没有状态管理，先取第一个）
func markdownCurrentTask(appName, planName string) (interface{}, error) {
	// 获取 Phases（已按 Phase.Id 升序排序）
	phases, err := listPhaseOfPlan(appName, planName)
	if err != nil {
		return nil, err
	}

	if len(phases) == 0 {
		return "No phases found", nil
	}

	// 取第一个 Phase
	currentPhase := &phases[0]

	// 获取该 Phase 下的所有 Tasks
	tasks, err := listTaskOfPhase(appName, planName, currentPhase.Id)
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
