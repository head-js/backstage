# Gitea API Wrapper CLI

> updated_by: Cline - Claude-Sonnet-3.5
> updated_at: 2026-03-24 22:11:00

## Requirements

<!-- 模板：`.agents/commands/plan/references/template-requirements.md` -->

### Goals

- **G-001**: 实现一个轻量 CLI 工具 `backstage-gitea`，通过 Gitea HTTP API 封装操作 Gitea
- **G-002**: 提供稳定的 JSON 输出供上层 Agent 解析
- **G-003**: 提供可扩展的通用 API 命令层，支持灵活调用 Gitea API

### Non-Goals

- **NG-001**: 不实现 Plan/Phase/Task 业务逻辑（由其他模块负责）
- **NG-002**: 不实现 Workflow state machine
- **NG-003**: 不实现 Task dependency graph
- **NG-004**: 不实现 Custom fields
- **NG-005**: 不实现 Advanced reporting
- **NG-006**: 不实现 UI

### Scope

- **S-001**: 通用 API 命令层（类 RESTFul 风格路径）
- **S-002**: 支持 macOS / Windows 跨平台
- **S-003**: 提供可扩展的 Adapter 层供上层调用

### Non-Scope

- **NS-001**: Plan/Phase/Task 业务逻辑（由其他模块负责）
- **NS-002**: Git 操作（push/pull/clone 等）

### Functional Requirements

#### 常规需求

- **FR-001**: CLI 应从环境变量 `BACKSTAGE_GITEA_URL` 和 `BACKSTAGE_GITEA_TOKEN` 读取 Gitea 连接信息
- **FR-002**: CLI 非 0 exit code 表示错误，并输出简洁错误消息
- **FR-003**: CLI 应支持通过类 RESTFul 风格路径调用任意 Gitea API 端点
- **FR-004**: Adapter 层应提供可扩展的接口供上层服务调用
- **FR-005**: CLI 输出格式由内部按需自行控制

### Success Metrics

| Metric | Current | Target | How to Measure |
|--------|---------|--------|----------------|
| CLI 命令执行成功 | N/A | 所有命令可正常执行 | 运行测试用例验证 |
| 跨平台构建成功 | N/A | macOS/Windows 均可编译 | Makefile 构建验证 |

### Dependencies

- **D-001**: go-gitea SDK
- **D-002**: Cobra CLI 框架
- **D-003**: ucarion/urlpath 路由解析

### Constraints

- **C-001**: 单文件 binary 发布，无需安装 runtime
- **C-002**: 二进制名称为 `backstage-gitea`（Windows 自动添加 .exe）
- **C-003**: 每个 Plan 对应一个独立 Repo

### Assumptions

- **A-001**: Gitea API 服务可访问
- **A-002**: 用户已配置有效的 API Token

## Specs

- [ ] **SPEC-001**: CLI 命令行接口
  - **背景 / 目标**：定义 CLI 命令结构，统一用户交互方式，提供通用 API 调用能力
  - **范围**：api 子命令（通用类 RESTFul 风格路径）
  - **关键决策**：采用 Cobra 框架，使用类 RESTFul 风格路径作为参数
  - **实现约束**：
    - 使用 spf13/cobra 构建 CLI
    - 使用 ucarion/urlpath 进行 RESTful 路径匹配
    - 支持任意 Gitea API 端点调用
  - **验收**：
    - [ ] `backstage-gitea --help` 正常输出帮助信息
    - [ ] `backstage-gitea api --help` 正常输出帮助信息
    - [ ] 可调用任意 Gitea API 端点

- [ ] **SPEC-002**: Gitea Adapter 层
  - **背景 / 目标**：统一管理 API 调用，隔离 SDK 细节，为上层提供可扩展接口
  - **范围**：所有 Gitea API 操作
  - **关键决策**：使用 go-gitea SDK，Adapter 层专注于 API 调用封装
  - **实现约束**：
    - 基于 go-gitea SDK 实现
    - 方法命名与 SDK 保持一致
    - 提供统一的错误处理和输出格式化
  - **验收**：
    - [ ] Adapter 可成功调用 ServerVersion
    - [ ] 可扩展支持更多 API 端点

