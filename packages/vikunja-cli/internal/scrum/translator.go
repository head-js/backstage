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

// TranslateVikunjaTask2Task 将 VikunjaTask 转换为 Task。
// bucket 是调用方已经确定的当前查询 bucket；不能从跨视图的 Buckets 数组猜测。
func (st *ScrumTranslator) TranslateVikunjaTask2Task(v *internalVikunja.VikunjaTask, bucket string) Task {
	return Task{
		Id:          v.Id,
		Name:        v.Title,
		Context:     v.Description,
		Bucket:      bucket,
		Priority:    v.Priority,
		WorkspaceId: v.ProjectId,
	}
}
