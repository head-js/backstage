---
name: plan
description: 初始化或修订单文件计划 `.plans/{plan_id}-{需求名}.md`（Requirements / Specs / Design / Phases 四段）。产出可执行、可验收、可追溯的计划文件；仅修订计划，不执行计划。
metadata:
  version: 0.0.3
---

# Plan - 计划 & 修订计划

**MUST** 仅修订计划，不执行计划，不改动除 `.plans/{plan_id}-{需求名}.md` 之外的文件。

## Workflow

维护约定：各部分的结构示例见 `references/` 下的 `template-xxx.md`；本文件（`COMMAND.md`）是规则口径的权威来源。若后续需要调整规则，应直接修改本文件；若需要同步示例，再相应更新 `template-xxx.md`，避免在 how-to/example 中复制出多份规则导致口径漂移。

### 1. 初始化 (Initialize)

此模式用于从零开始创建一份新的 Plan 文件。

- **输入**：用户提供尽可能完整的 Requirements、Design 和相关约束。
- **过程**：
  1. **收集输入**：确认 Requirements + Design + 约束。
  2. **执行拆分**：基于输入和模板，产出包含 Requirements / Specs / Design / Phases 的单文件 `.plans/{plan_id}-{需求名}.md`。
  3. **用户确认**：与用户确认单文件内容是否可执行、可验收、可追溯。
  4. **循环改进**：若用户不接受，则回到第 2 步，迭代优化直至满足要求。
- **输出**：一份可执行的 Plan 文件。

如果输入不足以完成初始化，你只能指出缺口并向用户索要最小必要信息，再开始。

### 2. 自动修订 (Auto Revise)

此模式用于 AI 基于「新信息」对现有 Plan 进行自动修订。

新信息来源包括但不限于：代码实现的既定事实、补充/更改的 Phases、执行反馈（日志/现象/截图等）。

- **输入**：上述新信息 + 你希望调整的范围（仅 Phases / 同步 Specs + Design / 同步 Requirements + Specs + Design + Phases）。
- **过程**：
  1. **识别变化点**：提炼新增/变更的需求点、约束、风险与验收标准。
  2. **优先更新 Phases**：将变化落到 Phases（必要时新增/调整 Phase），并保证每个 Phase 的目标、边界、交付物与验收点可验证。
  3. **反向对齐 Specs/Design**：当 Phases 的变化引入新设计决策或改变既有决策时，同步更新 Specs 与 Design，确保与 Phases 一致。
  4. **必要时更新 Requirements**：仅当变化影响目标/范围/约束/验收口径时，才更新 Requirements，并同时回查 Specs/Design/Phases 的一致性。
  5. **产出差异摘要**：给出本次修订的变更点清单（改了什么、为什么、影响哪些验收点）。
- **输出**：更新后的 `.plans/{plan_id}-{需求名}.md`（唯一权威来源）。

### 3. 人工修订 (Human Revise)

此模式用于修订现有的 Plan 文件。

- **输入**：用户提供对现有 Plan 文件的修改意见或新增需求。
- **过程**：
  1. **收集输入**：确认修改意见或新增需求。
  2. **执行修订**：基于输入和模板，更新 `.plans/{plan_id}-{需求名}.md` 文件。
  3. **用户确认**：与用户确认修订后的内容是否可执行、可验收、可追溯。
  4. **循环改进**：若用户不接受，则回到第 2 步，迭代优化直至满足要求。
- **输出**：一份可执行的 Plan 文件。

- **触发**：用户在评审计划时，提出质疑、补充需求或新想法。
- **联动原则**：修订不是 `Requirements → Specs → Design → Phases` 的单向过程，而是四者的往复联动；当用户在任一部分形成了新决策或已有实现时，AI 必须基于该事实反向补全其余三部分，并保持整份文档的一致性、可执行性与有效性。
- **约束**：遵循 Plan-First 原则（Plan 变更流程）
  1. **暂停并讨论**：立即暂停，与用户讨论质疑点/变更点。
  2. **更新 Plan**：将讨论结果更新到 `.plans/{plan_id}-{需求名}.md` 文件中。
     - 对于小范围调整，应更新 `Specs`、`Design` 和 `Phases`。
     - 对于涉及根本性变更的大范围调整，可能需要修订 `Requirements` 和 `Design`。
  3. **确认 Plan**：与用户确认更新后的 Plan 内容。
  4. **修改代码**：只有在 Plan 确认后，才能根据新的 Plan 开始或继续修改代码。严禁跳过 Plan 更新直接修改代码。

如果输入不足以完成修订，指出缺口并向用户索要必要信息，再开始。

### 4. Phase 拆分规范

使用 **World Model** 中的 Task 定义，即：

Task 是项目执行的最小粒度单元，AA 可以完整执行一个 Task 并完成其所有步骤。Task 应当遵循如下数据结构：

