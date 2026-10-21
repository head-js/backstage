# Hub HTTP MCP 概念验证实施方案

> updated_by: Codex - GPT-6
> updated_at: 2026-11-02 19:06:50

## 目标与当前状态

在现有 Wails3 Hub 内增加使用 mark3labs/mcp-go 的 HTTP MCP 服务，以一个简单工具验证客户端接入可行性。本文完整覆盖此前的 Hub 初始化任务。

本轮仅更新本文档，不实施代码。后续 Agent 只负责 MCP 代码及必要的依赖声明和生命周期接线；MCP 安装、客户端配置和全部验证由人工执行。

使用顺序固定为：用户先手动启动 Hub 并保持应用运行，再启动 DSH 或使用 Codex 连接。常驻仅指用户不关闭应用，不增加托盘、关闭窗口隐藏、开机启动或系统守护。

## 职责与执行边界

| 工作 | 执行主体 | 说明 |
| --- | --- | --- |
| 编写 MCP package | Agent | 创建 MCP server、注册工具、提供 HTTP handler |
| HTTP 路由及退出接线 | Agent | 复用现有 listener，增加挂载与清理 |
| 必要的 Go 依赖声明 | Agent | 限于本功能，保留现有框架版本 |
| MCP 安装、注册和客户端配置 | 人工 | Codex、DSH 的配置与启用 |
| 安装依赖、构建与打包 | 人工 | 沿用现有 Hub 构建链路 |
| 启动与全部验证 | 人工 | 编译、测试、协议连接、工具发现、调用和退出检查 |

Agent 不执行安装、构建、测试、启动或验证命令，包括 go get、go mod tidy、go build、go test、Wails 构建、MCP Inspector、HTTP 验证请求、客户端调用和验证脚本。不新增测试专用程序或服务入口。

Agent 不修改 ~/.codex、~/.dsh 或其他客户端配置，不执行 MCP 注册或启用操作。依赖解析、下载及校验记录生成由人工处理，不编造 go.sum。

代码交付时列出修改内容和人工待办，明确说明“未构建、未运行、未验证”。人工尚未确认前，不得声称握手、客户端接入或验收通过。

## 最小协议契约

| 项目 | 约定 |
| --- | --- |
| 框架 | github.com/mark3labs/mcp-go，锁定已调研版本 v1.5.0 |
| 应用 | 现有 packages/hub Wails3 应用 |
| Go module | 复用 backstage.lisitede.com/hub |
| 传输 | Streamable HTTP |
| 路径 | /mcp |
| 默认端点 | http://127.0.0.1:20080/mcp |
| 地址来源 | 复用现有 HUB_HOST、HUB_PORT 和实际 listener 地址 |
| 服务名称 | backstage-hub |
| 服务版本 | 复用 config.Version |
| 唯一工具 | mcp_version |
| 参数 | 无业务参数，输入 schema 为 object，客户端传 {} |
| 结果 | 单个标准 MCP 文本块，内容严格为 Hello Backstage MCP |

旧讨论中的 mcp_tool_version 和 Hello World 已由上述名称和文本替代。结果使用标准 content 文本块，成功时 isError 为 false 或由 SDK 省略。文本不附加版本号、时间或其他前后缀。工具声明为只读、无破坏性、幂等，描述说明它用于检查 Backstage MCP 是否可调用。

初始化、能力发现、协议版本协商和连接状态交由 SDK 实现，不手写 JSON-RPC 分发或握手，不要求客户端固定使用某一版握手报文。不启用工具列表动态变更能力。

## 架构与文件改动

保持一个 Wails 应用进程、一个 Go module 和一个业务 HTTP listener。应用启动时创建一份 MCP server 和 handler，挂到现有 mux 的 /mcp 路径。客户端连接、断开或结束对话不启动或停止 Hub。

| 文件 | 具体改动 |
| --- | --- |
| packages/hub/src-wails/internal/mcpserver/server.go | 新增 Go package，创建 MCP server、注册唯一工具并构造 HTTP handler，提供挂载与清理能力 |
| packages/hub/src-wails/internal/server/server.go | 在现有 mux 挂载 /mcp，持有 MCP 组件并纳入现有关闭流程 |
| packages/hub/src-wails/go.mod | 声明固定版本的 mcp-go 依赖，保持现有 module 和 Wails/Huma 版本 |
| packages/hub/src-wails/go.sum | 仅接收真实依赖解析生成的必要记录；需人工生成时列为待办 |

mcpserver 不依赖 Wails UI。项目内部 server 与 mcp-go/server 同时导入时使用清晰别名。通过 mcp.NewTool 注册无参数工具，通过 mcp.NewToolResultText 返回约定文本，通过 server.WithToolCapabilities(false) 声明静态工具列表。

