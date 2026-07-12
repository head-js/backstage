# Desktop 双窗口方案

> updated_by: HBR - GPT-5
> updated_at: 2026-07-12 16:42:38

## 目标

Desktop App 使用两个 Tauri native window 实现异质窗口效果：

- `menu` 是右侧主窗口，也是系统意义上的主窗口。
- `workspace` 是左侧内容窗口，是 `menu` 的 native 副窗口。
- `menu` 不做整窗透明度衰减，native alpha 保持 `1.0`；窗口背景允许透明像素透出桌面，已绘制的 UI 按自身透明度显示。
- `workspace` 保持整窗 80% 不透明，即 native alpha 为 `0.8`。
- `workspace` 由 `menu` 控制显示、隐藏和导航。
- `workspace` 可以交互，但默认焦点应留在 `menu`，避免切换菜单时被动抢焦点。

两个窗口使用不同的透明实现：`menu` 使用透明窗口背景保留逐像素透明度，`workspace` 使用整窗 native alpha。两种方式都能透过 App 窗口看到桌面下方内容，不使用毛玻璃，也不使用 CSS `opacity` 模拟整窗透明度。

## 窗口结构

窗口定义在 `packages/desktop/src-tauri/tauri.conf.json`：

- App
  - `macOSPrivateApi: true`，在 macOS 下启用 Tauri 的透明窗口背景能力

- `menu`
  - 右侧固定宽度 `208`
  - `decorations: true`
  - `alwaysOnTop: true`
  - `transparent: true`
  - `skipTaskbar: false`
  - 不 responsive，不随窗口变窄收起
- `workspace`
  - 左侧内容窗口，宽度 `992`
  - `decorations: false`
  - `alwaysOnTop: true`
  - `transparent: false`
  - `skipTaskbar: true`
  - `visible: false`
  - `parent: "menu"`

`parent: "menu"` 是关键约束：`workspace` 必须是 `menu` 的 native child/subwindow，不是两个无关系的普通窗口。

## 显示和导航

前端采用两个独立入口：

- `menu.html` → `main-menu.tsx` → `App-menu.tsx`
- `workspace.html` → `main-workspace.tsx` → `App-workspace.tsx`

两个窗口不再共享同一个 `App` 组件，也不通过运行时 window label 判断渲染分支。`App-menu.tsx` 负责菜单状态以及 workspace 的显示、隐藏和导航；`App-workspace.tsx` 负责 workspace 导航监听、路由和页面外壳。

`menu` 窗口渲染 `Sidebar`，维护 `activeMenuPath`：

- 点击未选中 item：显示 `workspace`，并向 `workspace` 发送 `desktop:navigate`。
- 点击已选中 item：清空选中状态并隐藏 `workspace`。
- 点击另一个 item：保持 `workspace` 显示并切换内容。

`workspace` 窗口监听 `desktop:navigate`，收到后修改 hash route。

## 焦点策略

目标行为：

- 用户通过 `menu` 显示或切换 `workspace` 时，默认焦点保持在 `menu`。
- 用户可以连续点击或间隔点击多个 Menu item。
- `workspace` 不应因为 `show()` 或切换内容被动抢焦点。
- 用户主动点击 `workspace` 区域时，`workspace` 仍可获得焦点并正常交互。

实现方式：

1. `menu` 发起显示 `workspace` 前，临时执行 `workspace.setFocusable(false)`。
2. 执行 `workspace.show()` 和 `desktop:navigate`。
3. 短暂等待，避开系统 show 激活窗口的时序。
4. 执行 `workspace.setFocusable(true)`，恢复用户主动交互能力。
5. 执行 `currentWindow.setFocus()`，把默认焦点明确拉回 `menu`。

这不是永久禁用 workspace focus，只是抑制 Menu 驱动 show/switch 时的被动抢焦点。

相关 capability 在 `packages/desktop/src-tauri/capabilities/default.json`：

- `core:window:allow-show`
- `core:window:allow-hide`
- `core:window:allow-set-focus`
- `core:window:allow-set-focusable`

## 透明度实现

### 右侧 Menu：透明窗口背景

`menu` 使用透明窗口背景保留页面逐像素透明度：

1. `packages/desktop/src-tauri/tauri.conf.json` 设置 `app.macOSPrivateApi: true`。
2. 同一文件为 `menu` 设置 `transparent: true`。
3. `packages/desktop/src/index.html` 令 `html`、`body`、`#root` 背景透明。
4. `packages/desktop/src/styles.css` 使用 `html[data-theme]` 覆盖 DaisyUI 写入根节点的主题底色，否则主题白底会挡住透明窗口背景。
5. `packages/desktop/src/layout/Sidebar.tsx` 的 Sidebar 背景保持透明；黄色标题区、文字、按钮等已绘制 UI 继续按自身样式显示。

Rust 侧仍对 `menu` 调用 `setAlphaValue(1.0)`。这里的 `1.0` 表示不对 Menu 最终画面做整窗衰减，不代表把页面中的透明像素强制变成不透明。

这条链路只作用于 `transparent: true` 的 `menu`。`workspace` 保持 `transparent: false`，不使用该逐像素透明方式。

