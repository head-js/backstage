package scrum

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"com.lisitede.backstage.vikunja/framework"
	internalVikunja "com.lisitede.backstage.vikunja/internal/vikunja"
)

// bucketSpec 描述一个目标 kanban bucket。
// RenameFrom 非空表示把已有同名 bucket 改名；为空表示插入新 bucket。
type bucketSpec struct {
	Name       string
	Position   float64
	RenameFrom string
}

// kanbanBucketSpec scrum 业务的 kanban bucket 规范（按 position 升序）。
var kanbanBucketSpec = []bucketSpec{
	{Name: "TO VISION", Position: 50},
	{Name: "IN VISION", Position: 75},
	{Name: "TO DO", Position: 100, RenameFrom: "To-Do"},
	{Name: "IN DO", Position: 200, RenameFrom: "Doing"},
	{Name: "IN TEST", Position: 260},
	{Name: "IN PRD", Position: 270},
	{Name: "DONE", Position: 300, RenameFrom: "Done"},
}

type labelSpec struct {
	Name  string
	Color string
}

var defaultLabelSpec = []labelSpec{
	{Name: "创造", Color: "22C55E"},
	{Name: "迭代", Color: "3B82F6"},
	{Name: "维护", Color: "F59E0B"},
}

// CreateWorkspace 创建 scrum 工作区：先创建默认项目，再按规范调整 kanban buckets。
func CreateWorkspace(name string) (any, error) {
	adapter, err := internalVikunja.NewAdapter()
	if err != nil {
		return nil, err
	}

	created, err := adapter.CreateProject(name)
	if err != nil {
		return nil, err
	}

	projectID, kanbanViewID, err := extractKanbanView(created)
	if err != nil {
		return nil, err
	}
	if err := configureKanbanBuckets(adapter, projectID, kanbanViewID); err != nil {
		return nil, err
	}

	return framework.RestOK, nil
}

// Initialize configures the built-in Inbox and ensures the default labels and saved filters exist.
func Initialize() (any, error) {
	adapter, err := internalVikunja.NewAdapter()
	if err != nil {
		return nil, err
	}

	projects, err := adapter.ListProjects()
	if err != nil {
		return nil, err
	}
	var resp internalVikunja.VikunjaResp
	if err := decode(projects, &resp); err != nil {
		return nil, err
	}

	for _, item := range resp.Items {
		var project internalVikunja.VikunjaProject
		if err := json.Unmarshal(item, &project); err != nil {
			return nil, err
		}
		if project.Title != "Inbox" {
			continue
		}

		projectID := strconv.FormatInt(project.Id, 10)
		views, err := adapter.ListProjectViews(projectID)
		if err != nil {
			return nil, err
		}
		kanbanViewID, err := findKanbanViewID(views)
		if err != nil {
			return nil, err
		}
		if err := configureKanbanBuckets(adapter, projectID, kanbanViewID); err != nil {
			return nil, err
		}
		if _, err := adapter.UpdateProjectTitle(projectID, "Backlog"); err != nil {
			return nil, err
		}
		break
	}
	labelIDs, err := ensureDefaultLabels(adapter)
	if err != nil {
		return nil, err
	}
	if err := ensureDefaultSavedFilters(adapter, resp.Items, labelIDs); err != nil {
		return nil, err
	}
	if err := configureDefaultSavedFilterBuckets(adapter); err != nil {
		return nil, err
	}

	return framework.RestOK, nil
}

func ensureDefaultLabels(adapter *internalVikunja.Adapter) (map[string]int64, error) {
	labelIDs := make(map[string]int64, len(defaultLabelSpec))
	for _, spec := range defaultLabelSpec {
		labels, err := adapter.SearchLabels(spec.Name)
		if err != nil {
			return nil, err
		}

		var resp internalVikunja.VikunjaResp
		if err := decode(labels, &resp); err != nil {
			return nil, err
		}

		for _, item := range resp.Items {
			var label internalVikunja.VikunjaLabel
			if err := json.Unmarshal(item, &label); err != nil {
				return nil, err
			}
			if label.Title == spec.Name {
				if existingID := labelIDs[spec.Name]; existingID != 0 && existingID != label.Id {
					return nil, fmt.Errorf("multiple labels named %q found", spec.Name)
				}
				labelIDs[spec.Name] = label.Id
			}
		}
		if labelIDs[spec.Name] != 0 {
			continue
		}

		created, err := adapter.CreateLabel(spec.Name, spec.Color)
		if err != nil {
			return nil, err
		}
		var label internalVikunja.VikunjaLabel
		if err := decode(created, &label); err != nil {
			return nil, err
		}
		if label.Id <= 0 {
			return nil, fmt.Errorf("created label %q has invalid id %d", spec.Name, label.Id)
		}
		labelIDs[spec.Name] = label.Id
	}
	return labelIDs, nil
}

