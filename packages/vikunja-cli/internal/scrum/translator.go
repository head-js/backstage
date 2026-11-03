package scrum

import internalVikunja "com.lisitede.backstage.vikunja/internal/vikunja"

// ScrumTranslator 用于将 Vikunja 对象转换为 Scrum 领域模型
type ScrumTranslator struct{}

// NewScrumTranslator 创建 ScrumTranslator 实例
func NewScrumTranslator() *ScrumTranslator {
	return &ScrumTranslator{}
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
