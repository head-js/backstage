package scrum

import internalVikunja "com.lisitede.backstage.vikunja/internal/vikunja"

// ScrumTranslator 用于将 Vikunja 对象转换为 Scrum 领域模型
type ScrumTranslator struct{}

// NewScrumTranslator 创建 ScrumTranslator 实例
func NewScrumTranslator() *ScrumTranslator {
	return &ScrumTranslator{}
}

// workspacePriorityBase 估算的 position 上界（id 200 × 2^16），
// 用于把 VikunjaProject.Position 翻转为 Priority（数字大=优先）。
// 仅作排序键：单调递减变换，溢为负值不影响排序正确性。
const workspacePriorityBase = 200 * 65536

// TranslateVikunjaProject2Workspace 将 VikunjaProject 转换为 Workspace
// Id/Title/Description 直接透传；临时用 workspacePriorityBase-Position 填充 Priority（数字大=优先）
func (st *ScrumTranslator) TranslateVikunjaProject2Workspace(v *internalVikunja.VikunjaProject) Workspace {
	return Workspace{
		Id:          v.Id,
		Title:       v.Title,
		Description: v.Description,
		Priority:    workspacePriorityBase - v.Position,
	}
}

// TranslateVikunjaTask2Task 将 VikunjaTask 转换为 Task
// Name 用 VikunjaTask.Title 填充，Id 用 VikunjaTask.Id 填充，Context 用 VikunjaTask.Description 填充
func (st *ScrumTranslator) TranslateVikunjaTask2Task(v *internalVikunja.VikunjaTask) Task {
	return Task{
		Id:          v.Id,
		Name:        v.Title,
		Context:     v.Description,
		Priority:    v.Priority,
		WorkspaceId: v.ProjectId,
	}
}
