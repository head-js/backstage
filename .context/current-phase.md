---

## PHASE-200: 注入脚本实现

本 Phase 聚焦于基于选型结果实现可运行的注入脚本。

### Relevant Requirements

从 `.context/current-plan.md` 提取与 PHASE-200（注入脚本实现）相关的需求：

- **FR-001**：系统应通过 backstage 扩展（content_script 注入 `js/backstage.js`）在私有化 Gitlab 页面注入并执行 JS 脚本（原 Tampermonkey 语义由 backstage 扩展取代，见 Design 注入路线）。
- **FR-002**：系统应将读取到的 merge_requests 规范化为 JSON 并在浏览器 console 打印。
- **FR-003**：系统应使用 Gitlab 标准 API 字段，不得自定义或重命名字段。
- **FR-010**：当用户在 Gitlab 页面触发脚本（或脚本按匹配规则自动执行）时，系统应发起 merge_requests 读取请求。
- **FR-011**：当鉴权凭证有效且目标项目存在 merge_requests 时，系统应获取数据并打印 JSON。
- **FR-012**：当鉴权凭证缺失或失效时，系统应在 console 输出可识别的错误提示，且不得在提示中暴露凭证明文。
- **FR-020**：在脚本执行期间，系统应保持只读，不得发起任何写操作。
- **FR-030**：如果读取目标项目无 merge_requests 数据，则系统不得抛出未捕获异常，应在 console 输出空结果提示。
- **FR-031**：系统不得执行任何写回操作（创建/修改/删除 issues、MR、评论等）。
- **C-001**：只读，禁止任何写回操作。
- **C-003**：鉴权凭证（token / cookie）不得在 console 明文打印或写入日志。
- **C-004**：脚本以原生 JS 实现，不引入外部依赖。
- **D-001~D-003**：私有化 Gitlab 访问权限、backstage 扩展已安装、至少一个含 MR 数据的项目（验证时由人工指定，PHASE-100 已用 project 133 实证）。

### Relevant Specs

从 `.context/current-plan.md` 提取与 PHASE-200 相关的 Specs（选型结论已由 PHASE-100 回填至 Design）：

- **SPEC-001 鉴权机制选型**：已选 cookie 路线（复用浏览器已登录会话，同源 fetch 自动携带 cookie）。
- **SPEC-002 接口选型**：已选公开 REST API（`GET /api/v4/projects/:id/merge_requests`）。
- **SPEC-003 注入方案**：已选 backstage 扩展（content_script），非 Tampermonkey。
- **SPEC-004 issues 数据规范化**：直接使用 Gitlab API 原始响应字段，不自定义/重命名，输出合法 JSON。

### Relevant Design

从 `.context/current-plan.md` 提取与 PHASE-200 实现直接相关的 Design 结论：

- **鉴权路线**：cookie 路线。复用浏览器已登录会话，同源 `fetch` 由浏览器自动携带 cookie 调 REST API v4；凭证不入脚本源码（C-003 内在满足）。失效兜底：会话过期（401）时 console 输出可识别提示（不含凭证明文），用户重新登录后恢复，无需改代码。
- **数据来源路线**：公开 REST API `GET /api/v4/projects/:id/merge_requests`。直接消费原始响应字段，不自定义/不重命名。路径编码约束：项目路径须 `encodeURIComponent` 单层编码为 `%2F`；未编码或双编码 `%252F` 均 404；数字 id 可用（PHASE-100 TASK-100 实证 project 133 MR 端点 200 + 有数据）。
- **注入路线**：backstage 扩展（`crx/manifest.json`，MV3）。`content_scripts` 在 `<all_urls>` 于 `document_end` 注入 `js/backstage.js` → Gitlab 页自动被注入；`route`/`verb` API 由 `frontstage.js`（web_accessible_resource）提供，暴露于 `window.backstage`。业务脚本以 backstage 站点脚本约定承载，但脚本本体不在本 backstage 仓库——归属独立的 `frontstage` 仓库（monorepo，每个 host 一个 `packages/<host>/` 包；参考 `frontstage` 仓库 `packages/zhihu.com/` 样板，用 `@head/frontstage-utils/ready` 的 `$ready` 包裹在 `READY` 后执行）。本仓库 `packages/frontstage/` 下 `www.example.com` 为旧规范遗留，不作为新增 host 落点。
- **数据目标语义纠错**：数据目标端点为 `merge_requests`（非 `/issues`）；「issues」为统称，具体读取对象为 merge_requests。
- **边界与错误策略**：只读边界（仅 GET）、凭证安全（不进 console/日志）、错误提示（凭证缺失/失效、无数据、接口不可达输出可识别提示而非抛未捕获异常）、数据模型（直接消费标准 API 响应字段，不二次建模）。

