//go:build ignore

package main

type Workspace struct {
	Id          int64   `json:"Id"`
	Name        string  `json:"Name"`
	Description string  `json:"Description"`
	Priority    float64 `json:"Priority"`
}

type VikunjaProject struct {
	Id          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Position    float64 `json:"position"`
}

type VikunjaProjectList struct {
	Items []VikunjaProject `json:"items"`
}

type VikunjaTask struct {
	Id           int64                `json:"id"`
	ProjectId    int64                `json:"project_id"`
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
	Id        int64 `json:"id"`
	ProjectId int64 `json:"project_id"`
	Index     int64 `json:"index"`
}

// VisitEdge represents source blocking target using public VisitNode Id.
type VisitEdge struct {
	Source string `json:"Source"`
	Target string `json:"Target"`
}

type VikunjaTaskList struct {
	Items []VikunjaTask `json:"items"`
}

type VisitNode struct {
	Id          string              `json:"Id"`
	Workspace   string              `json:"Workspace"`
	Session     string              `json:"Session"`
	Title       string              `json:"Title"`
	Description map[string][]string `json:"Description"`
	Progress    string              `json:"Progress"`
}

type VisitNodeCreateRequest struct {
	Workspace   string              `json:"Workspace"`
	Title       string              `json:"Title"`
	Session     string              `json:"Session"`
	Description map[string][]string `json:"Description"`
}

// VisitNodeUpdateRequest carries optional fields. A nil pointer means the
// caller did not send the field and it must not be touched; a non-nil zero
// value (Progress 0, Session "") is an explicit update. Progress is a
// percentage integer (0-100).
type VisitNodeUpdateRequest struct {
	Id       string  `json:"Id"`
	Progress *int64  `json:"Progress"`
	Session  *string `json:"Session"`
}
