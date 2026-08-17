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

	buckets, err := adapter.ListProjectViewBuckets(projectID, kanbanViewID)
	if err != nil {
		return nil, err
	}
	idByTitle, err := indexBucketsByTitle(buckets)
	if err != nil {
		return nil, err
	}

	for _, s := range kanbanBucketSpec {
		if s.RenameFrom != "" {
			bucketID, ok := idByTitle[s.RenameFrom]
			if !ok {
				return nil, fmt.Errorf("default bucket %q not found", s.RenameFrom)
			}
			if _, err := adapter.UpdateBucket(projectID, kanbanViewID, bucketID, s.Name, s.Position); err != nil {
				return nil, err
			}
			continue
		}
		if _, err := adapter.CreateBucket(projectID, kanbanViewID, s.Name, s.Position); err != nil {
			return nil, err
		}
	}

	return framework.RestOK, nil
}

// ListWorkspaces 列出全部工作区，并转换为 Workspace 列表。
// 排除 Vikunja 内置项目（Inbox、My Open Tasks）。
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
		if vp.Title == "Inbox" || vp.Title == "My Open Tasks" {
			continue
		}
		workspaces = append(workspaces, translator.TranslateVikunjaProject2Workspace(&vp))
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
	bucketID, ok := idByTitle[status]
	if !ok {
		return nil, fmt.Errorf("bucket %q not found in workspace %s", status, workspaceID)
	}

	if _, err := adapter.MoveTaskToBucket(workspaceID, kanbanViewID, bucketID, strconv.FormatInt(taskIDValue, 10)); err != nil {
		return nil, err
	}
	return framework.RestOK, nil
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
