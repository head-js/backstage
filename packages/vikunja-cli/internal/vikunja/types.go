package vikunja

import "encoding/json"

// VikunjaResp Vikunja 分页响应
type VikunjaResp struct {
	Schema     string            `json:"$schema"`
	Items      []json.RawMessage `json:"items"`
	Page       int64             `json:"page"`
	PerPage    int64             `json:"per_page"`
	Total      int64             `json:"total"`
	TotalPages int64             `json:"total_pages"`
}

// VikunjaTask Vikunja 任务（对应 Vikunja API Task）
type VikunjaTask struct {
	ProjectId   int64  `json:"project_id"`
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int64  `json:"priority"`
}

// VikunjaProject Vikunja 项目（对应 Vikunja API Project）
type VikunjaProject struct {
	Id              int64   `json:"id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	ParentProjectId int64   `json:"parent_project_id"`
	Position        float64 `json:"position"`
}