沿用此前讨论的 server.WithStateful(true) 作为旧协议客户端的 SDK session 配置；session 仅属于协议层，不拥有应用或 listener。新协议行为遵循 SDK 协商结果，不额外构建业务 session。

前端、Wails bindings、窗口关闭行为、开发入口和平台打包配置无需修改。保留现有 ping、version 和 OpenAPI 路由。MCP handler 直接挂载到 mux，不经过 Huma 的 REST 接口注册，也不通过 Wails bindings 对外提供。

## HTTP 生命周期接线

1. 沿用现有配置加载和同步端口绑定，在现有 mux 注册 /mcp。
2. 由现有 http.Server.Serve 处理 MCP 请求，不调用 MCP handler 自带的 Start() 再建 listener。
3. 保留当前启动及运行错误传播。绑定失败应明确报错，不能报告服务已就绪；初始化失败时清理已创建的 MCP 资源。
4. 在现有幂等 Stop() 中，使用独立且有界的清理 context 关闭 MCP sessions 和内部资源，再关闭外层 HTTP server，避免流式连接阻碍退出。使用 SDK 的 Shutdown(ctx) 或对应清理 API，外层 listener 仍由现有组件拥有。
5. 保留现有超时与强制关闭行为，确保 Wails OnShutdown 和 defer Stop() 重复调用不会重复释放或死锁。

## 客户端接入约定（仅供人工执行）

Codex 和 DSH 连接同一个已启动的 Hub。以下是人工接入依据，不属于 Agent 实施动作。

Codex 以 Streamable HTTP URL 注册，建议客户端名称为 backstage，不配置启动 Hub 的 command。

DSH 使用 @deepseek-ai/dsh-mcp-client，配置约定如下：

| 字段 | 值 |
| --- | --- |
| serverName | backstage |
| transport | streamable-http |
| url | 默认 http://127.0.0.1:20080/mcp，以实际地址为准 |
| failOnStartupError | true |
| 加载策略 | 可选 MCP 配置项，插件错误明确上报，局部故障不终止 DSH 应用启动 |

failOnStartupError: true 表示首次连接或工具发现失败时向 DSH 抛出 MCP 插件加载错误，不是整个 DSH 应用是否退出的开关。人工接入时保留可选加载策略，使错误明确且不会成为 DSH 的致命启动依赖。Hub 代码不控制此客户端策略，也不修改 DSH 本体或插件。

按先启动 Hub、后启动 DSH 的正常流程实施，本任务不讨论或扩展重试、异常启动顺序或后台异步连接。

## 后续代码实施顺序

1. 补充 Hub Go module 中的固定版本依赖声明。
2. 新增 internal/mcpserver，实现 MCP server、唯一工具和 HTTP handler。
3. 在现有 HTTP mux 挂载 /mcp，保留原有接口和监听配置。
4. 接入初始化失败清理和幂等退出路径。
5. 交付源码修改、依赖解析待办和人工验收说明，不执行安装、构建、启动或验证。

本轮仅覆盖本文档，不执行上述代码工作。

## 人工验收步骤与完成标准

以下步骤全部由人工执行，Agent 不代为执行或确认：

1. 安装必要依赖，按现有 Wails 链路构建真实 Hub 应用。
2. 手动启动 Hub 并保持应用运行，确认实际监听地址。
3. 人工注册或启用 Codex 的 HTTP MCP 配置，在加载配置的会话中连接 Hub。
4. 确认 Codex 发现 mcp_version，以 {} 调用并获得 Hello Backstage MCP。
5. 人工在 DSH 的生效配置中添加可选 MCP 配置项，在 Hub 已启动的前提下启动 DSH。
6. 确认 DSH 正常进入可用状态，发现 mcp__backstage__mcp_version 并实际调用，得到相同文本。
7. 确认两个客户端使用同一份 Hub 服务，客户端结束不影响 Hub，现有 HTTP API 仍可使用。
8. 用户退出 Hub，确认应用结束、HTTP 端口释放。

只验收真实 Wails 应用，不以独立 HTTP 服务、隐藏 CLI 或测试入口替代。只有人工确认后，才能记录连接、工具发现、客户端接入和运行验收通过。配置存在或普通 HTTP 请求成功均不等于 MCP 验收通过。

## 本次不包含

- Vikunja、Gitea 或其他业务工具。
- 独立 MCP 子项目、CLI、sidecar、额外 HTTP 端口或测试专用入口。
- 托盘、关闭窗口隐藏、开机启动、守护进程或自动拉起 Hub。
- OAuth、数据库、通配 CORS、动态工具注册、resources 或 prompts。
- Agent 执行的 MCP 安装、客户端配置和任何验证。
- DSH 启动机制修改、重试策略或极端场景扩展。