func ensureDefaultSavedFilters(adapter *internalVikunja.Adapter, projects []json.RawMessage, labelIDs map[string]int64) error {
	existing := make(map[string]bool, len(defaultLabelSpec))
	for _, item := range projects {
		var project internalVikunja.VikunjaProject
		if err := json.Unmarshal(item, &project); err != nil {
			return err
		}
		if project.Id < -1 {
			existing[project.Title] = true
		}
	}

	for _, spec := range defaultLabelSpec {
		if existing[spec.Name] {
			continue
		}
		labelID := labelIDs[spec.Name]
		if labelID <= 0 {
			return fmt.Errorf("default label %q has no valid id", spec.Name)
		}
		filter := "labels = " + strconv.FormatInt(labelID, 10)
		if _, err := adapter.CreateSavedFilter(spec.Name, filter); err != nil {
			return err
		}
		existing[spec.Name] = true
	}
	return nil
}

func configureDefaultSavedFilterBuckets(adapter *internalVikunja.Adapter) error {
	projects, err := adapter.ListProjects()
	if err != nil {
		return err
	}
	var resp internalVikunja.VikunjaResp
	if err := decode(projects, &resp); err != nil {
		return err
	}

	wanted := make(map[string]bool, len(defaultLabelSpec))
	for _, spec := range defaultLabelSpec {
		wanted[spec.Name] = true
	}

	configured := make(map[string]bool, len(defaultLabelSpec))
	for _, item := range resp.Items {
		var project internalVikunja.VikunjaProject
		if err := json.Unmarshal(item, &project); err != nil {
			return err
		}
		if project.Id >= -1 || !wanted[project.Title] {
			continue
		}
		if configured[project.Title] {
			return fmt.Errorf("multiple saved filters named %q found", project.Title)
		}

		projectID := strconv.FormatInt(project.Id, 10)
		views, err := adapter.ListProjectViews(projectID)
		if err != nil {
			return err
		}
		kanbanViewID, err := findKanbanViewID(views)
		if err != nil {
			return err
		}
		if err := configureKanbanBuckets(adapter, projectID, kanbanViewID); err != nil {
			return err
		}
		configured[project.Title] = true
	}

	for _, spec := range defaultLabelSpec {
		if !configured[spec.Name] {
			return fmt.Errorf("default saved filter %q not found", spec.Name)
		}
	}
	return nil
}

func configureKanbanBuckets(adapter *internalVikunja.Adapter, projectID, kanbanViewID string) error {
	buckets, err := adapter.ListProjectViewBuckets(projectID, kanbanViewID)
	if err != nil {
		return err
	}
	idByTitle, err := indexBucketsByTitle(buckets)
	if err != nil {
		return err
	}

	for _, s := range kanbanBucketSpec {
		if bucketID, ok := idByTitle[s.Name]; ok {
			if _, err := adapter.UpdateBucket(projectID, kanbanViewID, bucketID, s.Name, s.Position); err != nil {
				return err
			}
			continue
		}
		if s.RenameFrom != "" {
			bucketID, ok := idByTitle[s.RenameFrom]
			if !ok {
				return fmt.Errorf("default bucket %q not found", s.RenameFrom)
			}
			if _, err := adapter.UpdateBucket(projectID, kanbanViewID, bucketID, s.Name, s.Position); err != nil {
				return err
			}
			continue
		}
		if _, err := adapter.CreateBucket(projectID, kanbanViewID, s.Name, s.Position); err != nil {
			return err
		}
	}
	return nil
}

