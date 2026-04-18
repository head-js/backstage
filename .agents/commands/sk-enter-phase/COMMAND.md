---
name: sk-enter-phase
description: 找到并进入 Current Phase，更新相关 Metadata 和 Context，将当前 Phase 拆分成具体的 Tasks 并上传。完成后立刻停止并退出。其他 Agent 会跟进并完成具体的任务。
metadata:
    version: 0.0.3
---

# Enter Phase

**FAIL-FAST REQUIREMENT** Check if the tools exist, stop immediately and exit on any error.

```bash
$ backstage-gitea agent --help
```

## Initialization

**MUST** 读取相关的 `Plan` 信息在 `.context/current-plan-metadata.yaml` 和 `.context/current-plan.md`。读取失败时必须立刻停止并退出。

## Workflow

- Step 1: 使用 `backstage-gitea` 工具找到并下载 Current Phase。

更新 Current Phase Metadata `.context/current-phase-metadata.md`。

```bash
$ backstage-gitea agent HEAD /:appId/:planId/current-phase > .context/current-phase-metadata.yaml
```

更新 Current Phase 内容 `.context/current-phase.md`。

```bash
$ backstage-gitea agent GET /:appId/:planId/current-phase > .context/current-phase.md
```

- Step 2: 按照以下格式将 Current Phase 拆分为具体的 Tasks，并更新到 `.context/current-phase.md`。

```markdown
- **TASK-101**: {task description}
  - **Do**:
    - {要执行的动作} -> {预期产出物/交付物}
  - **Check**:
    - {验收/验证标准}
    - {观测信号/指标}
```

- Step 3: 在 Task 的 `Do` 步骤中，所有产出物写入 `.context/current-task.md` 对应章节。**不允许在 `.context/` 下新建文件**。如果任务涉及 Design 变更（包括但不限于：调研、盘点、思考类任务产生的跨 Phase 结论；对架构、方案、技术选型的调整），应在 `Do` 步骤中注明写入 `.context/current-plan.md` → `Design` 对应章节。

    1. 明确产出物的格式和必填字段，明确待写入章节的位置（例如 `.context/current-task.md#Additional Thoughts`）

- Step 4: 按照 `Output Example` 的格式输出到用户指定文件。

- Step 5: 使用 `backstage-gitea` 工具将 `.context/current-phase.md` 的 `Tasks Breakdown` 章节创建到 Gitea 项目中。每个 Task 调用一次创建工具。

```bash
$ backstage-gitea plan POST /:appId/:planId/:phaseId/tasks --name "Task Name"
```

注意：不要在 Task 名称中包含 "Task" 前缀；taskId 应该由工具自动生成，你需要将这个 taskId 回填到原文档中。

注意：没有批量创建功能，必须逐个创建。即应该调用创建命令多次，例如：

```bash
$ backstage-gitea plan POST /cms-mgr/PLAN-102/PHASE-200/tasks --name "修复素材模块前导斜杠问题"

$ backstage-gitea plan POST /cms-mgr/PLAN-102/PHASE-200/tasks --name "修复商品模块前导斜杠问题"
```

- Step 6: 将 `.context/current-phase.md` 的内容追加到 `.context/current-plan.md` 文件尾部，然后使用 `backstage-gitea` 工具上传更新后的计划文件。

```bash
$ backstage-gitea plan PUT /:appId/:planId --context ".context/current-plan.md"
```

## Output Example

```markdown
<!-- // 注意 Phase 的头尾要保留分割线 -->
---

## PHASE-100: Phase Name

### Relevant Requirements

从 `.context/current-plan.md` 中提取和 Current Phase 相关的 Requirements。

### Relevent Specs

从 `.context/current-plan.md` 中提取和 Current Phase 相关的 Specs。

### Relevent Design

从 `.context/current-plan.md` 中提取和 Current Phase 相关的 Design。

### Tasks Breakdown

将 Current Phase 拆分为具体的 Tasks，按照 Step 2 规定的格式排列。

- **TASK-110**: {任务：例如补充日志/补齐缺失用例/小范围重构}
  - **Do**:
    - {执行动作} -> {产出物/中间产物}
  - **Check**:
    - {通过标准}
    - {观测信号}

- **TASK-120**: {task description}
  - **Do**:
    - {要执行的动作} -> {预期产出物/交付物}
  - **Check**:
    - {验收/验证标准}
    - {观测信号/指标}

- **TASK-130**: {task description}
  - **Do**:
    - {要执行的动作} -> {预期产出物/交付物}
  - **Check**:
    - {验收/验证标准}
    - {观测信号/指标}

---
<!-- // 注意 Phase 的头尾要保留分割线 -->
```
