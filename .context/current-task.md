# VISIT-026：Whiteboard 通过原生 HTTP API 读取 Markdown 描述

> 状态：HTTP 改造已实施；从宿主配置读取 PublicURL 的修正已实施，重新安装及运行验证待完成。
> 目标宿主：原版 Vikunja v2.7.1 + Yaegi，不 fork、不修改宿主、不扩展 symbols。
> 插件源码：`packages/vikunja-whiteboard-plugin/01_http.go`、`packages/vikunja-whiteboard-plugin/main.go`。
> 对外接口：`GET /api/v1/plugins/whiteboard/tasks/:taskId/description`。
> 内部调用：`GET /api/v2/tasks/{taskId}?format=markdown`。
> 地址配置：使用宿主配置 `config.ServicePublicURL.GetString()`。
> 认证：仅透传调用方的 Authorization，插件没有自己的 Token。
> 当前不修改宿主配置、不安装、不重启、不创建 Token；代码变更需用户重新部署后验证。

## 1. 已确认的需求

1. 保留现有 Whiteboard Ping 和 description 路由。
2. description 改用 HTTP 调用原生 v2 任务详情接口，固定请求 Markdown。
3. 插件不再直接调用 `Task.ReadOne()`，不自行转换 HTML。
4. 对外只提供 Markdown，不接受格式、筛选或其他查询参数；保留必需的路径参数 `taskId`。
5. 服务根地址直接读取 Vikunja 内置配置 `config.ServicePublicURL.GetString()`，不依赖插件环境变量。
6. 认证只使用当前请求透传的 Authorization，不读取 `BACKSTAGE_VIKUNJA_TOKEN`，不使用内置、全局或备用 Token。
7. HTTP 基础能力拆到 `01_http.go`；入口仍叫 `main.go`，不加数字前缀。
8. 不 fork Vikunja，不修改宿主源码或 Yaegi symbols。

本方案替代原先“在插件中使用数据库会话、当前用户、CanRead 和 ReadOne 读取 HTML”的方案。

## 2. 已核对的宿主能力

### 2.1 原生 Markdown 接口

Vikunja v2.7.1 支持：

```http
GET /api/v2/tasks/{taskId}?format=markdown
Authorization: Bearer <调用方凭据>
Accept: application/json
```

原生处理链：

```text
任务详情 HTTP handler
→ 校验认证及 tasks.read_one 路由权限
→ handler.DoReadOne：任务读取权限检查及数据加载
→ convertTasksToMarkdown
→ 返回任务 JSON，description 为 Markdown
```

存储仍为 HTML。Markdown 转换发生在 v2 API 层，`Task.ReadOne()` 本身不支持格式参数。

`pkg/richtext` 提供 HTML/Markdown 转换函数，但 v2.7.1 未将该包注册到 Yaegi 的宿主符号表。插件不直接 import 该包，不复制转换代码，使用原生 HTTP API 复用能力。

原生转换函数在转换出错时可能保留 HTML；插件固定请求原生 Markdown 表示，不额外实现 HTML 转换或自行回退到 HTML 请求。该宿主行为作为原生 API 的限制记录，不承诺 Markdown 的无损往返。

### 2.2 已有实测

2026-11-02，原版 Vikunja v2.7.1：

- 原生 `/api/v2/tasks/6?format=markdown` 返回 HTTP 200，描述为 `# 设计一个数学学习路径/路线图`。
- 原生同一任务 `format=html` 返回 HTML 描述。
- 原生 `format=text` 返回 HTTP 422，只支持 `html` 和 `markdown`。
- 空描述任务通过原生 Markdown 接口返回空字符串。
- 当前插件 description 仍返回 HTML；给插件添加 `format=markdown` 不会自动应用原生转换。
- 当前插件 description 的正常读取、空描述、非法 ID、不存在任务及无效认证已验证；这些结果不代表本次 HTTP 改造已运行通过。

## 3. 对外接口约定

### 3.1 路由保持不变

源码注册相对路径：

```go
g.GET("/whiteboard/ping", handlePing)
g.GET("/whiteboard/tasks/:taskId/description", handleTaskDescription)
```

宿主添加 `/api/v1/plugins` 前缀。description 完整路径为：

```http
GET /api/v1/plugins/whiteboard/tasks/:taskId/description
```

不新增格式协商接口，不新增旧路径兼容路由。

### 3.2 输入限制

- `taskId` 使用 `strconv.ParseInt(..., 10, 64)` 解析，必须大于 0。
- 非数字、溢出、0、负数返回 HTTP 400。
- 不接受任何查询参数，检测到非空 `URL.RawQuery` 时返回 HTTP 400；不只检查 `format`，也不能忽略未知参数。
- 不转发调用方的查询字符串、请求体、格式请求头或目标地址。
- 内部路径使用解析后的数值 ID 构建，避免把原始路径片段直接拼接到上游 URL。
- Authorization 由宿主认证后透传；防御性检查缺失时返回 HTTP 401，不尝试其他凭据。