// ListWorkspaces 列出全部真实项目，并转换为 Workspace 列表。
// Vikunja 使用负 ID 表示 saved filter 等伪项目。
func ListWorkspaces() ([]Workspace, error) {
	adapter, err := internalVikunja.NewAdapter()
	if err != nil {
		return nil, err
	}

	projectsResp, err := adapter.ListProjects()
	if err != nil {
		return nil, err
	}

	var resp internalVikunja.VikunjaResp
	if err := decode(projectsResp, &resp); err != nil {
		return nil, err
	}

	translator := NewScrumTranslator()
	workspaces := []Workspace{}
	for _, item := range resp.Items {
		var vp internalVikunja.VikunjaProject
		if err := json.Unmarshal(item, &vp); err != nil {
			return nil, err
		}
		if vp.Id >= 0 {
			workspaces = append(workspaces, translator.TranslateVikunjaProject2Workspace(&vp))
		}
	}
	return workspaces, nil
}

// ListTasks 按 bucket 标题（filter，如 "TO DO"）列出工作区任务，并转换为 Task 列表。
func ListTasks(workspaceID, filter string) ([]Task, error) {
	adapter, err := internalVikunja.NewAdapter()
	if err != nil {
		return nil, err
	}

	views, err := adapter.ListProjectViews(workspaceID)
	if err != nil {
		return nil, err
	}
	kanbanViewID, err := findKanbanViewID(views)
	if err != nil {
		return nil, err
	}

	buckets, err := adapter.ListProjectViewBuckets(workspaceID, kanbanViewID)
	if err != nil {
		return nil, err
	}
	idByTitle, err := indexBucketsByTitle(buckets)
	if err != nil {
		return nil, err
	}
	bucketID, ok := idByTitle[filter]
	if !ok {
		return nil, fmt.Errorf("bucket %q not found in workspace %s", filter, workspaceID)
	}

	tasksResp, err := adapter.ListProjectTasksByBucket(workspaceID, bucketID)
	if err != nil {
		return nil, err
	}

	var resp internalVikunja.VikunjaResp
	if err := decode(tasksResp, &resp); err != nil {
		return nil, err
	}

	translator := NewScrumTranslator()
	tasks := []Task{}
	for _, item := range resp.Items {
		var vt internalVikunja.VikunjaTask
		if err := json.Unmarshal(item, &vt); err != nil {
			return nil, err
		}
		tasks = append(tasks, translator.TranslateVikunjaTask2Task(&vt))
	}
	return tasks, nil
}

// UpdateTaskStatus moves a task into the Scrum bucket named by status.
func UpdateTaskStatus(workspaceID, taskID, status string) (any, error) {
	workspaceID = strings.TrimSpace(workspaceID)
	workspaceIDValue, err := strconv.ParseInt(workspaceID, 10, 64)
	if err != nil || workspaceIDValue <= 0 {
		return nil, framework.InvalidFormatException("workspace ID must be a positive integer")
	}

	taskID = strings.TrimSpace(taskID)
	taskIDValue, err := strconv.ParseInt(taskID, 10, 64)
	if err != nil || taskIDValue <= 0 {
		return nil, framework.InvalidFormatException("task ID must be a positive integer")
	}
	taskID = strconv.FormatInt(taskIDValue, 10)

	status = strings.TrimSpace(status)
	if status == "" {
		return nil, framework.InvalidFormatException("task status is required")
	}
	if !isScrumStatus(status) {
		return nil, framework.InvalidFormatException(fmt.Sprintf("unsupported task status %q", status))
	}

	adapter, err := internalVikunja.NewAdapter()
	if err != nil {
		return nil, err
	}

	taskResp, err := adapter.GetTask(taskID)
	if err != nil {
		return nil, err
	}
	var task internalVikunja.VikunjaTask
	if err := decode(taskResp, &task); err != nil {
		return nil, err
	}
	if task.ProjectId != workspaceIDValue {
		return nil, framework.NotFoundException(fmt.Sprintf("task %s not found in workspace %s", taskID, workspaceID))
	}

	if err := moveTaskToStatus(adapter, workspaceID, taskID, status); err != nil {
		return nil, err
	}
	if err := syncTaskStatusToLabelSavedFilters(adapter, task.Labels, taskID, status); err != nil {
		return nil, err
	}
	return framework.RestOK, nil
}

