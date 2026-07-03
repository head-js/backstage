# Chrome Extension 架构说明

> updated_by: HBR - GLM-5.2
> updated_at: 2026-07-04 00:07:00
> 适用范围：本仓库 `crx/`（`@head/backstage` 扩展，MV3）的注入与协调机制；业务脚本（host frontstage）归属独立的 `frontstage` 仓库。

## 1. 概述

backstage 扩展是一个**中枢协调器**：它让「无法直接连接的世界」——多个浏览器标签页、扩展后台（service worker）、扩展页（options）——能够互相发起调用。其设计准则：

- **中枢不执行业务**：后台/扩展页只做寻址与转发，业务逻辑落在目标 host 页面内执行。
- **同源优先**：对目标站点的 API 调用，由注入到该站点页面的脚本着页面身份发起，天然复用浏览器登录态（cookie），凭证不入源码。
- **route / verb 模型**：页面侧用 `backstage.route` 注册路由，调用方用 `backstage.verb` / `backstage.invoke` 发起；跨世界投递由扩展桥接。

## 2. 三个世界与注入物

| 世界 | 运行位置 | 注入物 | 来源 |
|---|---|---|---|
| 主世界（page main world） | 任意 http(s) 页面 | `frontstage.js` + 联邦 vendors，挂 `window.backstage` | `crx/js/backstage.js` 用 `<script>` 注入 |
| 隔离世界（content script） | 同上页面，但隔离 | `backstage.js`（content script） | `manifest.json` 的 `content_scripts` |
| 扩展后台 | service worker | `background.js` | `manifest.json` 的 `background` |
| 扩展页 | `chrome-extension://<id>/options.html` | `frontstage.js`（直接 `<script>` 引入）+ `options.js` | `options.html` |

关键事实：

- `crx/manifest.json:47-53`：`content_scripts` 在 `<all_urls>` 于 `document_end` 注入 `js/backstage.js`（外加 `css/shoelace-light.css`）。
- `crx/manifest.json:31-46`：`vendors/backstage-vendors.js`、`js/frontstage.js` 等列为 `web_accessible_resources`，供主世界 `<script>` 加载。
- `crx/js/backstage.js:6-21`：`requirejs(src)` 创建 `<script src=chrome.runtime.getURL(src)>` 注入主世界。
- `crx/js/backstage.js:159-166`：`initBackstage()` 顺序加载 `vendors/backstage-vendors.js` → `js/frontstage.js`。
- **不存在 `window.frontstage` 全局**：注入主世界的 bundle 叫 `frontstage.js`，但它挂到 `window` 上的变量是 `backstage`（`packages/backstage/src/frontstage.js:61-62`）。

## 3. 注入到页面主世界的东西：`window.backstage`

`packages/backstage/src/frontstage.js` 在主世界暴露的成员：

| 成员 | 含义 | 源码 |
|---|---|---|
| `backstage.route(method, path, handler)` | 注册路由（绑定 `router.verb`） | `:71` |
| `backstage.verb(method, path, search, form, headers)` | 发起本页路由调用（绑定 `client.verb`） | `:74` |
| `backstage.invoke(svc, method, ep, search, form)` | 调用**另一个 host 的服务**（经后台转发） | `:65-68` |
| `backstage.toast(msg, variant, duration)` | Shoelace alert 提示 | `:77-98` |
| `backstage.drawer(html, width)` | Shoelace drawer 抽屉 | `:101-124` |
| `backstage.onClick(handler)` | 在页面右下角注入一个原生 JS 圆形浮层按钮（`position: fixed; right; bottom`，深灰色近黑色），点击触发 `handler`；返回该按钮元素。无外部依赖（C-004）。每次调用新增一个独立按钮。 | `:127-178` |

> 注：`onClick` 为通用 UI 触发能力，与具体业务解耦；行为注册（如点击触发 `verb('GET','/dashboard')`）由各 host 业务脚本在其 `$ready` 回调内调用 `backstage.onClick(handler)` 完成（见 TASK-131）。

就绪信号：`packages/backstage/src/frontstage.js:127` `window.postMessage({ source: '@head/backstage', type: 'READY', from: 'FRONTSTAGE' })`。业务脚本（`frontstage` 仓库 `packages/utils/src/ready.js`）监听该 `READY` 后执行回调，回调内通过 `window.backstage` 访问扩展 API。

主世界与隔离世界的桥：`frontstage.js` 不直接使用 `chrome.*` API，而是 `window.postMessage` 投递 `{from:FRONTSTAGE, to:BACKSTAGE}`（`packages/backstage/src/frontstage.js:21-32`）；隔离世界里的 `backstage.js`（content script）监听该消息（`crx/js/backstage.js:103-131`），再经 `chrome.runtime.sendMessage` 与后台通信。

## 4. 完整跨世界调用链路

以「在 gitlab 页 console 调 `await window.backstage.verb('GET','/merge-requests/133',{}, {})`」为例（本页路由，最短路径）：

```
console ──verb──▶ frontstage.js(client.verb) ──▶ router 分发到 route handler
  handler 内同源 fetch('/api/v4/projects/133/merge_requests', {credentials:'same-origin'})
  ──cookie 自动携带──▶ Gitlab ──MR JSON──▶ ctx.body = {status, data}
```

以「跨 host 调用 `backstage.invoke('gitlab.com','GET','/merge-requests/133')`」（例如从别的标签页或扩展页发起）为例：

