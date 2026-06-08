# tauri_edge 方案

> updated_by: HBR - GPT-5
> updated_at: 2026-06-08 19:29:34

## 背景

需要新增一个 Tauri 侧能力，将前端 JS 调用包装到内部固定 CLI 调用，并把 CLI 返回值转换后返回给 JS。

当前能力名为 `tauri_edge`。它不是通用命令执行器，而是一个固定业务入口。当前阶段已经完成 `GET /blames` 链路、Rust DTO 转换、CLI helper 拆分，以及前端 Blames 列表和详情抽屉。

## 前端调用方式

前端调用固定为：

```ts
invoke("tauri_edge", {
  method: "GET",
  pathname: "/blames",
  query: {},
})
```

约束：

- `query` 必传。
- 当前没有参数时传 `{}`。
- `method`、`pathname`、`query` 作为请求结构保留，方便后续扩展。

## 核心边界

- `tauri_edge` 不做通用 `run command`。
- JS 不允许传入 CLI 命令、子命令或任意 shell 字符串。
- `GET /blames` 到 CLI 的映射写死在 edge 内部。
- `backstage-gitea` 被视为系统命令，edge 不负责查找、安装或打包处理。
- edge 相关路由和业务映射保留在 `commands.rs`。
- CLI stdout/stderr/exit_code 的通用解析逻辑拆到 `src/cli.rs`。
- 当前不设计通用 envelope，`GET /blames` 直接返回数组。

## 三层结构

```text
JS
  invoke("tauri_edge", { method: "GET", pathname: "/blames", query: {} })
    |
    v
Rust command: tauri_edge
    |
    v
[Layer 1] Edge Router
  写死识别 GET /blames
    |
    v
[Layer 2] Command Adapter
  写死执行 backstage-gitea plan GET /blames
    |
    v
[Layer 3] Transformer
  解析 CLI 返回并转换为 JS 可消费数据
```

### Layer 1: Edge Router

职责：

- 接收 `method`、`pathname`、`query`。
- 当前只处理固定能力：`GET /blames`。
- 命中后选择内部动作 `ListBlames`。
- 不做动态路由，不做通用转发。

当前写死映射：

```text
method = GET
pathname = /blames
=> EdgeAction::ListBlames
```

### Layer 2: Command Adapter

职责：

- 根据 edge 内部动作执行固定系统命令。
- 当前只执行：

```text
backstage-gitea plan GET /blames
```

约束：

- 不接收 JS 传入的命令。
- 不做通用 shell runner。
- 不负责 `backstage-gitea` 的安装、查找或打包。

### Layer 3: Transformer

职责：

- 读取 CLI stdout。
- 解析 JSON。
- 转换为 JS 可消费的数据结构。

当前实现：

- 上游结构为 `GiteaBlame`。
- 顶层 `id`、`name`、`gitea` 为必填字段。
- `gitea.updatedAt` 使用 `#[serde(rename(deserialize = "updatedAt"))]` 显式标注为上游反序列化字段。
- 输出结构为 `EdgeBlame`，使用 `#[serde(rename_all(serialize = "camelCase"))]` 显式标注为给 JS 的序列化字段。
- 输出 `id` 为组合 id：`AppId-PlanId-BlameId`。
- 原上游 blame id 输出为 `blameId`。
- `appId` 来自 `gitea.owner`，`planId` 来自 `gitea.repo`。

## Rust 文件结构

```text
src/
  cli.rs
    parse_cli_output
    summarize_output

  cli_tests.rs
    CLI stdout/stderr/exit_code 解析测试

  commands.rs
    tauri_edge
    GET /blames edge router
    backstage-gitea command adapter
    GiteaBlame -> EdgeBlame transformer

  commands_tests.rs
    edge router 和 Blame DTO 映射测试
```

## 前端 Blames 页面

当前 Blames 页面调用：

```ts
listBlames()
```

`listBlames` 内部调用：

```ts
invoke("tauri_edge", {
  method: "GET",
  pathname: "/blames",
  query: {},
})
```

页面行为：

- Table 显示 `App Id`、`Plan Id`、`Name`、`Updated At`。
- Table 不显示 `Blame Id`。
- Row 可点击，也支持 Enter / Space 打开详情。
- 点击 Row 后从右侧打开 daisyUI `drawer-end` 抽屉。
- 抽屉宽度在桌面为 `70vw`。
- 抽屉标题显示 Blame Title。
- 抽屉内容显示 `Blame Id` 和 `Blame Context`。
- `drawer-side` 使用 `overflow-x-hidden`，避免关闭状态下抽屉平移到屏幕外导致页面横向滚动条。

## 已完成

- 新增 Tauri command：`tauri_edge`。
- 定义请求结构：`method`、`pathname`、`query`。
- edge router 写死处理 `GET /blames`。
- command adapter 写死执行 `backstage-gitea plan GET /blames`。
- CLI stdout 解析拆到 `src/cli.rs`。
- Blame DTO 转换为 JS 可消费结构。
- Rust 测试拆为 sibling test 文件。
- 前端 Blames 页面接入 `listBlames`。
- 前端 Blames 页面实现右侧详情抽屉。

## 暂不处理

- 不实现通用命令执行能力。
- 不扩展其他 pathname。
- 不处理 `query` 的业务语义，只保留结构并传入 `{}`。
- 不处理 `backstage-gitea` 的命令发现、安装或 sidecar 打包。
- 不设计通用 envelope，先以整体调通为目标。
- 暂不重构 transform / map 的命名和职责边界。
- 暂不强制 `gitea.owner`、`gitea.repo`、`gitea.body`、`gitea.updatedAt` 必填。
- 暂不处理 `summarize_output` 的 UTF-8 边界安全问题。
- 暂不处理 in-app browser 对 localhost 的阻塞问题；此前验证时 `127.0.0.1` 和 `localhost` 均返回 `net::ERR_BLOCKED_BY_CLIENT`。

## 验收标准

- JS 能调用 `tauri_edge`。
- Rust 能接收到 `{ method: "GET", pathname: "/blames", query: {} }`。
- edge 内部能命中写死的 `GET /blames` 能力。
- command adapter 能执行 `backstage-gitea plan GET /blames`。
- CLI stdout 能被解析。
- transformer 能返回 JS 可消费数据，其中 `id` 为 `AppId-PlanId-BlameId`，`blameId` 为上游 blame id。
- JS 不具备传入任意命令的能力。
- Blames 页面能展示列表。
- Blames 页面点击 row 能打开右侧详情抽屉。
