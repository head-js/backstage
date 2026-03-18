package plan

// GiteaExtra Gitea 对象的附加信息（Repo / Milestone 通用）
type GiteaExtra struct {
	Type        string `json:"type"`        // "REPO" | "MILESTONE" | "ISSUE"
	Id          int64  `json:"id"`          // Gitea 对象 ID
	Name        string `json:"name"`        // 名称
	Description string `json:"description"` // 描述
	Body        string `json:"body"`
	State       string `json:"state"`
	CreatedAt   string `json:"createdAt"` // 创建时间（RFC3339）
	UpdatedAt   string `json:"updatedAt"` // 更新时间（RFC3339）
}

// Plan 计划（对应 Gitea Repo）
type Plan struct {
	Id     string     `json:"id"`
	Name   string     `json:"name"`
	Gitea  GiteaExtra `json:"gitea"`
	Phases []Phase    `json:"phases"`
}

// Phase 阶段（对应 Gitea Milestone）
type Phase struct {
	Id    string     `json:"id"`    // Phase-xx 格式
	Name  string     `json:"name"`  // 阶段名称
	Tasks []Task     `json:"tasks"` // 任务列表
	Gitea GiteaExtra `json:"gitea"` // Gitea Milestone 信息
}

// Task 任务（对应 Gitea Issue）
type Task struct {
	Id    string     `json:"id"`    // TASK-xxxx 格式
	Name  string     `json:"name"`  // 任务标题
	Gitea GiteaExtra `json:"gitea"` // Gitea Issue 信息
}