```
调用方 frontstage.invoke
  └─ window.postMessage {FRONTSTAGE→BACKSTAGE, call}        # 主世界→隔离世界
     └─ backstage.js(callBackground) chrome.runtime.sendMessage {BACKSTAGE→BACKGROUND}
        └─ background.js: chrome.tabs.query({url: call.svc})   # 找一个 gitlab.com 标签页
           └─ chrome.tabs.sendMessage(tabId, {BACKGROUND→BACKSTAGE, call})
              └─ 该 tab 的 backstage.js(callFrontstage) window.postMessage {BACKSTAGE→FRONTSTAGE}
                 └─ 该 tab 的 frontstage.js: client.verb → route handler → 同源 fetch
                    └─ 结果沿 CALLBACK 原路返回
```

关键源码：`crx/js/background.js:4-30`——后台用 `chrome.tabs.query({ url: call.svc })` 寻址匹配 host 的**已打开标签页**，再 `chrome.tabs.sendMessage` 转发进去。因此跨 host 调用**要求目标 host 有一个活着的标签页**，否则寻址失败。

## 5. 为何 gitlab 读取必须落在 gitlab.com 页面（cookie 同源论证）

cookie 路线（PLAN-101 鉴权选型结论）选择「同源页面 fetch」的根本原因：

- Gitlab 会话 cookie 通常为 `SameSite=Lax` + `HttpOnly`。
- **`SameSite=Lax` 的 cookie 不会随跨源 XHR/fetch 携带**（只在顶层导航时带）。因此从 `chrome-extension://`（options/background）或别的 host 页面直接 `fetch('https://gitlab.com/api/v4/...')`，即便扩展有 `host_permissions: <all_urls>`、即便 CORS 可被扩展绕过，cookie 也大概率不带 → 401。
- 而在 gitlab.com 页面主世界里 `fetch('/api/v4/...', { credentials: 'same-origin' })` 是**同源**，Lax cookie 自动携带，且凭证不进源码（C-003 内在满足）。
- 若改用 `chrome.cookies` API 读 cookie 再注入 `PRIVATE-TOKEN` 头，则退化回 token 路线，凭证进入代码/请求头，破坏 C-003。

结论：gitlab 的读取逻辑**只能**落在 gitlab.com 页面（注入式 frontstage route），由 background 做跨标签页/跨世界寻址转发。这是唯一既满足 C-003（凭证不入源码）又满足 C-004（原生 JS、无外部依赖）的落点。

## 6. `options.html` 的定位与限制

`crx/options.html` 是 `chrome-extension://<id>/options.html` 扩展页，用于扩展自身配置 UI（`app-account` / `app-options` umi 应用）。它**不直接处理 gitlab.com**：

- `crx/options.html:23` 虽然也加载 `frontstage.js` 拿到 `window.backstage`，但 `manifest.json:47-53` 的 `content_scripts` 只 `matches: <all_urls>`（http/https/file/ftp），**不匹配 `chrome-extension://`**——扩展页里**没有 `backstage.js` 内容脚本桥**。
- 因此 `options.html` 里 `backstage.invoke` 走的 `window.postMessage`→BACKSTAGE 这条路是断的（无监听者）。扩展页要与服务通信只能用 `chrome.runtime` 直连后台（特权页特权），是另一条独立路径。
- 即便直连后台，`background.js` 仍会把调用 `chrome.tabs.query({url: svc})` 转发进匹配的 gitlab 标签页执行，依旧不落在 options 侧。
- 且如第 5 节所述，options（`chrome-extension://` origin）对 gitlab.com 跨源 fetch 受 `SameSite=Lax` cookie 限制，不可靠。

## 7. 与 `frontstage` 仓库的关系

- **注入机制**（content_script、`route`/`verb`、后台桥接）由本 backstage 仓库提供，暴露于 `window.backstage`。
- **业务脚本**（各 host 的站点脚本）不在本仓库，归属独立的 `frontstage` 仓库（monorepo，每个 host 一个 `packages/<host>/` 包，如 `packages/zhihu.com/`、新增的 `packages/gitlab.com/`）。业务脚本用 `@head/frontstage-utils/ready` 的 `$ready` 包裹，在 `READY` 后注册路由。
- 本仓库 `packages/frontstage/` 下的 `www.example.com` 等为旧规范遗留，不作为新增 host 的落点。

## 8. 调试入口

- **页面侧调试**：在目标 host 页面（如 gitlab.com）console 直接用 `await window.backstage.verb('GET','/merge-requests/133', {}, {})`；查路由可用 `backstage.route`。
- **无独立 `window.frontstage` 调试壳**：如需更友好的调试别名（语法糖/REPL），属新增设计，当前仓库不存在。
- **后台调试**：service worker 控制台查看 `chrome.runtime.onMessage` 投递日志（`crx/js/background.js:39-55` 有 `console.debug`）。
- **就绪判断**：`window.backstage` 存在即代表 `frontstage.js` 已就绪；`$ready` 已保证回调在 `READY` 后执行。

## 9. 路由参数解码 footgun（实现注意）

`@head/edge` 的 `Layer.prototype.params` 对捕获值执行 `safeDecodeURIComponent`（源码佐证：`packages/backstage/node_modules/@head/edge/src/lib/router.js`，`safeDecodeURIComponent` 调用处）。因此 `ctx.params.id` **已被解码**；handler 内若要拼入含 `/` 的项目路径，必须对 `ctx.params.id` 再做一次 `encodeURIComponent`，以达到「单层编码」效果（`group%2Fsub`）。调用方契约：对 `group/sub` 类路径项目须预先传已单层编码的 `group%2Fsub`，数字 id（如 `133`）原样传即可。