### 左侧 Workspace：整窗 Native Alpha

Rust 入口在 `packages/desktop/src-tauri/src/lib.rs`。

macOS 下使用 `objc2-app-kit` 访问 `NSWindow`，对 `workspace` 调用 `setAlphaValue(0.8)`，使其背景、文字和内容统一参与 80% 整窗合成。

当前两个窗口的整窗 alpha 设置为：

- `menu` 调用 `setAlphaValue(1.0)`
- `workspace` 调用 `setAlphaValue(0.8)`

运行时通过 WindowServer 验证过：

- `desktop menu alpha=1`
- `desktop alpha=0.800000011920929`

## 跟随和生命周期

`workspace` 的位置和尺寸由 Rust 侧同步，不能用 CSS 伪造。

同步规则：

- `workspace` 左侧贴在 `menu` 左侧之外。
- `workspace` 上沿对齐到 Menu 实际菜单内容区域顶部。
- `workspace` 下沿对齐到 `menu` 外框下沿。

当前实现：

- 布局实现集中在 `packages/desktop/src-tauri/src/layout.rs`；`lib.rs` 只负责 Tauri 启动、事件注册和调用布局能力。
- 使用 `menu.outer_position()`、`menu.outer_size()`、`menu.scale_factor()` 计算物理坐标。
- macOS 下读取 `NSWindow.frame()` 和 `NSWindow.contentLayoutRect()` 计算 title bar inset。
- 叠加当前 Sidebar 内部到第一个实际菜单项的布局偏移 `MENU_FIRST_ITEM_TOP_OFFSET`。

生命周期同步：

- `menu` move/resize/scale change 时同步 `workspace`。
- `menu` close requested 时阻止默认 close，隐藏 `workspace`，再隐藏 `menu`。
- `menu` destroyed 时隐藏 `workspace`。
- 后台跟随循环检查 `menu.is_visible()` 和 `menu.is_minimized()`：
  - `menu` 隐藏或最小化时，强制隐藏 `workspace`。
  - `menu` 可见且 `workspace` 已显示时，只同步几何位置，不自动 show。

Tray 行为：

- 显示窗口：只显示并聚焦 `menu`，不自动显示 `workspace`。
- 隐藏窗口：先隐藏 `workspace`，再隐藏 `menu`。
- 置顶切换：`menu` 和 `workspace` 一起切换 always-on-top。

## 已验证行为

在 `pnpm --filter @head/backstage-desktop tauri dev` 下验证过：

- 初始只显示 `desktop menu`，不显示 `desktop` workspace。
- Menu 空白区域能够透出桌面背景，黄色标题区和文字仍按自身样式显示。
- 使用高对比棋盘背景验证过 Menu 的逐像素透明效果。
- 点击 Menu item 后显示 workspace。
- workspace native alpha 为 `0.8`。
- workspace 左侧贴着 Menu 左侧。
- workspace 下沿和 Menu 外框下沿平齐。
- 再次点击同一个 Menu item 后 workspace 隐藏。
- 点击 Menu 关闭按钮后 Menu 和 workspace 都不残留。
- Tray 恢复只显示 Menu，不自动显示 workspace。
- Menu 发起 workspace show/switch 后，默认焦点保持在 Menu。
- 间隔一段时间后继续点击另一个 Menu item 仍能切换。
- `pnpm --filter @head/backstage-desktop tauri build --debug --no-bundle` 构建通过。

## 当前边界

- 双窗口方案天然存在两个可聚焦窗口；当前策略是“Menu 默认持有焦点，workspace 仅在用户主动点击时获得焦点”。
- `workspace` 的顶部对齐依赖当前 Sidebar 内部布局偏移。如果 Sidebar 的 header/nav spacing 改动，需要同步检查 `MENU_FIRST_ITEM_TOP_OFFSET`。
- Menu 透明背景依赖 `macOSPrivateApi`；该能力使用 macOS 私有 API，不适用于需要提交 Mac App Store 的发行方式。
- DaisyUI 会在带 `data-theme` 的根节点写入主题底色，`html[data-theme]` 的透明覆盖规则不能删除，否则 Menu 会重新出现不透明底色。
- 已验证 debug 非打包构建；目前没有 release bundle 验证。

## 开发约束

- 禁止仅为缩短调用处代码、包装一次性日志或隔离实现细节而创建没有稳定语义的函数。抽象必须表达可复用的业务能力、明确的领域概念或独立的错误处理边界。
- Rust 和 `unsafe` 不构成过度抽象的例外。若确实需要用函数收敛 `unsafe` 边界，函数应返回调用方需要的结果或错误，并通过名称表达实际能力；不得创建只执行探测和日志输出的 `log_*` 包装函数。
- 日志应写在信息被使用或错误被处理的位置，使上下文、控制流和日志保持一致。只调用一次且只产生日志的逻辑应保留在调用处，或在没有实际需求时直接删除。
- 对 Tauri 官方支持的能力，应依据官方文档和正确配置使用，不在业务代码中重复验证框架承诺。若官方支持方式存在疑问，应先调研并确认正确接口；若需要验证框架自身行为，应在独立的最小复现项目中隔离验证，不把实验性探测代码带入产品实现。
