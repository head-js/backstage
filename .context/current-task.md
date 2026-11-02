# Whiteboard 插件：工作区列表接口

## 目标与约束

新增 `GET /api/v1/plugins/whiteboard/workspaces`，使用调用者身份读取全部真实项目，转换为 Workspace 数组后直接返回。

用户确定的实现要求：

1. 源码适配目标 Vikunja 宿主的 Yaegi 逐文件加载方式。
2. 不编写测试文件、不创建模拟宿主或隔离测试工程。功能验收必须安装到实际 Vikunja 后进行。
3. 按本地应用设计，不设置请求超时，不识别超时，不分类错误，不建立可恢复错误分支。
4. 现有错误直接返回，安装时出现的错误应定位并修复，不以额外错误处理框架掩盖问题。

## 接口契约

- 对外接口定义为无业务参数的固定 GET 路由；handler 不解析或主动检查 query/body，认证 Header 正常使用。
- 固定一次读取 100 条，无对外分页、搜索、筛选或排序参数。
- 内部使用 `/api/v2/projects?per_page=100`；安装时核对实际版本的数量参数及效果。
- 全量范围为当前调用者可见的项目，不使用 CLI 环境变量 Token 或管理员凭证。
- 成功直接返回 Workspace 数组，不封装 `{code, data}`，不输出分页元信息。
- 空结果返回 `[]`，不返回 `null`。
- 请求和 JSON 解码错误直接返回。

## Yaegi 加载规则

本地保存的宿主 `pkg/plugins/yaegi/loader.go` 用 `os.ReadDir` 按文件名读取全部 `.go` 文件，逐个调用 `i.Eval`，最后查找工厂。它不按普通 Go package 一次性编译，也不排除 `_test.go`。

- 前面的文件不能依赖后面才定义的类型、函数或变量。
- 跨文件重复导入的包使用不同的显式别名，避免解释器作用域内名称冲突。
- 插件目录只放正式源码，不放测试文件；build tag 不是让宿主跳过文件的机制。
- 源码保留 `package main` 和 `//go:build ignore` 约定。
- 不修改宿主加载器，不用伪造符号表代替安装验收。

## 文件组织

加载顺序：

```text
01_vikunja.go → 02_types.go → 03_services.go → main.go
```

| 文件 | 职责 |
| --- | --- |
| `01_vikunja.go` | 从 `01_http.go` 重命名；客户端配置、HTTP GET、固定 projects 请求。只依赖自身及宿主包，不引用后续类型 |
| `02_types.go` | Workspace、Vikunja Project、项目列表结构；不定义专用错误或错误响应结构 |
| `03_services.go` | 显式接收客户端，解码、过滤伪项目、转换 Workspace；不访问 main.go 的 singleton |
| `main.go` | 生命周期、工厂、单例、路由与 handler；直接返回调用错误 |

包别名：客户端沿用原包名；服务文件用 `servicecontext`、`servicejson`；入口用 `handlerjson`、`handlererrors`、`handlerstrconv`、`handlerstrings`。

## 数据转换

```go
type Workspace struct {
    ID          int64   `json:"id"`
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Priority    float64 `json:"priority"`
}
```

| 字段 | 映射 |
| --- | --- |
| `id` | Project.id |
| `title` | Project.title 原样复制 |
| `description` | Project.description 原样复制，保留空白和原格式 |
| `priority` | `200 * 65536 - Project.position`，float64，不截断或取整 |

保留已有 CLI 的 Workspace 语义：排除 id < 0 的伪项目，包含子项目，输出扁平数组并保留上游顺序。priority 越大越优先，可为负数。

上游列表只解码 `items` 为 `[]VikunjaProject`，不读取分页元信息，不对字段缺失、null 或列表完整性增加额外防御性检查。JSON 解码错误直接返回。

## 调用与错误传递

客户端方法：

```go
func (a *vikunjaHTTPClient) get(ctx context.Context, path, authorization string) ([]byte, error)

func (a *vikunjaHTTPClient) listProjects(ctx context.Context, authorization string) ([]byte, error)
```

- 创建请求、发请求、读正文发生错误时直接 `return nil, err`。
- 仅统一检查一次 HTTP 状态是否为 200；否则返回包含上游状态及原始正文的普通错误，不区分 3xx、4xx、5xx。
- 不配置 Timeout，不提供 isVikunjaHTTPTimeout，不使用 errors.Is/errors.As 分类。
- 保留请求上下文、调用者 Authorization 和现有禁止跟随重定向的行为。

服务方法：

```go
func listWorkspaces(
    client *vikunjaHTTPClient,
    ctx servicecontext.Context,
    authorization string,
) ([]Workspace, error)

func projectToWorkspace(project VikunjaProject) Workspace
```

- 调用方直接传入已初始化的客户端，不增加客户端为空的防御性分支。
- 获取及 JSON 解码错误直接返回，不包装成其他类别。
- 初始化非 nil 结果数组、过滤伪项目并转换。
- Service 不依赖 Echo，不读取 request，不写 response。

Handler：

- 注册 `/whiteboard/workspaces`，通过宿主前缀得到完整路径。
- 原样读取 Authorization 并传给上游，将客户端显式传给 service；不在插件中重复校验 token，也不读取 query 或 body。无效或缺少 token 的处理交给上游。
- 错误直接 `return err`，由宿主统一处理；不自行映射 400/401/502/503/504，不透传另一套错误响应。
- 正常结果 `return c.JSON(httpStatusOK, workspaces)`。
- 现有任务描述接口只保留路径解析、上游调用和响应转换，不增加 ID、query、token、客户端存在性或响应字段一致性防御分支。

## 安装后验收

本地只做源码审阅、gofmt 和 diff 检查，不编写或运行测试。

1. 确认实际安装目标、宿主版本和插件目录；安装四个正式 Go 文件，移除旧 `01_http.go` 及任何测试文件。
2. 按实际机制重新加载/重启，确认插件加载、初始化和路由注册成功。
3. 使用实际认证请求完整 workspaces 路径，核对数组响应、字段转换、真实项目集合与顺序。
4. 对照上游项目集合，确认固定数量 100 的参数有效、没有遗漏；安装中出现错误就直接定位并修复。
5. 检查原 ping、任务描述接口及空列表行为。

当前状态：源码已调整，待实际安装验收。未安装调用前不得宣称宿主加载或功能已验证通过。
