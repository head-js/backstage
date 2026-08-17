package scrum

// Workspace 工作区
type Workspace struct {
	Id          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Priority    float64 `json:"priority"`
}

// Task 任务
type Task struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Context     string `json:"context"`
	Priority    int64  `json:"priority"`
	WorkspaceId int64  `json:"workspaceId"`
}
