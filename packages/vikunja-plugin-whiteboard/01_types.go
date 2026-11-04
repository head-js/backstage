//go:build ignore

package main

type Workspace struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Priority    float64 `json:"priority"`
}

type VikunjaProject struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Position    float64 `json:"position"`
}

type VikunjaProjectList struct {
	Items []VikunjaProject `json:"items"`
}

type VikunjaTask struct {
	ID           int64                `json:"id"`
	ProjectID    int64                `json:"project_id"`
	Index        int64                `json:"index"`
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	PercentDone  float64              `json:"percent_done"`
	RelatedTasks VikunjaTaskRelations `json:"related_tasks"`
}

type VikunjaTaskRelations struct {
	Blocking []VikunjaTaskReference `json:"blocking"`
}

type VikunjaTaskReference struct {
	ProjectID int64 `json:"project_id"`
	Index     int64 `json:"index"`
}

// VisitEdge represents source blocking target using public VisitNode IDs.
type VisitEdge struct {
	Source string `json:"Source"`
	Target string `json:"Target"`
}

type VikunjaTaskList struct {
	Items []VikunjaTask `json:"items"`
}

type VisitNode struct {
	ID          string              `json:"Id"`
	Title       string              `json:"Title"`
	Progress    string              `json:"Progress"`
	Session     []string            `json:"Session"`
	Description map[string][]string `json:"Description"`
}

type VisitNodeCreateRequest struct {
	Workspace   string              `json:"Workspace"`
	Title       string              `json:"Title"`
	Session     []string            `json:"Session"`
	Description map[string][]string `json:"Description"`
}