### Tasks Breakdown

<!-- // 本章节仅记录 Task 的拆分结果与关键步骤，不用于执行、追踪或进度管理；禁止使用进度标记、CheckBox 或任何其他进度相关标记。 -->

#### TASK-100: 实现鉴权装配与 merge_requests 读取请求
- **Do**: 在 `frontstage` 仓库新建 `packages/gitlab.com/` host 包（含 `package.json`/`webpack.config.js`/`src/gitlab.com/frontstage.js`，参考 `packages/zhihu.com/` 样板；`<gitlab-host>` 取 `gitlab.com` 作为通用能力载体，私有化域名 `@match`/`match` 由人工确认）中实现 merge_requests 读取函数：按 cookie 路线用同源 `fetch` 调 `GET /api/v4/projects/:id/merge_requests`（浏览器自动携带 cookie，凭证不入源码，满足 C-003）；项目路径用 `encodeURIComponent` 单层编码为 `%2F`（数字 id 可用，PHASE-100 已实证 project 133）；读取函数实现为原生 JS、不引入外部依赖（C-004）。产出物为读取函数源码，关键实现说明写入 `.context/current-task.md` 对应 Task 章节。 -> `frontstage 仓库 packages/gitlab.com/src/gitlab.com/frontstage.js`
- **Check**: 请求返回 200 且响应体含 merge_requests 数据
- **Check**: 全程仅 GET，无任何写操作；脚本源码不含 token/cookie 凭证明文
- **实现**: `frontstage` 仓库新建 `packages/gitlab.com/` host 包（`package.json` / `webpack.config.js` / `src/gitlab.com/frontstage.js`），仿 `packages/zhihu.com/` 样板；使用 `vanilla.js/fetch/get` 的 `$get` + `@head/frontstage-utils/ready` 的 `$ready`（C-004 由用户裁决推翻改用 vanilla.js）；注册 backstage 路由 `GET /repos/:owner/:project/merge-requests`，handler 内 `encodeURIComponent(\`${owner}/${project}\`)` 单层编码拼 Gitlab API 端点（`/api/v4/projects/<owner%2Fproject>/merge_requests`），`$get(endpoint, ctx.request.query)` 发请求（vanilla.js 自动序列化 query + `credentials:'same-origin'` + `JSON.parse`），`ctx.body = { status: 200, data }` 返回原始 MR 数据；仅注册路由不自动触发，由 `window.backstage.verb(...)` 调用；package.json dependencies 含 `vanilla.js@0.2.11` + `@head/frontstage-utils`；构建产出 `dist/gitlab.com/frontstage.js`（22.3 KiB，含内联 vanilla.js/get.js）；V1-V8 均通过（V1-V4 Agent，V5/V6/V8 人工端到端）。状态：【PASS】

