# Vikunja Whiteboard 插件

> updated_by: HBR - GPT-6
> updated_at: 2026-11-19 22:37:06

## ext 扩展字段

向 Vikunja 原生 `tasks` 表人工执行（纯线下，不进代码、不进 migration）：

```sql
ALTER TABLE tasks ADD COLUMN ext_iid TEXT DEFAULT NULL;     -- 已于 2026-11-16 在本机执行
ALTER TABLE tasks ADD COLUMN ext_session TEXT DEFAULT NULL; -- 已于 2026-11-16 在本机执行
```

本机数据库：brew 安装的 Vikunja，sqlite，位于 `/opt/homebrew/var/vikunja/vikunja.db`（配置见 `/opt/homebrew/etc/vikunja/config.yml` 的 `database.path`）。

- `ext_iid`：承载 VisitNode.Id，规范为 `大写(Workspace.Name)-任务Index`（例如 `MYPROJECT-12`，不零填充）。
- `ext_session`：承载 VisitNode.Session，任意字符串。

> 命名约定：标识符一律不使用全大写 `ID`。PascalCase 场景用 `Id`（如 Go 导出字段、对外 JSON），camelCase 场景用 `id`；`ext_iid`、`ext_session` 是数据库列名，不在此列。

- **允许为空，默认值 NULL**（SQL 中显式声明 `DEFAULT NULL`）。允许空 + 默认值是为了让 Vikunja 原有读写（不提供该列）完全不受干扰。

## 存储层设计（03_store.go）

数据访问走 Vikunja 暴露给 yaegi 插件的 `code.vikunja.io/api/pkg/db`（xorm session 原生 SQL）。按**行**而非按字段设计，以容纳未来的更多 ext 列：

- `taskExt` 行类型聚合全部 ext 列。新增字段 = `ALTER TABLE` + 结构体加字段 + SELECT 列表加列名。
- **读**：`listExt()` 一次查询取出所有任务的 ext 列，服务列表接口（避免 N+1）；`findTaskExtByIid()` 按 VisitNode Id 反查单行（原生 id + 全部 ext 列）。NULL 列一律读为 `""`。
- **写**：按字段独立 setter（`setTaskExtIid` / `setTaskExtSession`），各校验 `RowsAffected == 1`；与 update 路由的部分更新语义一一对应。

## ext_iid 读写语义

- **不变量**：节点只能经本插件创建，创建成功后立即 `UPDATE tasks SET ext_iid`，因此读取时该字段必然存在；缺失或为空即视为数据完整性被破坏，返回错误而非静默处理。
- **读取路径**：`resolveVisitNodeTask` 经 `findTaskExtByIid`（`WHERE ext_iid = ?`）反查任务原生主键 `tasks.id`，再经 REST `GET /tasks/{id}` 取完整任务，由 `task.project_id` 关联项目；Id 字符串永不解析。
- **软删除**：所有 ext 查询带 `deleted_at IS NULL` 过滤；软删除任务的 Id 随之失效（404），其 `ext_iid` 残留不影响新任务（Vikunja 的 `(project_id, index)` 唯一索引包含已删除行，index 不会复用）。
- **稳定性收益**：Workspace 改名后 VisitNode.Id 不再变化（旧 computed 方案会因 Name 改变而失效）。

## ext_session 读写语义

- **唯一数据源**：`VisitNode.Session` 只从 `tasks.ext_session` 读取，String 类型（非数组）；不从 description 提取，创建时也不写入 description。
- **缺失语义与 ext_iid 不同**：`ext_session` 为可选列，NULL / 缺失一律读为 `""`，不报错——历史任务该列均为 NULL，报错会破坏现有功能。
- **创建**：`VisitNodeCreateRequest.Session` 非空时才 `setTaskExtSession`；为空则不写，保持 NULL。
- **description 中的 `## Session` 段**（历史数据残留，或调用方自行放入 Description）属于普通数据：插件不提取、不过滤、不隐藏，原样透传，随任务生命周期自然消失。

## VisitNode.Type 与原生颜色映射

- **业务类型**：创建、单节点读取、列表和更新的响应始终包含 `VisitNode.Type`，取值为 `"Task"` 或 `"Vision"`。这里的 Task 是 Whiteboard 的节点类型，与作为底层存储实体的 Vikunja Task 概念独立。
- **转换边界**：REST 和 MCP 只使用 Type。插件复用 Vikunja 的 `hex_color`，负责颜色映射与规范化。
- **存储**：复用原生字段，不增加 `ext_type`、数据库 CHECK、触发器或迁移。Type 不进入 Description 或 Session。
- **创建**：Create 无 Type 入参。插件向原生 API 创建任务时省略颜色，`hex_color` 自然留空。构造响应时，插件通过同一套读取规则将空色映射为 `Type: "Task"`。
- **读取**：移除可选的 `#` 并转大写后，`FFFCC8` 映射为 `Vision`。空色、`FFFFFF` 和其他颜色（包括旧值 `FFD304`）均映射为 `Task`。读取仅转换返回值。
- **写入**：显式提交 `Type: "Task"` 时写入 `FFFFFF`，提交 `Type: "Vision"` 时写入 `FFFCC8`。写入值为大写且省略 `#`。显式写入白色会让深色主题下的卡片也显示白底，空色卡片则随主题显示。

## visit-nodes.update（PUT）

- **方法**：`PUT /whiteboard/visit-nodes.update`。
- **必填字段**：每次更新都必须提供 `Id` 和 `Type`。Type 只接受 `"Task"` 或 `"Vision"`。缺失、`null`、空字符串或其他值均返回 400。
- **可选字段**：省略 `Progress` 或 `Session` 时保留对应值。`Progress` 为 0 到 100 的整数，`Progress: 0` 会将进度设为零，`Session: ""` 会清空会话。字段校验在写入前完成。Id 未找到时返回 404。
- **原生部分更新**：插件通过原生批量接口的 `fields` 列表写入 Type 对应的颜色，以及本次提供的 Progress，Session 通过既有扩展字段存储层写入。
- **MCP**：`whiteboard_update_visit_node` 使用相同的必填字段、可选字段和校验规则。
