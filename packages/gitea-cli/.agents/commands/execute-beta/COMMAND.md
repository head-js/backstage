---
name: execute-beta
description: 执行并完成任务 `.context/current-task.md`，将任务标记为完成，将相关的过程文档上传更新。
metadata:
  version: 0.0.3
---

# Execute Beta

**FAIL-FAST REQUIREMENT** Check if the tool exists, stop immediately and exit on any error.

```bash
$ backstage-gitea plan --help
```

## Workflow

- Step 1: 和该任务相关的 `Plan` 在 `.context/current-plan.md`。读取 `## Metadata` 章节，获取 `Plan` 的元信息。
- Step 2: 使用 `backstage-gitea` 工具下载最新的 `Plan` 并更新到 `.context/current-plan.md`。

```bash
$ backstage-gitea plan GET-MARKDOWN /:appName/:planName
```

- Step 3: 和该任务相关的 `Phase` 在 `.context/current-phase.md`。读取 `## Metadata` 章节，获取 `Phase` 的元信息。
- Step 4: 使用 `backstage-gitea` 工具下载最新的 `Phase` 并更新到 `.context/current-phase.md`。

```bash
$ backstage-gitea MARKDOWN /:appName/:planId/:phaseId
```

- Step 5: 和该任务相关的 `Task` 在 `.context/current-task.md`。读取 `## Metadata` 章节，获取 `Task` 的元信息。
- Step 6: 使用 `backstage-gitea` 工具下载最新的 `Task` 并更新到 `.context/current-task.md`。

```bash
$ backstage-gitea MARKDOWN /:appName/:planId/:phaseId/:taskId
```

- Step 7: 开始执行任务 `.context/current-task.md`。
- Step 8: 执行过程中如果发现需要更新 `Plan` 或 `Phase` 或 `Task`，应该在任务步骤中明确说明，并更新到 `.context/current-plan.md` 或 `.context/current-phase.md` 或 `.context/current-task.md` 或其它相关文档或代码。

- Step 9: 任务完成后，使用 `backstage-gitea` 工具将 `.context/current-task.md` 上传更新。

```bash
# IF Task Success
$ backstage-gitea SUCCESS /:appName/:planId/:phaseId/:taskId --context .context/current-task.md

# IF Task Failed
$ backstage-gitea FAILED /:appName/:planId/:phaseId/:taskId --context .context/current-task.md
```

- Step 10: 如果执行过程中有更新 `Phase` 或 `Plan`，使用 `backstage-gitea` 工具将 `.context/current-phase.md` 和 `.context/current-plan.md` 上传更新。

```bash
$ backstage-gitea PUT /:appName/:planId/:phaseId --context .context/current-phase.md

$ backstage-gitea PUT /:appName/:planId --context .context/current-plan.md
```

## Rules

**MUST** `Agent` 应该在项目根目录执行，即和 AGENTS.md 同目录。

**MUST** 应该在一个 `Agentic Session` 中执行。如果 Context 即将用尽，必须立刻停止并 Handoff 给 Human。这种场景说明 Task 拆分有误。

**MUST** 执行过程中如识别到关键变更，必须立刻停止并 Handoff 给 Human。这种场景说明 Design 需要更新。

- 按顺序执行任务，从第一个 `- [ ]` 开始。
- 一次只执行一个步骤条目。
- 完成后将对应条目从 `- [ ]` 更新为 `- [x]`。
- 确认所有步骤条目已按预期完成并标注完成状态。