#### TASK-110: 实现 merge_requests 数据规范化
- **Do**: 将读取响应以 Gitlab 标准 API 字段为数据源，用 jsonata 投影为 camelCase key 的对象并作为 `ctx.body.data` 返回；输出 JSON key 统一采用 **camelCase**（`url` / `author` / `sourceBranch` / `targetBranch` / `createdAt` / `updatedAt`），value 仍取自 Gitlab MR 标准字段（满足 FR-002；FR-003 / SPEC-004 的「不重命名」约束按本指令放宽为「key 统一 camelCase、value 取自标准字段」）；不二次建模、不下沉字段细节。**console 打印行保留为注释（调试噪声，按用户裁决不打印）**。产出物为规范化逻辑源码，写入 `.context/current-task.md` 对应 Task 章节。 -> `frontstage 仓库 packages/gitlab.com/src/gitlab.com/frontstage.js`
- **Check**: `ctx.body.data` 为合法 JSON 可序列化对象（可被 `JSON.stringify` / `JSON.parse` 往返）
- **Check**: 输出 key 统一为 camelCase，value 对照 Gitlab API 文档均取自标准字段
- **实现**: `packages/gitlab.com/src/gitlab.com/frontstage.js`——`$ready` 回调签名升级为 `async (Ajv, jsonata, rxjs) => {}`、deps `['ajv','jsonata','rxjs']`（jsonata 由 ready 注入，无 import/无新增依赖，C-004）；回调顶层 `const transformer = jsonata(expression)` 编译一次，表达式为模块顶层 `const expression`（camelCase key：`iid`/`title`/`state`/`url`(←web_url)/`author`(←author.name)/`sourceBranch`/`targetBranch`/`createdAt`/`updatedAt`，value 取自 Gitlab MR 标准字段）；handler 内 `$get`→`const picked = await transformer.evaluate(data)`→`ctx.body={status:200,data:picked}`；`console.log` 行保留为注释（用户裁决）；未引入 try/catch（错误处理归 TASK-120）。端到端实拨由用户 2026-07-03 通过（jsonata 注入缺陷已连带 `packages/utils/src/ready.js` 二次解包修复）。状态：【PASS】

#### TASK-120: 实现 dashboard 接口与 MR 列表 drawer 渲染
- **Do**: 在 `frontstage 仓库 packages/gitlab.com/src/gitlab.com/frontstage.js` 的 `$ready` 回调内，与 MR 路由并列新增注册 `backstage.route('GET', '/dashboard', handler)`；handler 内 `await window.backstage.verb('GET', '/repos/jingyingzi-shengsheng/mall-view-consumer/merge-requests', { state: 'opened' }, {})` 调同包已注册 MR 路由获取 camelCase MR 列表，用 `<table><tr><td>` 字符串模板拼装表格（表头 IID/Title/Author/Source Branch/Target Branch/Updated At，`url` 渲染为 `<a target="_blank">`，纯文本字段经 `escapeHtml` 转义），调用 `window.backstage.drawer(tableHtml)` 渲染；新增顶部 `TIME_PICTURE` 常量与 `escapeHtml` 纯函数，jsonata 表达式升级为数组化（尾部 `[]`）+ author 取括号前 + 时间格式化 `MM-DD HH:mm`。原「错误提示与边界处理」需求已撤销（由用户澄清重做），错误处理不在本 Task 范围。
- **Check**: 注册 `GET /dashboard` 后调用该路由能拉起 drawer 并展示 MR 列表表格；表格数据来自已注册 MR 路由的 camelCase 字段；全程仅 GET、无写操作、无凭证明文。
- **实现**: `packages/gitlab.com/src/gitlab.com/frontstage.js`——`$ready` 回调内新增 dashboard 路由（51-88 行），复用 MR 路由经 `window.backstage.verb` 取数，`data.map` 拼 `<tr>`，`escapeHtml` 转义纯文本防 drawer 渲染破坏，`window.backstage.drawer(tableHtml)` 拉起抽屉，`ctx.body={status:200,data:{count}}`；jsonata 表达式尾部加 `[]` 强制数组化绕过 singleton 展平，`$substringBefore($v.author.name,"（")` 取括号前姓名，`$fromMillis($toMillis(iso),'${TIME_PICTURE}')` 格式化时间。用户 2026-07-03 简单验收通过（drawer 拉起 + 表格渲染正常）。状态：【PASS】

