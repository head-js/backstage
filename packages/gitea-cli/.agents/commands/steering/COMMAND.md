---
name: steering
description: 从当前项目中总结和提炼现状，形成项目公约（Steering）文档，为 AI 提供关于本项目的核心指导规则。
version: 0.0.1
---

# Steering - 项目公约

本能力用于从当前项目中总结和提炼现状，形成项目公约（Steering）文档。这些文档为 AI 提供关于本项目的背景知识、项目公约、产品方向、技术栈约束和代码结构约定，作为 AI 在后续任务中进行推理和决策的核心指导。

这些文档是项目的核心“规则”，比易变的“记忆”或临时的“日志”更具指导性和稳定性。AI 在执行任务时，会首先加载这些 Steering 文档来校准其行为，确保与项目目标和约束保持一致。

## 工作流程

本能力通过一个自主分析的工作流来创建或更新项目公约文档。

### 步骤 1: 全仓库分析

AI 将首先对整个仓库进行全面分析，阅读和理解包括代码、现有公约、计划文档 (`.plans/`)、命令和技能在内的所有内容，以形成对项目当前状态的整体认知。

### 步骤 2: 生成更新提案

基于分析结果，AI 将为 `.steering/` 目录下的四份核心公约文档生成更新提案。在生成提案时，AI 会严格遵守各项文档的既定稳定性：

-   **`product.md` (产品方向)**: 极稳定。除非项目发生重大业务转向，否则不应变更。
-   **`constitution.md` (项目公约)**: 非常稳定。通常只在团队遇到新问题并需要固化新规约时进行增补。
-   **`structure.md` (代码结构)**: 相对稳定。根据项目的实际演化进行调整。
-   **`design.md` (设计原则)**: 变化较多。反映跨项目的关键设计决策，并与 `.plans/` 中的细节设计保持同步，但不重复其细节。

### 步骤 3: 用户确认与写入

AI 将向您展示所有生成或修改的文档内容，供您审阅。在获得您的确认后，AI 才会将变更写入对应的文件。

**重要约定：** 所有 `steering` 文档均为 **动态文档 (Living Documents)**。在生成或更新这些文件时，末尾的 `Changelog` 部分必须遵循以下格式，不记录详细的变更历史：

```markdown
## Changelog

<!-- // 这是一个 Living Document，如无必要，无需维护变更历史。 -->
```

### 输出与模板

所有公约文档均基于参考模板生成，并输出到项目根目录的 `.steering/` 文件夹下。您可以在 `.agents/commands/steering/references/` 目录查看和修改这些模板：

| 文档 | 输出路径 | 模板路径 |
|---|---|---|
| `product.md` | `.steering/product.md` | `references/product.md` |
| `constitution.md` | `.steering/constitution.md` | `references/constitution.md` |
| `structure.md` | `.steering/structure.md` | `references/structure.md` |
| `design.md` | `.steering/design.md` | `references/design.md` |