### 3.3 成功响应

继续返回 HTTP 200、JSON，只包含两个字段：

```json
{
  "task_id": 6,
  "description": "# 设计一个数学学习路径/路线图"
}
```

- `task_id` 来自上游任务 `id`，必须与请求的数值 ID 一致。
- `description` 来自原生 Markdown 响应，不添加 HTML 标签或进行二次转换。
- 空描述返回 `""`。
- “只提供 Markdown”指 description 内容的固定表示；不改成裸字符串、`text/plain` 或 `text/markdown` 响应。
- 不返回完整 Task，不增加 title、code、message、data 或格式字段。

## 4. 文件布局与加载顺序

最终仓库结构：

```text
packages/vikunja-whiteboard-plugin/
├── 01_http.go
└── main.go
```

两个文件都使用：

```go
//go:build ignore

package main
```

不写 `import "./01_http.go"`，同包代码直接调用。每个文件自行 import 它实际使用的标准库或已暴露宿主包。

### 4.1 显式顺序约定

- 固定文件名 `01_http.go` 和 `main.go`，不使用 `_http.go` 或 `02_main.go`。
- HTTP 辅助文件用数字前缀明确排在入口前，不依赖业务名称偶然排序。
- 底层机制仍是宿主 `os.ReadDir()` 的文件名排序，再逐个 `i.Eval()`；宿主没有加载清单配置。
- `01_http.go` 不得引用 `main.go` 中定义的插件类型、变量、错误或函数。
- `main.go` 可以引用 `01_http.go` 中已定义的辅助能力。
- 辅助文件中的声明也按依赖顺序组织；不通过导入子目录解决顺序问题。
- Yaegi 根据 `main.NewPlugin` 等工厂查找入口，不依赖文件名，但本项目保留用户要求的 `main.go`。

### 4.2 01_http.go 的职责

封装业务无关的 HTTP 基础能力：

- 读取、校验服务根地址并构建 `/api/v2` 基础 URL。
- 构造可复用 HTTP client，统一超时及重定向策略。
- 使用传入的 context、Authorization 和内部相对路径发起 GET。
- 读取并关闭响应体，提供状态码、必要的响应类型和原始响应体。
- 区分连接/读取失败和超时。

该文件不持有 Token，不知道 Whiteboard 插件类型，不处理任务领域逻辑。请求的认证必须作为每次调用的参数传入，不能保存到共享对象上。

可采用如下接口形态，具体命名可在实施时调整：

```go
type vikunjaHTTPClient struct {
    baseURL string
    client  *http.Client
}

func newVikunjaHTTPClient() (*vikunjaHTTPClient, error)

func (a *vikunjaHTTPClient) get(
    ctx context.Context,
    path string,
    authorization string,
) (statusCode int, contentType string, body []byte, err error)
```

客户端接受的 path 只能由插件内部构造，不允许从外部请求指定任意地址。

### 4.3 main.go 的职责

- 保留 WhiteboardPlugin、Name、Version、生命周期、singleton 和两个工厂。
- Name 继续返回 `whiteboard`，Version 保持 `0.1.0`。
- 保留 Ping 原有路由及 `{"code":0,"message":"ok"}` 响应。
- 在 Init 中调用辅助构造函数并保存客户端，正常初始化仍记录日志。
- 配置无效时 Init 返回明确错误；考虑宿主可能只记录 Init 错误而继续运行，description handler 必须对客户端未初始化返回 HTTP 503，避免 nil panic 或隐式 fallback。
- 参数及查询字符串校验、提取调用方 Authorization。
- 固定调用 `/tasks/{数值ID}?format=markdown`。
- 将上游成功 JSON 解码为最小结构，校验并裁剪响应。
- 处理或透传错误响应，不暴露完整上游任务内容。

删除 `db`、`models`、`user` imports 和以下调用：

```text
创建数据库 Session
GetCurrentUserFromDB
Task.CanRead
Task.ReadOne
Commit
```

任务权限由携带同一凭据的原生 v2 HTTP 接口检查，插件不重复执行数据库读取或权限检查。

## 5. 宿主配置、HTTP 与认证

### 5.1 地址配置

使用 Yaegi 已注册的 Vikunja 配置包：

```go
import "code.vikunja.io/api/pkg/config"

configuredURL := strings.TrimSpace(config.ServicePublicURL.GetString())
```

v2.7.1 的 `pkg/yaegi_symbols/vikunja_config.go` 将 `config.ServicePublicURL` 注册到符号表，插件可以通过 `config` 包读取宿主当前有效的 `service.publicurl` 配置。该值由 Vikunja 自身配置系统加载，因此无需额外设置 `BACKSTAGE_VIKUNJA_URL`，也无需让服务进程注入插件专用 URL 环境变量。