<Task type="Model">
- [ ] **TASK-xxx**: {Task Description}
  - **Dependencies**:
    - TASK-xxx
    - TASK-xxx
  - **Do**: 
    - {明确“要执行的动作”以及“预期产出物/交付物”（推荐用 `动作 -> 产出` 的形式）}
    - {要执行的动作}
    - {要推进的事项}
  - **Check**: 
    - {明确“如何验收/验证”，并给出可观测信号/指标与通过/失败阈值}
    - {验收/验证标准}
    - {观测信号/指标}
    - Conditions:
      - IF PASS   将 Act 设置为: PASS and Continue to 下一条 TASK-xxx
      - ELSE FAIL 将 Act 设置为: Pause and HITL，{阻塞点/需要决策的问题}
  - **Act**: {根据 Check Conditions 设置}
  - **禁止添加字段，使用 Do 和 Check 描述和执行任务**
</Task>

## MUST Follow Rules

### Rule 1 - Core Principles

1. **先产出再确认**：先给出单文件产物（`Requirements` / `Specs` / `Design` / `Tasks`），再确认是否可执行、可验收、可追溯；避免只讨论原则。
2. **结果与验证优先**：关注“可执行、可验收、可追溯”，验收与检查必须落到可测试/可检查标准（可复现、可度量、可对照）。
3. **粒度可控**：默认以最小可落地且可验收的形式输出，避免不必要的冗长。
4. **循环改进**：若用户不接受结论，回到“怎么拆”继续迭代，而不是在原地争论。
5. **单文件交付**：最终产物只有一个 Markdown 文件（例如 `.plans/{plan_id}-{需求名}.md`），内部只包含 `Requirements` / `Specs` / `Design` / `Tasks` 四段（结构以模板为准）。
6. **初始化合并模板**：初始化 `.plans/{plan_id}-{需求名}.md` 时，将 `references/` 下多个模板片段合并为单文件对应段落；`template-xxx.md` 为模板来源，应从中学习与复制结构，避免自行发明口径导致漂移。

### Rule 2 - 拼装规则（来源 / 过程 / 原则 / 输出目标）

本命令的核心工作不是“确认模板”，而是把来自多个参考文件的结构与口径**拼装**为单文件 `.plans/{plan_id}-{需求名}.md`，并在修订时保持四段内容一致。

- 拼装来源（必须遵循其结构与口径）：

- `.agents/commands/plan/references/template-requirements.md`
- `.agents/commands/plan/references/template-specs.md`
- `.agents/commands/plan/references/template-design.md`
- `.agents/commands/plan/references/example-phases.md`

补充参考（需要时参考，但不作为强制模板）：

- 方法论：`.agents/commands/plan/references/what-is-ears-format.md`
- 示例：`.agents/commands/plan/references/001-example.md`

拼装过程（初始化 / 修订共用）：

1. 将上述 `template-xxx.md` 的段落结构合并到单文件中，确保只包含 `Requirements` / `Specs` / `Design` / `Phases` 四个段落。
2. 将用户提供的信息落到对应段落；若用户只给了其中一段，也必须联动检查另外三段并做最小必要更新。
3. 当代码实现或相关决策已成为既定事实且无争议时，以代码事实为准反推并修订 `Requirements` 与 `Design`。但当用户对代码实现提出质疑、或识别出新需求/新想法时，必须遵循 Plan-First 原则（Plan 变更流程）。

拼装原则：

- 单一事实来源：所有结论都必须落回 `.plans/{plan_id}-{需求名}.md`，不要在聊天里维护另一个版本。
- 最小必要变更：只改与新信息/新决策相关的最小范围，避免无关重写导致迭代发散。
- 一致性优先：任何一处修改都必须检查四段（R/D/S/T）是否仍对齐、是否可执行与可验收。

输出目标：

- 产出或更新单文件 `.plans/{plan_id}-{需求名}.md`（唯一权威来源）。
- 让该文件达到“可执行、可验收、可追溯”的最低可用形态；不足之处通过 Tasks 明确下一步补齐。

## Single File Example（模板）

**Single File of Truth（单一事实来源）**：`.plans/{plan_id}-{需求名}.md` 是唯一权威来源。

- 任何讨论、迭代与修改建议，都必须以该文件的当前内容为准。
- 不要维护“另一个版本”的 Requirements/Specs/Design/Tasks；需要变更就直接反映到该文件结构对应段落里。

最终只输出一个文件（示例文件名：`.plans/{plan_id}-{需求名}.md`），其内容结构建议如下（段落标题大小写与 `template-xxx.md` 保持一致）：

```md
# {需求名}

## Requirements

目标、范围、约束、验收标准

（模板：`.agents/commands/plan/references/template-requirements.md`）

## Specs

规格说明与验收标准

（模板：`.agents/commands/plan/references/template-specs.md`）

## Design

可落地的设计信息

（模板：`.agents/commands/plan/references/template-design.md`）

## Phases

关键任务列表（每个任务写清"做什么、为什么、完成如何验收"）

（模板：`.agents/commands/plan/references/example-phases.md`）
```
