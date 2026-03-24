---
name: execute
description: 执行单文件计划 `.plans/{plan_id}-{需求名}.md`。仅执行计划，将任务标记为完成，不修订计划。
metadata:
  version: 0.0.2
---

# Execute - 执行计划

**MUST** 仅执行计划，将任务标记为完成，不修订计划。

## Workflows

### 1. 顺序执行 - Execute Sequentially

此模式用于根据已有的 Plan 文件，按部就班地执行任务。

- **输入**：用户的执行指令，通常到 `Phase` 粒度（例如“执行 Phase 1”）。
- **过程**：
    1.  AI Agent 在指定 `Phase` 下，按顺序执行 `Task`。
    2.  在执行中，若发现异常依赖或阻塞，应立即提示并停下等待用户确认。
- **输出**：代码实现、测试结果或其他交付产物。

**每个 Context Window 只执行一个 Phase**

- **默认**：一个 Phase 在一个新的 Context Window 中执行。
- **边界**：一个 Context Window 内默认只推进一个 Phase。
- **例外**：当某个 Phase 过大（例如超过上下文容量或需要大量并行信息）时，由 Human 明确指示拆分到多个 Context Window 连续执行；拆分后仍视为同一个 Phase，必须从上一次的断点继续，且不得跳过未完成的 Task。

### 2. 指定执行 - Execute Specified

如果用户明确指出需要执行特定的 Task，则使用此模式。此模式用于根据已有的 Plan 文件，执行指定的 `Task`。

执行完指定的 Task 后暂停。

## 执行模式

按 Phase / Task 两层结构拆分。

- **MUST** 每个 Phase 的第一条 Task 必须用于复核“上一个 Phase 的完成情况”（确认所有 Task 已标注完成、验收点已满足、以及上一个 Phase 的收尾动作已完成），且编号必须为 `TASK-a00`（其中 `a` 为 Phase 编号，例如 Phase 1 为 `1xx`）；若不满足则暂停并回到上一个 Phase 补齐。
- **MUST** 每个 Phase 的最后一条 Task 必须用于判断“当前 Phase 是否需要返修 Design/Specs”（例如发现关键变更、设计口径不一致、验收标准需调整），且编号必须为 `TASK-a99`（其中 `a` 为 Phase 编号，例如 Phase 1 为 `1xx`）；若需要返修则暂停并先更新 Specs 与用户确认，再重新拆分/调整后续 Tasks。
- **MUST** Phase 0 为特殊 Phase：用于确认每个 Phase 的执行准备是否已就绪（尤其是拆分是否合理、验收是否可验证、依赖是否明确）。其中 `HITL-000` 用于确认整体目标，`HITL-00x`（例如 `HITL-001`）用于确认对应的 Phase a（例如 Phase 1）的拆分是否合理。

- **MUST** 必须严格按顺序执行任务，并从第一个 `- [ ]` 开始。
- **MUST** 必须一次只执行一个 Task，完成后暂停并等待下一步指示。
- **MUST** 必须在完成 Task 后将对应条目从 `- [ ]` 更新为 `- [x]`。
- **MUST** 发生错误时必须立即停止执行，并等待用户指示。
- **MUST** 执行过程中如识别到关键变更，必须立即暂停，先更新 Design 并与用户确认，再重新按 Phase / Task 拆分并调整顺序。

- **MUST NOT** 跳过任务。
- **MUST NOT** 不按顺序执行。
- **MUST NOT** 执行任务列表之外的工作。
- **MUST NOT** 出错后继续执行。

##  Context Window / Phase 执行约束：

- **默认**：一个 Phase 在一个新的 Context Window 中执行。
- **边界**：一个 Context Window 内只推进一个 Phase。
- **例外**：如果执行过程中识别到某个 Phase 过大（例如超过上下文容量或需要大量并行信息），暂停并提示 HITL 指示拆分到多个 Phase。

当一个 Phase 完成后暂停，执行下面的检查：
- 确认所有 Task 已按预期完成并在文档中标注完成状态。
- 如需测试：将测试拆为独立 Task 并在 Tasks 中跟进；完成后在此确认已全部通过。
- 如需更新文档：将文档更新拆为独立 Task 并在 Tasks 中跟进；完成后在此确认已更新。
