# Hub Vault 凭据代理实施方案

> updated_by: Qoder - Claude (Qoder)
> updated_at: 2026-11-19 18:20:53

## 目标

在 packages/hub 内新增 `internal/vault`，实现**本地凭据代理 + 内存审计日志**：外部调用方持 secret 通过 HTTP 从 macOS keychain 取值，hub 首页负责发放 secret 并展示审计记录。本文覆盖此前的 Hub HTTP MCP 概念验证方案。

## 已锁定的需求边界（实施时不得扩大）

以下均为 2026-11-19 用户澄清的定稿，实施 Agent 不得重新讨论或"加固"：

| 项 | 定稿 |
| --- | --- |
| 对外接口 | 仅 `GET /vault-api/get?role=&secret=&n=&k=` 一个 HTTP 接口，非 MCP，不走 huma |
| 验证 | 只校验 secret 是否为 Vault 的 key；role 不校验，仅作审计标签 |
| 万能钥匙 | secret 持有者可通过 n（service）、k（key）读取任意 keychain 条目，**这是需求本身，不是缺陷** |
| 错误契约 | 成功 `{"value": ...}`；**其余一律 401**，响应体固定 `{"error": "unauthorized"}`；原因区分只写日志（secret miss / keyring not found / keyring 其他错误），不透出到响应 |
| secret | 6 位随机数字，`crypto/rand` 生成（`[0, 1e6)` 前导零补齐），可多个并存；key 撞车则重新生成 |
| 生命周期 | 纯内存，无过期、无撤销（撤销留待未来，本期不做）、进程重启即清 |
| 审计记录 | 验证通过且 keychain 取值成功后 append `{role, n, k, timestamp}`；失败的请求不落记录 |
| 前端取数 | 前端与 Go 同进程，**通过 Wails binding 直读内存**，不新增任何 HTTP 接口 |

## 架构

hub 首页是 Wails webview 加载的 Vite 前端（`packages/hub/src/App.jsx`），与 Go 进程同应用。数据面仅两条通道：

```
外部调用方 ──HTTP──▶ GET /vault-api/get ──▶ Store（mutex + map）
                                             ▲
Wails webview ──binding直调──▶ VaultService ─┘
                              (AddSecret / ListRecords)
```

- 加号按钮：前端调 binding `AddSecret()`，Go 侧生成 6 位 secret 插入 map 并返回前端展示
- Table：前端调 binding `ListRecords()`，拍平遍历所有 secret 的记录，得到 `{role, secret, n, k, timestamp}[]`，按 timestamp 倒序，2s setInterval 轮询
- 挂载：现有 mux 监听 `127.0.0.1:20080`，vault handler 以裸 `http.Handler` 注册

## 模块设计

```
internal/vault/
  vault.go    — type Store struct { mu sync.Mutex; entries map[string][]Record }
                type Record struct { Role, Secret, N, K string; Timestamp time.Time }
                New() *Store；AddSecret() string；Append(secret, role, n, k)；All() []Record
  handler.go  — func (s *Store) Handler() http.Handler
server.go     — Start() 中挂载一行：mux.Handle("GET /vault-api/get", store.Handler())
main.go       — application.Services 追加 application.NewService(store)，暴露 VaultService
src/App.jsx   — 首页改造：左上 "Vault"、右上「+」按钮、横线、下方 Table
```

binding 生成沿用 `src/bindings` 现有流程（dev 构建自动生成），前端仿照 `GetStatus` 的 import 方式新增 `AddSecret` / `ListRecords`。

## 并发设计

一把 `sync.Mutex` 足够：锁内原子化「查 secret → Append」；`keyring.Get` 放锁外执行（瓶颈是 keychain IPC 而非内存结构），避免外部调用方互相串行。不需要 RWMutex / 分片 / atomic map。

## 接口契约明细

| 情况 | 状态码 | 响应体 | 日志 |
| --- | --- | --- | --- |
| secret 命中且 keychain 取到值 | 200 | `{"value": "<值>"}` | — |
| secret 未命中 / query 缺参 | 401 | `{"error": "unauthorized"}` | `secret miss` / `missing param` |
| keychain `ErrNotFound` | 401 | `{"error": "unauthorized"}` | `keyring not found: n=... k=...` |
| keychain 其他错误（授权拒绝等） | 401 | `{"error": "unauthorized"}` | 完整错误文本 |

其他实施默认：不加 CORS（调用方为本地进程）；n/k 按标准 URL query 解码。

## 实施顺序

```
① internal/vault/vault.go（Store + Record + AddSecret/Append/All）
② internal/vault/handler.go（/vault-api/get，200/401 契约）
③ server.go 挂载一行 + main.go 注册 VaultService
④ wails3 生成 bindings
⑤ App.jsx 首页改造（Vault 标题 / 「+」按钮 / 横线 / Table）
⑥ pnpm run build 验证编译产物
```

## 背景与既有事实（供实施 Agent 参考）

- go-keyring v0.2.6 已在依赖中，现有读取方为 `internal/config/config.go` 的 `loadCredential`，其 fail-fast 行为正确，**本期不改它**
- keychain Service 名 `com.lisitede.backstage.hub` 与 Wails productIdentifier 一致；vault 的 n/k 由调用方指定，不限于该 Service
- App.jsx 现有首页展示 `HubService.GetStatus` 的状态面板（版本/地址/运行状态），本期改造为 Vault 页，状态面板去留由实施时按视觉权衡
