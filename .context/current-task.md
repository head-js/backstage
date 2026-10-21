# Hub Wails 同进程应用第一步实施方案

> updated_by: HBR - GPT-6
> updated_at: 2026-10-21 12:50:00

## 当前任务与交付原则

将原计划的“建立 Wails 入口”和“接入同进程 Go 服务”合并为第一步：在 `packages/hub` 中建立能够真实构建、启动、提供现有 HTTP API 并正常退出的 Wails v3 应用。

本文覆盖原来的独立 Go 服务与 Tauri 接入计划。总体架构参考 [Design.md](Design.md)，以用户最新澄清为准。本次仅更新方案，尚未实施以下代码改造。

- 唯一交付产物是 Wails 应用，窗口、Go HTTP 服务和业务任务运行在同一个应用进程中。
- 不交付、不保留、不测试独立 CLI；不提供隐藏 CLI 或仅供测试启动内部服务的命令。
- 只测会交付的东西，只测真实的东西。未来 E2E 必须启动真实构建的 Wails 应用。
- 继续使用 `packages/hub`，新增 `src`、`src-wails` 及必要的开发、构建配置。
- 与 `packages/desktop` 无关，不读取其实现作为接入依据，不修改其代码、配置或资源，不复制任何产物到该目录。
- 不使用 Tauri，不启动 sidecar，不设计跨进程保活。真正退出应用时，HTTP 服务和后台任务一并结束。

## 参考方案与目录

参考 `/Users/shuaizhexu/head-js/preset-wails` 的前后端目录划分、React/Vite 开发链路、Wails 本地工具和平台构建配置。第一步先与参考仓库的 Wails `v3.0.0-beta.32` 对齐，显式锁定版本。

```text
packages/hub/
├── package.json
├── vite.config.js
├── scripts/
│   └── dev.mjs
├── src/
│   ├── index.html
│   ├── main.jsx
│   ├── App.jsx
│   └── styles.css
└── src-wails/
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── Taskfile.yml
    ├── internal/
    │   ├── config/
    │   │   └── config.go
    │   └── server/
    │       └── server.go
    ├── frontend/
    │   └── dist/
    └── build/
        ├── config.yml
        ├── Taskfile.yml
        ├── darwin/
        └── windows/
```

- 将现有 `go.mod`、`go.sum` 和业务 `internal` 代码迁入 `src-wails`，module 名称按 Go 域名正序惯例定为 `backstage.lisitede.com/hub`，维持一个 Go module，无需 `go.work`。
- 删除现有 `cmd/hub/main.go`、`internal/cmd/root.go`，移除 Cobra 依赖及失效的 CLI 构建、启动配置。
- `src-wails/main.go` 是唯一应用入口，直接组装并调用内部 Go 组件。
- `frontend/dist` 是嵌入应用的前端构建资源，生成产物忽略提交；`build` 下的平台配置和源资源正常提交。
- 原 `build-mac.sh` 调整或替换为 Wails 应用构建入口，不能继续把独立 HTTP 可执行程序作为构建结果。
- 新增依赖产生的工作区锁文件更新属于必要的依赖变更；应用实现与配置保持在 `packages/hub`。

## 应用组成与前端

应用进程包含原生窗口、系统 WebView、嵌入前端资源、Go HTTP 服务和后续业务任务。HTTP 服务仅作为应用内部组件存在，没有独立启动命令和交付产物。

前端采用 React + Vite。第一步提供最小界面，显示应用版本、实际 HTTP 服务地址和运行状态，便于观察真实应用的运行结果。

通过薄层 Wails 方法绑定读取应用内部状态；现有 HTTP 接口继续面向外部调用者。Wails 内部绑定不等同于外部 RPC、SSE 或 WebSocket。

## HTTP 服务与生命周期

现有阻塞式 `server.Run` 改造成可由应用管理的内部生命周期：

| 能力 | 职责 |
| --- | --- |
| 启动 | 同步监听端口，成功后在 goroutine 中提供 HTTP 服务 |
| 地址 | 返回实际监听地址，支持端口 `0` |
| 等待 | 获取 HTTP 服务意外终止的结果 |
| 停止 | 幂等优雅关闭，超时后强制关闭连接并等待服务退出 |