仍需进行 URL 安全校验：

1. 必须为 HTTP/HTTPS 绝对 URL，Host 非空。
2. 不允许 query、fragment 或 userinfo 凭据。
3. 值应为 Vikunja 服务根 URL，不含 `/api/v1`、`/api/v2` 或插件路由。
4. 保留部署根路径并追加 `/api/v2`。
5. 不根据调用方 Host 或转发头推导上游，不回退到另一 URL。

`service.publicurl` 必须能由 Vikunja 服务进程访问；如配置为外部域名、反向代理或 TLS 地址，应确认服务本机可达且不会重定向到登录门户。客户端禁止跟随重定向。

### 5.2 凭据透传

- 仅从当前请求读取 `Authorization`，原样设置到内部请求。
- 不读取或使用任何 Token 环境变量；`BACKSTAGE_VIKUNJA_TOKEN` 即使存在也必须忽略。
- 不内置 Token、不保存 Token 到插件字段、不使用管理员或服务账号备用 Token。
- 不重新解析、签发或校验 JWT；认证仍由宿主完成。
- 不透传 Cookie、其他请求头或调用方的 `X-Vikunja-Format`；内部只设置必要的 Authorization 和 `Accept: application/json`。
- 不将凭据写入 URL、日志或错误消息。

### 5.3 HTTP 行为

- 复用 `http.Client`，Timeout 为 30 秒。
- `http.NewRequestWithContext(c.Request().Context(), ...)` 保留调用方取消信号。
- `CheckRedirect` 返回 `http.ErrUseLastResponse`，不跟随重定向；3xx 作为上游异常处理，不把 Location 传给调用方。
- 每次请求独立创建 Header，避免并发调用之间串用 Authorization。
- 及时关闭 response.Body。
- 不直接把 CLI 的错误处理照搬成统一字符串错误；保留上游 HTTP 状态和错误体的结构。

## 6. 权限要求

API Token 必须同时具有：

```json
{
  "plugins": ["whiteboard_tasks_description"],
  "tasks": ["read_one"]
}
```

两次请求分别校验：

1. 外层插件接口需要 `plugins.whiteboard_tasks_description`。
2. 内层原生任务接口需要 `tasks.read_one`，再检查用户对目标任务/项目的读取权限。

任务不存在、用户无权读取、Token 缺少任一权限均不能通过重试、数据库直读或备用凭据绕过。

Ping 权限仍为 `plugins.whiteboard_ping`，与 description 独立。路由本身未改变，因此不新增插件权限名称；如果旧 Token 缺少 `tasks.read_one`，需由用户重新创建合适权限的 Token。

## 7. 返回与错误处理

| 场景 | 行为 |
| --- | --- |
| taskId 非法 | HTTP 400，不发起上游请求 |
| 传入查询参数（包括 format） | HTTP 400，不发起上游请求 |
| 无有效认证或缺少插件权限 | 外层宿主认证中间件拒绝 |
| handler 中缺少 Authorization | HTTP 401，无其他凭据 fallback |
| 上游客户端未初始化 | HTTP 503 |
| 上游 HTTP 200，任务 JSON 有效 | 返回最小 JSON，description 为原生 Markdown 表示 |
| 上游描述为空字符串 | HTTP 200，保留空字符串 |
| 上游 4xx/5xx | 保留上游状态码及错误响应体，不转成成功响应 |
| 上游 3xx | HTTP 502，不跟随或转发重定向 |
| 建连、响应体读取等失败 | HTTP 502，使用不包含凭据的错误消息 |
| HTTP 请求或读取超时 | HTTP 504 |
| 上游成功响应不是有效任务 JSON、id 不匹配或 description 类型错误/缺失 | HTTP 502 |

成功响应解码时保留对必需字段存在性的检查：允许 description 为空字符串，不把缺失、null 或错误类型静默视为空描述。JSON 字段顺序不作为验证条件。

错误响应仅透传必要的 Content-Type 和响应体，不复制所有上游响应头，不泄露 Set-Cookie、Location 或其他会话信息。

原生错误 401/403/404 不改写成空 description。不给未知格式或不支持 Markdown 的旧宿主做 HTML 请求 fallback。

## 8. 变更范围与实施顺序

允许改动：

```text
新增 packages/vikunja-whiteboard-plugin/01_http.go
修改 packages/vikunja-whiteboard-plugin/main.go
更新 .context/current-task.md 的实际进度及验证记录
```

实施顺序：

1. 检查已有未提交差异，保留用户现有代码和文档范围之外的改动。
2. 新增 `01_http.go`，实现地址校验、共享 client、请求和响应基础能力。
3. main 中初始化客户端，将 description handler 改为上述 HTTP 链路。
4. 删除 description 的数据库读取和相关 imports。
5. 保留 Ping、路由、名称、版本、工厂和 singleton。
6. 执行格式、差异、配置和认证检查，记录尚未运行验证的部分。

