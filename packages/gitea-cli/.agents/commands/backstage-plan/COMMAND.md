---
name: backstage-plan
description: 使用命令行工具 backstage-gitea 管理项目计划
metadata:
    version: 0.0.1
---

# Backstage Plan - Using `backstage-gitea` Cli

**MUST** 所有的动作必须通过命令行工具 `backstage-gitea` 执行。

**MUST** 如果遇到执行报错，必须立即停止。

**MUST** 如果遇到缺少功能，必须立即停止。

# Design

从 `.workspace/README.md` 中获取元数据。**MUST** 保持原始数据结构。

Phase Id 的格式是 `PHASE-xxx`。Phase Name 的格式是 `PHASE-xxx: Phase Name`。

Task Id 的格式是 `TASK-xxx`。Phase Name 的格式是 `TASK-xxx: Task Name`。

# Workflow

## Step 1: Workspace Metadata Alignment - 同步工作区元数据

1. 检查 Current Plan 的 Phase 是否已经创建完成
    - 检查元数据，还没有 Id 的 Phase 是待创建状态
    - 依次询问是否需要创建，在得到肯定的回答后，依次创建；每轮对话只处理一项

2. 检查 Current Phase 的 Task 是否已经创建完成
    - 检查元数据，还没有 Id 的 Task 是待创建状态
    - 依次询问是否需要创建，在得到肯定的回答后，依次创建；每轮对话只处理一项

3. 找到 Current Phase 的第一个待执行 Task，并下载详情
    - `status: PENDING` 即为待执行状态
    - 下载该 Task 详细内容，将 Task.Context 字段的内容更新到 `.workspace/Current-Task.md`

# Commands

## Show a Plan

使用命令 `backstage-gitea plan GET /{project}/{plan}` 获取计划的详情。

## Create a new Phase

使用命令 `backstage-gitea plan POST /{project}/{plan}/phases --name "{name}"` 新增一个 Phase，系统会自动分配 Phase Id。

## Create a new Task

使用命令 `backstage-gitea plan POST /{project}/{plan}/{phaseId}/tasks --name "{name}"` 新增一个 Task，系统会自动分配 Task Id。

## Show full content of a Task

使用命令 `backstage-gitea plan GET /{project}/{plan}/{phaseId}/{taskId}` 获取 Task 详情。
