# [DEMO] Upload Image 接口字段命名应采用 snake_case 还是 camelCase

> updated_by: Kilo - GLM-5.1
> updated_at: 2026-06-13 11:58:00

### Category

决策歧义

### Where

cms-mgr/PLAN-102 / api/upload_image.go:42

### Observed

现有接口同时存在 snake_case（image_url）与 camelCase（imageUrl）两种命名，比例约 6:4。Upload Image 是 PLAN-102 新增的第一个端点，无既定先例可循。

### Expected

PLAN-102 或 cms-mgr 仓库内应有统一的 JSON 字段命名规范，便于直接套用。

### Tried

1. grep cms-mgr 内 *.go 发现两种风格共存；
2. 查 docs/api-style.md 未命中；
3. 询问 codex agent，无明确答复。

### Ask

- PLAN-102 的对外 JSON 字段应采用 snake_case 还是 camelCase？
- 是否需要同时落地一份 docs/api-style.md 作为后续基准？

### Impact

任务退出：blame 后立即退出，等待人工裁决后再继续；当前 PR 不推进。

### Refs

N/A
