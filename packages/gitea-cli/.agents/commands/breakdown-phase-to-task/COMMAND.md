---
name: breakdown-phase-to-task
metadata:
  version: 0.0.3
---

## Breakdown Phase to Task

## Workflow

- Step 1: 预期用户输入 `.workspace/plan-context.md` 或其它文件，加载并获得和 Plan 相关的项目知识。
- Step 2: 预期用户输入 `.workspace/phase-context.md` 或其它文件，加载并获得和 Current Phase 相关的项目知识。
- Step 3: 按照 `./references/example-task` 的格式将 Current Phase 拆分为具体的 Tasks。
- Step 4: 如果 Task 的内容是调研、盘点、思考类的任务，应该在任务步骤中说明输出到 `.workspace/phase-context.md` 的相应章节。
- Step 5: 按照 `Output Example` 的格式输出到用户指定文件。

## Output Example

```markdown
## Current Phase

### Metadata

- Id: PHASE-100
- Name: PHASE-100: Phase Name

### Relevant Requirements

从 `.workspace/current-context.md` 中提取和 current phase 相关的 Requirements。

### Relevent Specs

从 `.workspace/current-context.md` 中提取和 current phase 相关的 Specs。

### Relevent Design

从 `.workspace/current-context.md` 中提取和 current phase 相关的 Design。

### Tasks Breakdown

根据 `./references/example-task.md` 将 current phase 拆分为具体的 Tasks。

- [ ] **TASK-110**: {任务：例如补充日志/补齐缺失用例/小范围重构}
  - **Dependencies**:
    - TASK-102
  - **Do**:
    - {性动作} -> {产出物/中间产物}
  - **Check**:
    - {通过标准}
    - {观测信号}
  - **Act**:
    - IF  SUCCESS: 将 Act 更新为 Success and Continue
    - ELIF FAILED: 将 Act 更新为 FAILED and Handoff；立即停止执行，并报告失败原因、阻塞点、需要人工确认的决策点

- [ ] **TASK-120**: {task description}
  - **Dependencies**:
    - TEMP-101
  - **Do**:
    - {要执行的动作} -> {预期产出物/交付物}
  - **Check**:
    - {验收/验证标准}
    - {观测信号/指标}
  - **Act**:
    - IF  SUCCESS: 将 Act 更新为 Success and Continue
    - ELIF FAILED: 将 Act 更新为 FAILED and Handoff；立即停止执行，并报告失败原因、阻塞点、需要人工确认的决策点

- [ ] **TASK-130**: {task description}
  - **Dependencies**:
    - TASK-103
  - **Do**:
    - {要执行的动作} -> {预期产出物/交付物}
  - **Check**:
    - {验收/验证标准}
    - {观测信号/指标}
  - **Act**:
    - IF  SUCCESS: 将 Act 更新为 Success and Continue
    - ELIF FAILED: 将 Act 更新为 FAILED and Handoff；立即停止执行，并报告失败原因、阻塞点、需要人工确认的决策点
```