#### TASK-130: 实现 frontstage 注入浮层按钮触发 dashboard（方案 B）
- **Do**: 在 `frontstage` 仓库 `packages/gitlab.com/src/gitlab.com/frontstage.js` 的 `$ready` 回调内，向 gitlab.com 页面注入一个 `position: fixed; right; bottom` 的全局浮层按钮（原生 JS，无外部依赖，满足 C-004；按钮样式从简）。按钮点击 `await window.backstage.verb('GET', '/dashboard')` 调**本页**已注册 dashboard 路由（TASK-120 实现），handler 同源 fetch 携带 Lax cookie 取 camelCase MR 列表，经 `window.backstage.drawer(tableHtml)` 渲染。按钮与业务同处 gitlab.com 主世界，本页 `verb` + 同源 cookie 天然契合（见 `docs/Chrome-Extension.md` §5，唯一既满足 C-003 又满足 C-004 的落点）；无需跨标签页寻址，直接复用现有 `drawer`/`verb` 能力，改动面最小。需处理 SPA 路由切换后按钮是否存活（`$ready` 一次性注入 vs 监听路由）。产出物为浮层按钮注入逻辑源码，关键实现说明写入 `.context/current-task.md` 对应 Task 章节。 -> `frontstage 仓库 packages/gitlab.com/src/gitlab.com/frontstage.js`
- **Check**: gitlab.com 页面右下角出现浮层按钮；点击按钮能拉起 drawer 并展示 MR 列表表格
- **Check**: 按钮点击调用 `window.backstage.verb('GET', '/dashboard')` 为本页路由直达；全程仅 GET、无写操作、无凭证明文外泄
- **人工验收**: UI 表现与交互可行性需人工在浏览器实测确认

#### TASK-140: 实现 crx options.html 扩展页触发 dashboard（方案 A）
- **Do**: 在 `crx/options.html`（`chrome-extension://<id>/options.html`）放置触发按钮（按钮样式从简），点击发起 dashboard 调用。已知约束（见 `docs/Chrome-Extension.md` §5/§6）：扩展页 origin 为 `chrome-extension://`，对 gitlab.com 跨源 fetch 受 `SameSite=Lax` cookie 限制不可靠；`manifest.json:47-53` 的 `content_scripts` 只 `matches: <all_urls>`，不匹配 `chrome-extension://`，扩展页内**无 `backstage.js` 内容脚本桥**，`backstage.invoke` 的 `window.postMessage`→BACKSTAGE 路是断的（无监听者）。可行路径：扩展页用 `chrome.runtime.sendMessage` 直连后台 → `background.js` 用 `chrome.tabs.query({url: 'gitlab.com'})` 寻址**已打开的 gitlab 标签页** → `chrome.tabs.sendMessage` 转发进去 → 该页 frontstage.js 的 route handler 执行同源 fetch 返回。即业务仍落在 gitlab 页，UI 与业务分离；且**强依赖 gitlab 标签页处于打开状态**，否则寻址失败。此路径走 `backstage.invoke('gitlab.com', 'GET', '/dashboard')`（跨 host），**不是** `window.backstage.verb('GET', '/dashboard')`（本页路由），语义需在实现中明确。产出物为 options 页触发逻辑与后台转发链路调整（如需），关键实现说明写入 `.context/current-task.md` 对应 Task 章节。 -> `crx/options.html` / `crx/js/options.js` / `crx/js/background.js`
- **Check**: options 页点击按钮能拉起 gitlab 页 drawer 并展示 MR 列表表格（需 gitlab 标签页已打开）
- **Check**: 调用链路为 `backstage.invoke` 跨 host 寻址，非本页 `verb`；全程仅 GET、无写操作、无凭证明文外泄
- **人工验收**: UI 表现与交互可行性需人工在浏览器实测确认

---