- [ ] **SPEC-003**: 输出格式
  - **背景 / 目标**：CLI 输出格式由内部按需自行控制，供上层 Agent 解析
  - **范围**：所有命令的输出格式
  - **关键决策**：输出格式由 cmd 内部根据调用场景自行决定
  - **验收**：
    - [ ] 输出格式可被上层解析
    - [ ] human-readable 输出友好

## Design

### Architecture Overview

```
┌─────────────────┐
│  Upper Layer    │
│ (上层 Agent/服务) │
└────────┬────────┘
         │ 调用
         ▼
┌─────────────────┐
│ backstage-gitea │
│      CLI        │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│      Cmd        │
│  (Commands)     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Gitea Adapter │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   go-gitea SDK  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Gitea REST API  │
└─────────────────┘
```

> **说明**：本项目提供底层封装，上层的 Plan/Phase/Task 业务逻辑由其他模块负责。

### CLI Command Model

```
backstage-gitea
  └── api
        ├── <HTTP_METHOD> <API_PATH>
        │
        ├── GET /repos
        ├── GET /users/:username/repos
        ├── GET /repos/:owner/:repo
        ├── GET /repos/:owner/:repo/issues
        ├── GET /repos/:owner/:repo/milestones
        └── ... (可扩展任意 Gitea API)
```

### 项目结构

```
backstage-gitea
 ├ cmd
 │   ├ root.go
 │   └ api.go
 │
 ├ internal
 │   └ gitea
 │         ├ adapter.go
 │         └ api.go
 │
 ├ main.go
 └ version.go
```

### Error Handling

- 非 0 exit code 表示错误
- human-readable 模式：直接输出错误消息
- JSON 模式：`{"error": "error message"}`

## Phases

### PHASE-100: 基础项目搭建

本 Phase 聚焦于创建基础 Go 项目结构，集成 Gitea SDK 并验证连接。

- [x] 创建 Go 项目并验证 Hello World
- [x] 创建 Makefile 并验证跨平台构建
- [x] 集成 Gitea SDK 并测试 ServerVersion
- [x] 添加项目版本管理（version.go + ldflags）

### PHASE-200: Gitea Adapter 层开发

本 Phase 聚焦于创建 Gitea Adapter 层，统一管理 API 调用。

- [x] 创建 internal/gitea/adapter.go
- [x] 验证 Adapter 可正常调用 Gitea API

### PHASE-300: 通用 API 命令层

本 Phase 实现一个通用的 HTTP API 命令层，支持通过类 RESTFul 风格路径调用 Gitea API。

- [x] 添加依赖（Cobra, urlpath）并初始化 CLI 框架
- [x] 创建 cmd/api.go，实现通用 API 命令
- [x] 验证 `backstage-gitea api GET /repos` 正常工作
- [ ] 实现 ListUserRepos（GET /users/:username/repos）
- [ ] 实现 GetRepo（GET /repos/:owner/:repo）
- [ ] 实现 ListRepoIssues（GET /repos/:owner/:repo/issues）
- [ ] 实现 ListRepoMilestones（GET /repos/:owner/:repo/milestones）
- [ ] 实现 ListIssuesByMilestonePrefix（GET /repos/:owner/:repo/milestones/:prefix/issues）
- [ ] 实现 CreateRepo（POST /repos）
- [ ] 实现 DeleteRepo（DELETE /repos/:owner/:repo）

### PHASE-400: 验收与发布

本 Phase 完成最终验收和跨平台构建发布。

- [ ] 验收所有命令功能正常
- [ ] 验证输出格式正确
- [ ] 验证跨平台构建（macOS/Windows）
- [ ] 编写 README 文档
