package plan

// GiteaExtra Gitea 对象的附加信息（Repo / Milestone 通用）
type GiteaExtra struct {
	Type        string   `json:"type"`        // "REPO" | "MILESTONE" | "ISSUE"
	Id          int64    `json:"id"`          // Gitea 对象 ID
	Name        string   `json:"name"`        // 名称
	Description string   `json:"description"` // 描述
	Body        string   `json:"body"`
	State       string   `json:"state"`
	Labels      []string `json:"labels"`    // Label 名称列表（Issue 专用）
	CreatedAt   string   `json:"createdAt"` // 创建时间（RFC3339）
	UpdatedAt   string   `json:"updatedAt"` // 更新时间（RFC3339）
}

// Plan 计划（对应 Gitea Repo）
type Plan struct {
	Id     string     `json:"id"`
	Name   string     `json:"name"`
	Phases []Phase    `json:"phases"`
	Gitea  GiteaExtra `json:"gitea"`
}

// Phase 阶段（对应 Gitea Milestone）
type Phase struct {
	Id     string     `json:"id"`              // Phase-xx 格式
	Name   string     `json:"name"`            // 阶段名称
	Status string     `json:"status"`          // 状态：TODO / SUCCESS / UNKNOWN
	Tasks  []Task     `json:"tasks,omitempty"` // 任务列表
	Gitea  GiteaExtra `json:"gitea"`           // Gitea Milestone 信息
}

// Task 任务（对应 Gitea Issue）
type Task struct {
	Id      string     `json:"id"`
	Name    string     `json:"name"`
	Status  string     `json:"status"` // 状态：TODO / SUCCESS / FAIL / UNKNOWN
	Context string     `json:"context"`
	Gitea   GiteaExtra `json:"gitea"`
}

// Status 任务状态枚举
type Status string

const (
	StatusPass   Status = "PASS"
	StatusFail   Status = "FAIL"
	StatusTodo   Status = "TODO"
	StatusHold   Status = "HOLD"
	StatusUnknown Status = "UNKNOWN"
)

// StatusColor 状态与颜色映射
var StatusColor = map[Status]string{
	StatusPass:   "#009800",
	StatusFail:   "#e11d21",
	StatusTodo:   "#fbca04",
	StatusHold:   "#fef2c0",
}