启动顺序：读取应用配置 → 同步启动 HTTP 监听 → 创建 Wails 应用和窗口 → 运行 Wails 事件循环 → 退出时统一清理资源。

- 监听失败时明确报告错误并结束启动，不能显示一个服务不可用的正常窗口。
- 后续初始化失败时，必须清理已经启动的 HTTP 服务。
- HTTP 服务意外终止时记录错误并结束应用，不能静默失效；覆盖窗口事件循环尚未启动的错误情形。
- 关闭服务使用独立的超时 context，不能复用已经取消的应用 context；清理必须幂等。
- 第一阶段关闭主窗口即退出应用。后续加入托盘后，再实现关闭窗口隐藏、真正退出时统一清理。
- 应用配置统一读取并传给内部组件，不保留原 CLI 参数入口。保留监听地址与端口的配置能力，具体配置载体在实施时按最小需要确定。

## 接口兼容与资源边界

保留现有 `net/http` 与 Huma，以及以下接口：

| 接口 | 用途 |
| --- | --- |
| `GET /api/ping` | 检查 HTTP 响应 |
| `GET /api/v1/version` | 返回应用版本和 API 前缀 |
| `GET /api/openapi.json`、`GET /api/openapi.yaml` | Huma 生成的机器可读契约 |

保留默认监听 `127.0.0.1:20080`、可配置监听地址和端口，以及端口 `0` 的实际地址获取能力。不照搬参考项目的 `0.0.0.0`、通配 CORS 或 `/api/version` 路径。

业务 HTTP 服务继续关闭文档 HTML 和独立 schema 路由，不承载前端页面；未注册路径返回 404。前端资源由 Wails 资源处理机制提供，正式构建时嵌入 Go 应用。

第一步不增加业务接口、SSE、WebSocket 或外部 RPC 实现；后续按原始需求单独设计。

## 开发与正式构建

开发入口：

```sh
pnpm --filter @head/backstage-hub dev
```

开发脚本在 `src-wails` 中调用 `go tool wails3 dev`，由 Wails 开发链路启动 Vite 和应用，支持前端热更新。Vite 属于开发工具进程；Wails 与业务 HTTP 服务始终位于同一个应用进程。Wails 工具通过 Go module 管理，不依赖全局安装。

正式构建先生成前端静态资源，再嵌入并构建 Wails 应用。交付应用必须能够直接启动，不依赖 Vite、Node 或额外的 Go 服务程序。平台配置参考仓库整理，补齐真实前端构建与资源嵌入，不能照搬其空的前端构建任务。

第一步建立 macOS、Windows 必要构建配置；实际验收仅声明已在对应环境运行验证的结果。签名、安装器和自动更新不在本步展开。

## 测试与验收原则

未来 E2E 的唯一被测对象是构建后的 Wails 应用，测试路径为：

```text
构建交付应用
→ 启动该应用
→ 操作真实窗口、请求该应用提供的 HTTP 接口
→ 退出该应用
→ 检查进程结束、端口释放
```

- API 验收先启动 Wails 应用，再请求它实际提供的接口。
- 生命周期验收操作真实应用的启动、关闭和退出。
- 不独立启动内部 HTTP 组件代替应用 E2E，不构建或测试 CLI。
- 不引入测试专用应用入口、替代后端或仅供测试的服务启动模式。
- 开发热更新检查不能替代正式构建产物验收。HTTP 请求可以作为 E2E 检查手段，但服务必须来自正在运行的真实应用。
- 第一阶段进行必要的真实应用验证，不提前扩展建设完整 E2E 框架。

第一步验收项：

1. 开发命令可以打开 Wails 窗口，前端热更新正常。
2. 正式构建产物能够独立启动并显示内嵌页面。
3. 通过正在运行的 Wails 应用验证健康、版本和 OpenAPI 接口保持兼容。
4. 默认地址生效；通过应用配置验证监听地址、端口和动态端口，界面能够显示实际地址。
5. 端口占用时启动明确失败，初始化失败不遗留服务。
6. 关闭应用后进程结束、端口释放，再次启动成功。
7. 源码、构建配置及验证流程不存在独立 CLI 或 sidecar 入口，`packages/desktop` 无改动。

## 实施顺序与后续边界