func moveTaskToStatus(adapter *internalVikunja.Adapter, projectID, taskID, status string) error {
	views, err := adapter.ListProjectViews(projectID)
	if err != nil {
		return err
	}
	kanbanViewID, err := findKanbanViewID(views)
	if err != nil {
		return err
	}

	buckets, err := adapter.ListProjectViewBuckets(projectID, kanbanViewID)
	if err != nil {
		return err
	}
	idByTitle, err := indexBucketsByTitle(buckets)
	if err != nil {
		return err
	}
	bucketID, ok := idByTitle[status]
	if !ok {
		return fmt.Errorf("bucket %q not found in project %s", status, projectID)
	}

	_, err = adapter.MoveTaskToBucket(projectID, kanbanViewID, bucketID, taskID)
	return err
}

func syncTaskStatusToLabelSavedFilters(adapter *internalVikunja.Adapter, labels []internalVikunja.VikunjaLabel, taskID, status string) error {
	defaultLabels := make(map[string]bool, len(defaultLabelSpec))
	for _, spec := range defaultLabelSpec {
		defaultLabels[spec.Name] = true
	}

	wanted := make(map[string]bool, len(labels))
	for _, label := range labels {
		if defaultLabels[label.Title] {
			wanted[label.Title] = true
		}
	}
	if len(wanted) == 0 {
		return nil
	}

	projects, err := adapter.ListProjects()
	if err != nil {
		return err
	}
	var resp internalVikunja.VikunjaResp
	if err := decode(projects, &resp); err != nil {
		return err
	}

	synced := make(map[string]bool, len(wanted))
	for _, item := range resp.Items {
		var project internalVikunja.VikunjaProject
		if err := json.Unmarshal(item, &project); err != nil {
			return err
		}
		if project.Id >= -1 || !wanted[project.Title] {
			continue
		}
		if synced[project.Title] {
			return fmt.Errorf("multiple saved filters named %q found", project.Title)
		}

		projectID := strconv.FormatInt(project.Id, 10)
		if err := moveTaskToStatus(adapter, projectID, taskID, status); err != nil {
			return fmt.Errorf("move task %s in saved filter %q: %w", taskID, project.Title, err)
		}
		synced[project.Title] = true
	}

	for title := range wanted {
		if !synced[title] {
			return fmt.Errorf("saved filter %q not found for task label", title)
		}
	}
	return nil
}

func isScrumStatus(status string) bool {
	for _, bucket := range kanbanBucketSpec {
		if bucket.Name == status {
			return true
		}
	}
	return false
}

// findKanbanViewID 从项目 views 分页响应（Resp<View>）中提取 kanban view id。
func findKanbanViewID(views any) (string, error) {
	var resp internalVikunja.VikunjaResp
	if err := decode(views, &resp); err != nil {
		return "", err
	}
	for _, item := range resp.Items {
		var view struct {
			ID       int64  `json:"id"`
			ViewKind string `json:"view_kind"`
		}
		if err := json.Unmarshal(item, &view); err != nil {
			return "", err
		}
		if view.ViewKind == "kanban" {
			return strconv.FormatInt(view.ID, 10), nil
		}
	}
	return "", fmt.Errorf("kanban view not found")
}

// extractKanbanView 从创建项目响应中提取 projectID 与 kanban view id。
func extractKanbanView(created any) (projectID, kanbanViewID string, err error) {
	var proj struct {
		ID    int64 `json:"id"`
		Views []struct {
			ID       int64  `json:"id"`
			ViewKind string `json:"view_kind"`
		} `json:"views"`
	}
	if err := decode(created, &proj); err != nil {
		return "", "", err
	}

	projectID = strconv.FormatInt(proj.ID, 10)
	for _, v := range proj.Views {
		if v.ViewKind == "kanban" {
			return projectID, strconv.FormatInt(v.ID, 10), nil
		}
	}
	return "", "", fmt.Errorf("kanban view not found in created project")
}

// indexBucketsByTitle 从 bucket 列表响应中按 title 建立 id 映射。
func indexBucketsByTitle(buckets any) (map[string]string, error) {
	var list struct {
		Items []struct {
			ID    int64  `json:"id"`
			Title string `json:"title"`
		} `json:"items"`
	}
	if err := decode(buckets, &list); err != nil {
		return nil, err
	}

	m := make(map[string]string, len(list.Items))
	for _, b := range list.Items {
		m[b.Title] = strconv.FormatInt(b.ID, 10)
	}
	return m, nil
}

// decode 将 adapter 返回的 any 结果转换到类型化结构体。
func decode(src any, dst any) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}
