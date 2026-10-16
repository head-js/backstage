# 新建 Go hub 应用与后续 Tauri 接入方案

> updated_by: Kilo - cc-switch/managed
> updated_at: 2026-10-16 10:13:00

## 实施状态

第一阶段已实施完成（2026-10-16）：

- Go 模块 `backstage/hub` 建于 `packages/hub/`，构建产物 `.dist/backstage-hub`（不提交仓库）。
- 启动命令 `backstage-hub serve [--host] [--port]`，默认 `127.0.0.1:20080`，日志输出实际监听地址。
- 接口：`GET /api/ping`、`GET /api/v1/version`、`GET /api/openapi`；未注册路径 404，无文档 HTML。
- SIGINT/SIGTERM 优雅退出，启动失败非零退出码。
- **偏差**：Huma v2 依赖无法通过 GOPROXY 镜像或直连拉取（BLAME-182），当前以 `net/http` 手写等价路由与最小 OpenAPI JSON 交付；Huma 接入待依赖问题解决后补齐。

## 当前任务

整体工作分为两个阶段：

1. 新建 Go hub 应用。
2. 将 hub 接入当前 Tauri 桌面应用。

当前仅执行第一阶段。本次先记录方案，不创建应用代码或修改桌面配置。

原始架构参考 [Design.md](Design.md)。本文按当前会话已确认的范围和命名收敛任务；与原始设计不一致时，以用户最新澄清为准。

## 已确认的约束

- Go 应用称为 hub，独立 Go 模块放在 `packages/hub/`。
- Go 构建产物统一命名为 `backstage-hub`，Windows 也不添加 `.exe`。
- 第二阶段手动将产物复制到 `packages/desktop/src-tauri/binaries/backstage-hub`，名称保持一致。
- 最终桌面产物为 `BackstageDesktop.app` 和 `BackstageDesktop.exe`。
- 实际使用环境是 Apple Silicon macOS 和 x86 Windows；这只是当前环境事实，不转化为平台限制或兼容矩阵。
- 正常编写和构建程序，不增加架构检测、二进制校验或额外兼容层。编译工具需要目标信息时，正常指定即可。
- 不新增构建或复制脚本，包括 `.mjs` 脚本；需要复制二进制时由人工处理。
- 当前不设计 Go 与 Rust 的通信。

## 第一阶段：新建 Go hub 应用

### 交付目标

建立一个能够独立构建、独立启动的 Go HTTP 服务，为后续桌面拉起提供真实可执行程序。

本阶段采用前台运行方式作为最小交付：执行命令后，当前进程持续提供 HTTP 服务。后台化、运行记录、服务发现、启动锁、状态与停止子命令、空闲退出暂不纳入本任务；后续需要时再按原始设计单独展开。

### 目录结构

```text
packages/hub/
├── go.mod
├── go.sum
├── cmd/
│   └── hub/
│       └── main.go          # 应用入口
└── internal/
    ├── cmd/                # Cobra 命令与依赖组装
    ├── config/             # 监听地址、端口等基础配置
    └── server/             # HTTP 服务、Huma 路由和基础接口
```

按实际需要创建文件，不提前增加业务服务、数据库、存储接口或通用框架。

### 启动入口

使用 Cobra 组织命令，提供：

```text
backstage-hub serve
backstage-hub serve --host 127.0.0.1 --port 0
```

- 默认监听回环地址，端口为 `0` 时由操作系统分配可用端口。
- 通过普通日志记录实际监听地址和启动错误；本阶段不定义供 Rust 解析的输出协议。
- 启动失败返回非零退出码。
- 处理当前平台正常的终止信号，关闭 HTTP 服务并释放资源；不额外设计跨进程停止机制。

### HTTP 服务

沿用原始设计的 `net/http` 与 Huma，提供最小接口：

| 接口 | 用途 |
| --- | --- |
| `GET /api/ping` | 确认 HTTP 服务能够响应 |
| `GET /api/v1/version` | 返回应用版本及 API 版本 |
| `/api/openapi` 下的机器可读契约 | 由 Huma 生成，具体格式路径按配置确定 |

关闭 Huma 文档 HTML 和独立 schema 路由。Go 不提供 HTML 页面、静态资源、SPA 路由或浏览器自动打开逻辑，未注册路径返回 404。

本阶段不增加业务接口、SSE、WebSocket、桌面探测协议或认证体系。

### 构建产物

在 `packages/hub/` 中使用正常 Go 命令构建：

```sh
go build -o .dist/backstage-hub ./cmd/hub
```

两个系统都显式指定输出文件名为 `backstage-hub`。文件内容由 Go 编译目标决定，文件名不携带架构或系统信息。

构建产物不提交源码仓库。本阶段不自动复制产物到桌面目录。

### 实施顺序

1. 创建独立 Go 模块和最小命令入口。
2. 增加基础配置、HTTP 服务与 Huma 路由。
3. 增加健康和版本接口、启动日志及正常退出处理。
4. 构建 `backstage-hub`，验证前台运行和接口行为。

### 验收

- 可以独立构建得到文件名准确的 `backstage-hub`。
- 可以直接执行并持续提供 HTTP 服务，不依赖 Tauri 或前端。
- 指定监听地址和端口生效，动态端口能通过日志获知。
- 健康、版本和机器可读 OpenAPI 可以访问。
- 未注册页面路径返回 404，没有文档 HTML 或静态页面服务。
- 启动错误能够定位，正常终止能够释放监听端口。
- 验证限于实际构建与运行行为，不向应用加入架构或二进制校验逻辑。

## 第二阶段：接入当前 Tauri

以下只记录后续方向，当前不实施：

1. 人工将 `backstage-hub` 放入 `packages/desktop/src-tauri/binaries/`。
2. 通过 Tauri 的 `bundle.resources` 携带 `binaries/`，保留文件名，不采用要求目标平台后缀的 `externalBin` 输入方式。
3. 新增 Rust `sidecar.rs`，通过资源绝对路径启动 `backstage-hub`。
4. 在 Tauri 初始化入口中调用一次，记录进程输出和启动错误。
5. 配置桌面产物名称为 `BackstageDesktop.app` 和 `BackstageDesktop.exe`。
6. 实机确认 Windows 无扩展名程序的启动行为。

此阶段聚焦随包交付和进程拉起；通信、自动重启及更完整的生命周期管理另行讨论。