不修改 CLI 或直接 import CLI internal 包；仅参考提交中的配置方式。不新增 go.mod、README、安装脚本、业务写接口或其他模块文件。

## 9. 验证计划

### 9.1 实施时检查

- 对两个 Go 文件运行 gofmt，执行 git diff --check。
- 文件名、package、build tag、双工厂和根目录布局符合约定。
- 确认 `01_http.go` 不依赖 main 中的定义。
- 确认通过 `config.ServicePublicURL.GetString()` 读取宿主 URL，不读取任何 Token 环境变量。
- 确认每次调用透传当前 Authorization，不存储或跨请求复用凭据。
- 确认固定 format=markdown，拒绝所有外部查询参数。
- 确认不再存在数据库 Session、CanRead 或 ReadOne 调用。
- 使用隔离的 HTTP stub 验证内部请求路径、认证透传、固定参数、错误状态、超时及重定向拒绝；验证空描述和异常 JSON，不依赖真实任务写入。
- 如进行 Yaegi 兼容验证，按实际宿主方式先 Eval 01_http.go 再 Eval main.go；普通 go build 不能替代该检查。

### 9.2 部署与真实验证

安装、服务配置和重启由用户执行，未经后续请求不自动操作本地服务。

最终安装布局：

```text
<plugins.dir>/
└── whiteboard/
    ├── 01_http.go
    └── main.go
```

必须复制两个文件并保留文件名，不额外嵌套仓库目录，不遗留其他版本的重复 .go 文件。

插件重新部署并重启后：

1. 确认 `service.publicurl` 被读取，无配置、解释器、重复定义或符号加载错误。
2. 用同时具有插件权限和 tasks.read_one 的调用方 Token 请求 description。
3. 与原生 `/api/v2/tasks/{taskId}?format=markdown` 对比，ID 与 Markdown 必须一致。
4. 补测空描述、非法 ID、任何查询参数、不存在任务。
5. 补测无 Token、无效 Token、缺少插件权限、缺少 tasks.read_one、无任务读取权限。
6. 验证 Ping 仍正常，插件权限映射未改名。
7. 不创建或修改真实任务来完成验证；测试凭据不足的场景标为未验证。

## 10. 当前状态与非目标

- description HTTP 转调和 `01_http.go` 已实施，包含调用方 Authorization 透传及固定 Markdown 查询。
- 插件地址配置已改为 `config.ServicePublicURL.GetString()`，移除 `BACKSTAGE_VIKUNJA_URL` 依赖。
- 重装后实测曾因 URL 环境变量缺失返回 503；该失败发生在改用宿主配置之前，不能代表当前代码结果。
- 用户需重新部署、重启后验证插件加载和 description 转调。
- 原生 Markdown 读取已实测；插件封装后的 Markdown 读取尚未验证。
- 不 fork Vikunja、不修改宿主、不扩展 symbols。
- 不存储 Token、不使用服务账号替代调用方。
- 不新增 HTML/Markdown 转换器或格式切换功能。
- 不改 UI 保存方式、数据库结构、任务数据、MCP 或 Whiteboard 业务能力。
- 不要求为插件单独配置环境变量；不操作用户服务配置、安装目录或 Token 生命周期，除非后续用户明确要求。

## 11. 调查依据

- [Yaegi 中注册的 Vikunja 配置符号](https://github.com/go-vikunja/vikunja/blob/v2.7.1/pkg/yaegi_symbols/vikunja_config.go)，包括 `ServicePublicURL`。
- `BACKSTAGE_VIKUNJA_URL` 与 HTTP Adapter 设计曾参考提交 `5dfba51b22f102c88b12d09b21f830ea8b36f425` 的 `packages/vikunja-cli/internal/vikunja/adapter.go`；当前插件不再读取此环境变量。
- [v2.7.1 原生任务 HTTP handler](https://github.com/go-vikunja/vikunja/blob/v2.7.1/pkg/routes/api/v2/tasks.go)
- [v2.7.1 Markdown 格式协商与转换调用](https://github.com/go-vikunja/vikunja/blob/v2.7.1/pkg/routes/api/v2/richtext.go)
- [v2.7.1 HTMLToMarkdown 实现](https://github.com/go-vikunja/vikunja/blob/v2.7.1/pkg/richtext/htmltomarkdown.go)
- [v2.7.1 Yaegi 多文件加载和工厂查找](https://github.com/go-vikunja/vikunja/blob/v2.7.1/pkg/plugins/yaegi/loader.go)
- [v2.7.1 宿主符号注册目录](https://github.com/go-vikunja/vikunja/tree/v2.7.1/pkg/yaegi_symbols)