1. 迁移 Go module 和现有业务代码，删除 CLI 入口与 Cobra，清理旧构建链路。
2. 建立 React/Vite 前端、Wails 唯一入口、开发任务及平台构建配置。
3. 接入现有 HTTP 服务，落实启动错误传播、运行异常处理和统一退出清理。
4. 接入前端状态展示，完成正式前端构建和资源嵌入。
5. 构建并运行真实 Wails 应用，按上述验收项核验结果。

托盘、隐藏窗口继续运行、单实例、开机启动、自动更新及流式接口留到后续步骤。第一步完成标准是一个能够真实构建、启动、提供现有 API 并正常退出的 Wails 应用。

## 实施结果

已按上述方案完成实施，交付物为 `packages/hub` 下的 Wails v3 同进程应用。

关键实现：

- Go module 迁入 `src-wails`，module 名定为 `backstage.lisitede.com/hub`（Go 域名正序惯例；`packages/gitea-cli` 的 `com.lisitede.backstage.gitea` 属反向 DNS 误用，后续应另行迁移）；删除 `cmd/`、`internal/cmd/`、Cobra 依赖与 `build-mac.sh`；wails `v3.0.0-beta.32` 通过 go.mod `tool` 指令锁定，统一使用 `go tool wails3`。
- 命名规范：Go module 用域名正序（`backstage.lisitede.com/hub`）；原生包标识符（CFBundleIdentifier、Windows manifest）用反向 DNS（`com.lisitede.backstage.hub`）；npm 包用 `@head/backstage-hub`。Wails 绑定目录 `src/bindings/backstage.lisitede.com/hub/` 为 module 路径镜像，由生成器维护。
- 依赖版本对齐：huma v2.39.1 要求 uax29 v2.7.0，与 wails 工具链的 displaywidth v0.7.0 不兼容，huma 降为 v2.36.0（所用接口行为不变），uax29 对齐 wails 要求的 v2.4.0。
- `internal/server` 改为可管理生命周期：`Start` 同步监听后在 goroutine 提供服务，`Addr` 返回实际地址，`Err`/`Failure` 暴露意外终止，`Stop` 幂等优雅关闭（独立超时 context，超时强制关闭）。
- `src-wails/main.go` 为唯一入口：读取配置（`HUB_HOST`/`HUB_PORT` 环境变量，默认 `127.0.0.1:20080`）→ 同步启动 HTTP 服务 → 创建窗口 → 事件循环 → 退出统一清理；监听失败直接报错退出且不创建窗口；服务意外终止记录日志并结束应用。
- 前端 React + Vite 最小界面，通过 Wails 服务绑定 `GetStatus` 显示版本、实际 HTTP 地址与运行状态；生成绑定提交于 `src/bindings`。
- 正式构建链：`go tool wails3 build` 依次执行依赖整理、绑定生成、`vite build` 输出 `src-wails/frontend/dist`（go:embed 嵌入）与平台编译；`go tool wails3 task package` 产出 `.app`。开发构建通过 status 条件跳过前端重建。
- 仓库惯例为各包独立锁文件，`packages/hub/pnpm-lock.yaml` 依此惯例生成；未改动 `packages/desktop`。

已验证（macOS arm64 实测）：

1. `go tool wails3 build` 与 `task package` 成功，产物 `src-wails/bin/backstage-hub.app` 可独立启动。
2. 真实应用启动后 `/api/ping`、`/api/v1/version`、`/api/openapi.json`、`/api/openapi.yaml` 响应正常，未注册路径返回 404。
3. 默认地址 `127.0.0.1:20080` 生效；`HUB_PORT=0` 动态端口生效（实测分配 49906），日志与绑定均报告实际地址。
4. 端口占用时第二个实例启动明确报错退出；进程终止后端口释放，可再次启动。
5. `pnpm --filter @head/backstage-hub dev` 完整链路通过：绑定生成 → DEV 构建 → Vite（5173）→ 窗口启动 → API 可用；中断后进程树与端口全部清理。

人工验收项：

- 窗口界面实际显示（版本、地址、运行状态）与前端热更新效果。
- 关闭主窗口即退出应用的交互行为。
- Windows 平台构建：配置已建立（`build/windows`），未在 Windows 环境实测。
