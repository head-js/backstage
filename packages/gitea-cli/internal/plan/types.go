package plan

import (
	"fmt"
	"slices"

	"com.lisitede.backstage.gitea/framework"
)

// GiteaExtra Gitea 对象的附加信息（Repo / Milestone 通用）
type GiteaExtra struct {
	Type        string   `json:"type"` // "REPO" | "MILESTONE" | "ISSUE"
	Id          int64    `json:"id"`   // Gitea 对象 ID
	No          int64    `json:"no"`   // Issue Number
	Title       string   `json:"title"`
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
	Title  string     `json:"title"`
	Phases []Phase    `json:"phases"`
	Gitea  GiteaExtra `json:"gitea"`
}

// Phase 阶段（对应 Gitea Milestone）
type Phase struct {
	Id     string     `json:"id"`
	Name   string     `json:"name"`
	Title  string     `json:"title"`
	Status string     `json:"status"`          // 状态：TODO / SUCCESS / UNKNOWN
	Tasks  []Task     `json:"tasks,omitempty"` // 任务列表
	Gitea  GiteaExtra `json:"gitea"`           // Gitea Milestone 信息
}

// Task 任务（对应 Gitea Issue）
type Task struct {
	Id      string     `json:"id"`
	Name    string     `json:"name"`
	Title   string     `json:"title"`
	Status  string     `json:"status"` // 状态：TODO / SUCCESS / FAIL / UNKNOWN
	Context string     `json:"context"`
	Gitea   GiteaExtra `json:"gitea"`
}

// Status 任务状态枚举
type Status string

const (
	StatusPass    Status = "PASS"
	StatusFail    Status = "FAIL"
	StatusTodo    Status = "TODO"
	StatusHold    Status = "HOLD"
	StatusUnknown Status = "UNKNOWN"
)

var allStatus = []Status{StatusPass, StatusFail, StatusTodo, StatusHold, StatusUnknown}

// ValidateStatus 校验状态字符串是否合法，合法则返回对应的 Status，否则返回错误
func ValidateStatus(status string) (Status, error) {
	s := Status(status)
	if !slices.Contains(allStatus, s) {
		return "", framework.InvalidFormatException("invalid status: " + status)
	}
	return s, nil
}

// Translate2Label 将状态转换为 Label 名称
// prefix: 前缀，如 "TASK"，返回 "TASK-PASS" / "TASK-FAIL" / "TASK-TODO" / "TASK-HOLD" / "TASK-UNKNOWN"
func (s Status) Translate2Label(prefix string) string {
	return fmt.Sprintf("%s/%s", prefix, s)
}

// StatusColor 状态与颜色映射
var StatusColor = map[Status]string{
	StatusPass: "#009800",
	StatusFail: "#e11d21",
	StatusTodo: "#fbca04",
	StatusHold: "#fef2c0",
}
