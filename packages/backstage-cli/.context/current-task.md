# 当前任务 — CLI 调用方识别能力

> updated_by: Cascade - Claude-Sonnet-4.5
> updated_at: 2026-04-19 17:48:00

## 零、背景

Phase 1 原先的目标是「在 macOS 上证明进程树识别 shell 的思路可行」。TASK-101 完成后首次实测（VS Code + kilo 扩展场景）揭示：

- 进程树里可能**根本没有 shell**（本次：`kilo` → `Code Helper (Plugin)` → `Code`）
- env `SHELL` 表示「用户偏好」，不等于「本次实际调用者」
- 识别对象的语义需要从「shell 枚举」扩展到「任意调用方（shell / IDE / agent / CI / …）」

→ **Phase 1 性质调整为「信号探索」**：先把观测面铺开，在多环境收集样本，再基于数据决定归类模型，**不急着设计 enum 或 type**。

**2026-04-19 第二轮讨论后进一步校准**：目标不是输出"调用方是谁"的单标签，而是输出**一组多维特征**描述"调用方长什么样"。归类留给消费者。Phase 1 的正式产出是一套**特征 Schema v1**，而不是识别代码。

### 前提假设：调用方诚实

本工具**不与调用方对抗**——预设调用方会通过进程关系、环境变量、stdio 形态、专属 env 标记等如实暴露自己的身份。调用方不会主动伪装、也不会刻意隐藏；不同场景下的调用方式客观上不同，"复杂"是**场景差异**，不是对抗。

本工具的定位是**对调用方自述的补充与交叉印证**：相信观测面足够宽时，调用方自然留下的信号已足以拼出它的画像。因此 Schema 追求"宽度优先"——把所有能看到的都记下来，哪怕当前样本用不上。若未来出现对抗场景（例如安全检测），需重新评估本前提。

## 一、决策台账

| # | 决策 | 状态 | 备注 |
|---|------|------|------|
| 1 | 代码位置 | ✅ 定：`internal/agent/service.go` | |
| 2 | 触发入口 | ✅ 定：`backstage agent hasshin` | |
| 3 | 依赖边界 | ✅ 定：纯 stdlib + `ps` | 不允许新装 lib；`lsof` 暂不用（开销大） |
| 4 | 输出契约 | 🔄 过渡：诊断日志 → stderr，识别结果 → stdout | 返回值最终形态待探索收尾 |
| 5 | 落地档位 | ⏸ 绕过 | 探索期不区分 MVP / 完整版 |
| 6 | 归类枚举 | ⏸ 搁置 | 先收样本后设计 |
| 7 | 识别目标语义 | ✅ **特征抽取器，非分类器** | 输出一组多维特征，不是单一标签 |
| 8 | Schema 设计原则 | ✅ **可扩展、不丢信息、只增不减** | schema 服务于观测，不筛选 |
| 9 | Chain 层深限制 | ✅ **无硬数字上限**，用语义停止条件 | 原 `≤ 8` 作废 |
| 10 | 探索成本权衡 | ✅ **探索优先**，允许 fs IO 等重成本 | Phase 1 不做 flag gate |
| 11 | `self` 分组重命名 | ✅ 定：`process` | 重命名发生在 Schema v1 实现前，不触发「只增不减」约束 |
| 12 | 平台策略 | ✅ 单 struct、macOS ∪ Windows 并集 | 缺失字段留空/哨位（数值 -1 / 空串 / nil / false / Time 零）；Linux 不作为一等目标 |
| 13 | 新增 `runtime` 顶层 block | ✅ 放 process 前，描述宿主环境 | 探索阶段采最大信号覆盖；**排除 Go 构建相关字段**（属「探测器自己」不属「被探环境」） |

## 二、探索阶段的设计原则

1. **不是分类器，是特征抽取器** —— 输出一组多维特征，表达"调用方的特征长什么样"，而不是输出一个标签回答"调用方是谁"。归类留给消费者。
2. **不丢信息** —— Schema 是分组和命名，不是筛选。每个观测都必须有落脚点；归不了分组就进 `raw` / `uncategorized`。某一维度"两样本值相同"不代表它没价值——它仍在描述一个语义侧面。
3. **只增不减** —— 分组和字段只允许追加，不允许删除；字段观测不到就留 `null`，绝不省略。
4. **层数用语义停止条件**，不用硬数字 —— `ppid ≤ 1`（到 init / launchd）、环检测、查询失败、安全上限 64（防死循环）。
5. **探索优先，成本不重要** —— fs IO / 多次 ps 调用这类成本在 Phase 1 不做 gate，先采尽。

## 三、双样本对比与核心发现（Phase 1 截至目前）

### 样本清单

- **样本 1** Windsurf + Cascade (Claude Opus 4.7 Max)：`.context/observations/windsurf-claude-opus.log`
- **样本 2** VS Code + kilocode 扩展：**待保存到** `.context/observations/vscode-kilo.log`（见 TASK-104.1）

### 双样本对比

| 维度 | Windsurf + Cascade | VS Code + kilo |
|------|--------------------|-----------------|
| 直接父 comm | `zsh -il` | `kilo serve --port 0` |
| 祖先 L1 | `Windsurf Helper utility`（node service） | `Code Helper (Plugin)`（node service） |
| 祖先 L2 | `Windsurf.app` | `Visual Studio Code.app` |
| stdin | tty | tty |
| stdout | tty | **socket** |
| stderr | tty | **socket** |
| `SHELL` | `/bin/zsh` | `/bin/zsh` |
| `TERM` / `TERM_PROGRAM` | 有 / `vscode`，版本字串含 `windsurf` | 都空 |
| 专属 env | `WINDSURF_CASCADE_*`、`VSCODE_*`（fork 继承） | `VSCODE_CRASH_REPORTER_PROCESS_TYPE=extensionHost` 等 |
| 老 hasshin stdout 返回 | `zsh`（语义颠倒的"命中"） | `unknown`（语义正确） |

### 核心发现

1. **进程树里"没 shell"是常态，不是例外。** kilo 这条链 `kilo` → `Code Helper (Plugin)` → `Code` 全程不经 shell；Phase 1 最初"爬到 shell 就结束"的假设在 extension-host 模式下直接失效。
2. **stdio 的 `(tty, socket, socket)` 是 RPC agent 的指纹**，且完全独立于进程树和环境变量。Cascade 把命令塞进用户交互 tty，三路都是 tty；kilo 通过 socket 接管 stdout/stderr 回传给扩展。
3. **环境变量比进程树更可靠**。`WINDSURF_CASCADE_TERMINAL_KIND`、`VSCODE_CRASH_REPORTER_PROCESS_TYPE=extensionHost` 是二值判定、零误判；进程链依赖字符串模糊匹配，易被路径/版本差异打破。
4. **老 hasshin 的返回值语义颠倒**。kilo 下返回 `unknown` 是"对的"（链里真没 shell），Windsurf 下返回 `zsh` 是"错的"（Cascade 套的外壳 zsh 不代表用户在 zsh 里敲命令）。→ 输出不能合并成单一"shell 字符串"，必须分开描述调用方多个侧面。
5. **祖先链的"层数"不是语义单位，"角色"才是。** 两样本的 L1 都是 Electron helper utility、L2 都是 GUI app——这不是冗余，是"相同角色 + 不同身份"。但继续往上（L3 = launchd）在两样本里确实是无新信息的冗余层。所以 schema 要按**角色**组织层（`direct_parent` / `hosting_runtime` / `ui_root` / `intermediaries`），但**原始层数据全量保留**，角色是叠加在层上的 tag。

## 四、特征 Schema v1（Phase 1 正式产出）

### 骨架

```
{
  runtime: { os, arch, hostname, cpu_count, cpu_cores_physical, cpu_model,
             memory_bytes, os_version, os_build, kernel_version,
             distro, session_type, virtualization },
  process: { ids, identity, windows_identity, argv, exe, cwd, start_time, raw },
  tty:     { controlling, foreground_pgid, size },
  chain:   { depth, termination, layers: [ { pid, ppid, pgid, sid,
                                              uid, username, comm, command, tty,
                                              lstart, raw_ps_line } × N ] },
  env:     { shell_pref, terminal, ide, extension_host, session, ci,
             locale, editor, path, user, shell_version, remote,
             uncategorized, raw },
  time:    { now, self_start, parent_start, uptime },
  fs:      { cwd, home, tmpdir, repo_root?, ide_workspace_markers?, path_dirs },
  exec:    { argv0, resolved_exe, relation },
  raw:     { ... 任何还没归类的原始观测 ... }
}
```

### 各组字段设计

#### `runtime` — 宿主环境

采用 **macOS ∪ Windows 并集 schema**（同 `process`）；描述本 CLI 被执行时的宿主机器 / OS 上下文快照。Linux 不作为一等目标，合辙的顺带可用，不合辙的按字段自行留空。

**Scope**：仅记录「被探测环境」；不包含本 CLI binary 自身的构建身份（Go 版本、build target、cliVersion 等）——那些属于「探测器自己」，若需会走独立 block 承载。

字段默认跨平台；平台限定的字段在名字后标注 `(mac)` / `(win)` / `(linux)` / `(all)`。

- `os`: `runtime.GOOS`；如 `"darwin"` / `"windows"` / `"linux"`
- `arch`: `runtime.GOARCH`；如 `"arm64"` / `"amd64"`
- `hostname`: `os.Hostname()`
- `cpu_count`: `runtime.NumCPU()`；进程可用的 logical CPU 数
- `cpu_cores_physical`: 物理核心数；mac: `sysctl hw.physicalcpu`
- `cpu_model`: CPU 型号字串；mac: `sysctl machdep.cpu.brand_string`
- `memory_bytes`: 物理总内存字节数；mac: `sysctl hw.memsize`
- `os_version`: 用户可读版本号；mac: `"14.3"`；win: `"11 23H2"`
- `os_build`: 细粒度 build 号；mac: `"23D60"`；win: `"22621.3007"`
- `kernel_version`: 内核版本；mac: `uname -r`；win: NT 版本号
- `distro` (linux): Linux 发行版，如 `"Ubuntu 22.04"`；mac/win 留空
- `session_type` (linux): `XDG_SESSION_TYPE` —— `wayland` / `x11` / `tty`；mac/win 留空
- `virtualization` (all): 虚拟化 / 仿真标记；可能值 `"rosetta"` / `"wsl"` / `"docker"` / `"vm"` / `""`；TASK-105 / Phase 2 实现

> 哨位约定：整数字段采集失败或不适用留 **0**（本 block 的整数字段 0 天然非合法值，等价哨位）；字符串字段采集失败写入 `"<err: ...>"`。
>
> 采集策略：跨平台字段（`os` / `arch` / `hostname` / `cpu_count`）全部走 Go stdlib；平台专属字段（`os_version` 及 `sysctl` 采的 cpu/mem/...）按 OS 分支 exec 系统命令。

#### `process` — 运行身份

采用 **macOS ∪ Windows 并集 schema**；描述本进程的运行身份快照，字段均可从 kernel 直接读取。Linux 不作为一等目标，合辙的顺带可用，不合辙的按字段自行留空。

字段默认跨平台；平台限定的字段在名字后标注 `(mac)` / `(win)` / `(linux)`。

- `ids`: `pid` / `ppid` / `pgid` (mac) / `sid` (mac；Unix 会话号)
- `identity`: `uid` (mac) / `euid` (mac) / `gid` (mac) / `egid` (mac) / `username` / `groupname`
- `windows_identity` (win): `user_sid` / `group_sid` / `session_id`（Windows 登录会话，与 `sid` 无关）/ `domain` / `integrity_level` / `elevated`
- `argv`: 完整 argv 数组（不只 argv[0]）
- `exe`: 详见 `exec` 组
- `cwd`: 当前工作目录
- `start_time`: 本进程启动时刻（两边都需平台代码，stdlib 不直接给）

> 哨位约定：整数字段采集失败或平台不适用留 **-1**（0 在 uid/gid/pid 等字段上是合法值，不能当哨位）；字符串留空串，bool 留 `false`，slice 留 `nil`，`time.Time` 留零值；字符串字段采集失败写入 `"<err: ...>"`。

> stdio 形态不再作为独立 Schema 块；改由**特征指纹**承载，详见第八节 FP-01。

#### `tty` — 控制终端归属

- `controlling`: 本进程会话是否拥有控制 tty、设备名
- `foreground_pgid`: 本进程是否在 tty 的前台进程组
- `size`: 终端尺寸 `(rows, cols)` via `TIOCGWINSZ`

#### `chain` — 祖先链（仅 ps 直接返回的原始字段）

每层保留全部字段：

- `pid` / `ppid` / `pgid` / `sid`
- `uid` / `username`
- `comm` / `command`
- `tty`
- `lstart`
- `raw_ps_line`（原始 ps 输出，防止解析出错时丢信息）

> 角色分类 / 规则匹配 等二次判定不在 schema 里做——留给第八节特征指纹处理（候选素材见 FP-02 候选）。

**遍历策略**：不设硬层数上限，按以下语义停止条件结束：

- `ppid ≤ 1` → 到 init / launchd，自然终止
- 已访问 `pid` 出现第二次 → 环（数据异常），立即停
- `ps` 查询失败 → 记录原因、停
- 访问层数达到 64 → 安全上限保护（防死循环，**不是业务约束**）

`chain.depth` 本身是特征。`depth=3` 和 `depth=13` 描述完全不同的调用形态（直调 vs tmux 嵌套 vs ssh+tmux+sudo）。

`chain.termination` 记录终止方式，也是特征：

- `init`: 自然走到 launchd / init
- `depth_cap`: 触到 64 层保护
- `cycle`: 环检测命中
- `query_failed`: `ps` 失败
- `orphan`: `ppid=0`（父死了但未挂到 init）

#### `env` — 按语义家族分组 + 原始全量

**实现形态**：每家族是 `map[string]string`（家族内命中的 key/value 全量落入）；
未命中任何家族的 env 落 `uncategorized`；**全量 `os.Environ()` 同时保留在
`env.raw`**（`[]string`，原始顺序和可能的重复均保留不丢）。分组只是语义包装。

**归属规则**：每个 env 最多归一个家族（**单归属**），按特异性优先级匹配
（更窄的家族在前）；exact match 优先 prefix match。

**Windows 适配**：Windows 的 env key 大小写不敏感（`Path` / `PATH` / `path` 同一键），
已在 `classifyEnvKey` 加入 **UPPER(key) fallback**，mac 不受影响；下表的 Windows
列已填入 exact map。

**Windows User/System Scope 限制**：Windows 环境变量分 User scope 和 System scope，
加载顺序为 System → User，**同名变量 user 覆盖 system**。`os.Environ()` 只返回最终合并视图，
无法区分变量来源 scope。若需区分，需直接读取注册表（`HKCU\Environment` vs
`HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`）。当前实现
不区分 scope，`platform` 组的 Windows 键（`USERPROFILE` / `APPDATA` 等）来源可能是
user scope 或 system scope 继承。

| 家族 | Unix 覆盖的 env | Windows 覆盖的 env |
|------|-----------------|--------------------|
| `shell_pref` | `SHELL` | |
| `terminal` | `TERM` / `TERM_PROGRAM` / `TERM_PROGRAM_VERSION` / `COLORTERM` / `LC_TERMINAL` / `KITTY_WINDOW_ID` / `WEZTERM_EXECUTABLE` / `ALACRITTY_SOCKET` / `TERMINAL_EMULATOR` | |
| `ide` | 前缀：`VSCODE_*` / `WINDSURF_*` / `CURSOR_*` / `CASCADE_*` / `JETBRAINS_*` / `IDEA_*` / `PYCHARM_*` | 同 Unix（同前缀表） |
| `extension_host` | `VSCODE_CRASH_REPORTER_PROCESS_TYPE` / `VSCODE_PID` / `VSCODE_IPC_HOOK` / `VSCODE_ESM_ENTRYPOINT` / `VSCODE_NLS_CONFIG` | 同 Unix |
| `session` | `SSH_CONNECTION` / `SSH_CLIENT` / `SSH_TTY` / `TMUX` / `STY` | 同 Unix 子集 |
| `ci` | `CI` / `GITHUB_ACTIONS` / `GITLAB_CI` / `CIRCLECI` / `BUILDKITE` / `JENKINS_URL` / `TRAVIS` | 同 Unix |
| `locale` | `LANG` / `LC_*` | |
| `editor` | `EDITOR` / `VISUAL` / `PAGER` | |
| `platform` | `PATH` / `HOME` / `TMPDIR` | `PATH` / `USERPROFILE` / `TEMP` / `TMP` / `APPDATA` / `LOCALAPPDATA`（case-insensitive fallback） |
| `shell` | | `PSModulePath` |
| `user` | `USER` / `LOGNAME` | `USERDOMAIN` |
| `shell_version` | `BASH_VERSION` / `ZSH_VERSION` / `FISH_VERSION` | |
| `remote` | `KUBERNETES_SERVICE_HOST` / `IN_NIX_SHELL` | |
| `uncategorized` | 其他未归类（包括机器级 / 系统级键如 `COMPUTERNAME` / `SYSTEMROOT` / `WINDIR` 等） | |
| `raw` | 全量 `os.Environ()`（`[]string`，顺序 + 重复保留） | |

**Platform 组专用输出结构**：`EnvPlatformBlock` 是 `platform` 组的专用输出 struct，
通过 `printBlockByReflect` 打印（而非 `printEnvFamily`）。包含以下字段（均为 loglevel=debug）：
- `Path`：PATH 环境变量
- `HOME` / `USERPROFILE`：用户主目录（Mac 用 HOME，Windows 用 USERPROFILE）
- `TMPDIR` / `TEMP` / `TMP`：临时目录（Mac 用 TMPDIR，Windows 用 TEMP/TMP）
- `APPDATA` / `LOCALAPPDATA`：Windows 应用数据目录

字段命名使用原名（不跨平台统一），Mac 的 key 用原名，Windows 的 key 用原名。
`EnvBlockLayers.Platform` 是 `EnvPlatformBlock` struct 类型，直接填充并打印。

**Shell 组专用输出结构**：`EnvShellBlock` 是 `shell` 组的专用输出 struct，
通过 `printBlockByReflect` 打印。包含以下字段（均为 loglevel=debug）：
- `PSModulePath`：Windows PowerShell 模块路径

`EnvBlockLayers.Shell` 是 `EnvShellBlock` struct 类型，直接填充并打印。

#### `time` — 时间信号

- `now`: 观测时刻
- `self_start`: 本进程启动时刻
- `parent_start`: 直接父启动时刻
- `uptime`: 离系统启动多久
- （祖先的 `lstart` 落在 `chain.layers[].lstart`）

**语义用途示例**：`parent_start` 与 `self_start` 差距小 → 父刚 fork 我；差距大 → 我被长期运行的 daemon 收养。

#### `fs` — 文件系统上下文

- `cwd` / `home` / `tmpdir`
- `repo_root`: 沿 cwd 向上爬找 `.git`，找到即记（最近的 git 仓库根）
- `ide_workspace_markers`: 沿 cwd 向上爬，检测 `.vscode/` / `.idea/` / `*.code-workspace`
- `path_dirs`: `$PATH` 分割后的目录列表
- 注：fs IO 有成本，Phase 1 **不 gate**（探索优先）

#### `exec` — 被如何 exec 的

- `argv0`: 字面 `os.Args[0]`
- `resolved_exe`: 通过 macOS 的 `_NSGetExecutablePath` / Linux 的 `/proc/self/exe` 解析出的真实路径
- `relation` ∈ {`identical`, `via-symlink`, `via-path-lookup`, `relative-path`, `unknown`}

### 强制规则

1. **观测到就记录，观测不到留 `null`，绝不省略**
2. **新信号没归类 → 进 `raw` 或对应分组的 `.raw`**
3. **Schema 只增不减** —— 新增分组 / 字段需要在本文档同步追加，历史分组和字段保留
4. **诊断日志输出按 Schema 结构分块**，便于后续程序化解析

## 五、Phase 1 任务

> 目标：在多种调用环境下运行 `backstage agent hasshin`，收集每次 stderr 诊断输出作为样本，为 Phase 2 识别模型提供数据基础。

### TASK-101 (DONE)：macOS 原型 + 接入 hasshin

**实现侧**

- `internal/agent/service.go` 新增 `Hasshin() (string, error)`
- 进程树一级父进程查询 + 硬编码 shell 白名单匹配

**接入侧**

- `cmd/agent.go` 的 `hasshinCmd.RunE` 已挂接

### TASK-102 (DONE)：首版诊断日志

- 给 `Hasshin()` 加 stderr 诊断输出
- 字段：pid / ppid / parent comm / parent command / 祖先链 comm（≤ 8）/ env SHELL / TERM / TERM_PROGRAM
- 首次实测（VS Code + kilo）返回 `unknown`，暴露了「零、背景」里的核心发现

### TASK-103 (DONE)：扩充观测面到 5 块信号

**A. 自身**

- pid / ppid / uid / euid / username / cwd / argv[0]

**B. stdio fd 类型**

- stdin / stdout / stderr 各自：tty / pipe / file / socket / other
- 通过 `os.Stdin.Stat().Mode()` 的 `ModeCharDevice` / `ModeNamedPipe` / `ModeSocket` 判断

**C. 进程链（≤ 8 层）**

- 每级：pid / ppid / tty / user / lstart / comm / full command
- ps 字段：`-o pid=,ppid=,tty=,user=,lstart=,comm=,command=`

**D. env: shell & terminal**

- `SHELL` / `BASH_VERSION` / `ZSH_VERSION` / `FISH_VERSION` / `PSModulePath`
- `TERM` / `TERM_PROGRAM` / `TERM_PROGRAM_VERSION` / `COLORTERM` / `LC_TERMINAL`
- `TMUX` / `KITTY_WINDOW_ID` / `WEZTERM_EXECUTABLE` / `ALACRITTY_SOCKET` / `TERMINAL_EMULATOR`

**E. env: ide / ci / remote**

- 所有前缀为 `VSCODE_` / `CURSOR_` / `JETBRAINS_` / `IDEA_` / `PYCHARM_` 的变量（遍历 `os.Environ()` 匹配前缀）
- `CI` / `GITHUB_ACTIONS` / `GITLAB_CI` / `TRAVIS` / `CIRCLECI` / `BUILDKITE` / `JENKINS_URL`
- `SSH_CLIENT` / `SSH_TTY` / `SSH_CONNECTION`
- `KUBERNETES_SERVICE_HOST` / `IN_NIX_SHELL`

**产物**：`Hasshin()` 内部 stderr 诊断扩充为 5 个标题块。stdout 暂保留老逻辑（命中白名单返回名字，否则 `unknown`），便于现有验收命令继续可用。

### TASK-104 (DONE)：采样 Windsurf + Cascade 环境

- 样本落地：`.context/observations/windsurf-claude-opus.log`
- 关键观察与发现已汇总到第三节「双样本对比与核心发现」

### TASK-104.1 (TODO)：采样 VS Code + kilocode 环境

- 第二轮讨论期间已在终端捕获输出（父链为 `kilo serve` → `Code Helper (Plugin)` → `Code`），**尚未落盘**
- 目标文件：`.context/observations/vscode-kilo.log`
- 需要在 TASK-105 重构后重新采一次，以 Schema v1 格式落地（旧格式样本可先保留一份作为历史对照）

### TASK-105 (IN PROGRESS)：按 Schema v1 重构 `Hasshin()` 诊断输出

**约束**

- 诊断日志按第四节 Schema 骨架分块输出：`runtime` / `process` / `tty` / `chain` / `env` / `time` / `fs` / `exec` 各一块（stdio 原先的 `io` 块已废除，改由 FP-01 承载）
- 每块内部字段严格对齐 Schema；观测不到留 `null`，绝不省略
- `chain` 遍历改为语义停止条件（见第四节），废除硬 `≤ 8` 限制
- 新增 `env.raw` 落全量 `os.Environ()`；未归类 env 进 `env.uncategorized`
- stdout 老行为保留（白名单命中返回 shell 名，否则 `unknown`），便于现有回归命令继续可用
- 诊断块之间用固定分隔符（现已切换为 `[HASSHIN] ---- Xxx ----`），便于后续程序化解析

**已完成**

- ✅ **第一轮（2026-04-19 上午）**
  - `runtime` 块：`RuntimeBlock` + `dumpRuntimeBlock` + mac 采集器（sw_vers / uname / sysctl）
  - `process` 块：`dumpSelfBlock` → `dumpProcessBlock` + `ProcessBlock` 升级为 union schema（含 Windows 专属字段占位）
  - Schema v1 与 Go struct 双向对齐：field doc / 哨位约定 / union 策略
- ✅ **第二轮（2026-04-19 下午）**
  - stdio / B 块移除：改由 FP-01 特征指纹承载（见第八节），**不再新增 `IOBlock`**
  - `chain` 块：`ChainBlock` + `dumpChainBlock` + mac `queryChainLayerDarwin`（ps）；语义停止条件替换硬 8 层上限；角色分类下沉到 FP-02 候选
  - `env` 块：`EnvBlock` + `dumpEnvBlock`；D/E 二分合并为单块，按 family 重组，含 `uncategorized` 兼底 + `raw []string` 全量；stderr 末尾打全量 `all_keys (sorted)` 清单
  - Windows 代码预备：`queryChainLayer` GOOS 分派将 mac 实现圈定为 `queryChainLayerDarwin`，其他平台返回明确 stub err；`classifyEnvKey` 加 Windows `UPPER(key)` fallback；`envExactMap` 补 6 个 Windows 用户 env 键。mac 完全向后兼容。
- ✅ **第三轮（2026-04-20）**
  - Windows chain 完整实现：`chain_windows.go` 使用 `github.com/shirou/gopsutil/v3/process` 获取进程信息（PID/PPID/Name/Cmdline/Username/CreateTime）
  - 代码架构优化：`queryChainLayer` 中间函数移除，直接调用 `queryChainLayerOS`；`rawAfterNTokens` 移至 `chain_darwin.go`；`chain_other.go` 直接返回 `framework.NotImplementedException`
  - stderr 格式调整：`[HASSHIN]` 单独一行在开头，各 block 用 `---- Xxx ----`；env 输出只打印 key 不打印 value

**待办（后续轮次）**

- `tty` / `time` / `fs` / `exec` 四个新 block
- Windows RuntimeBlock 和 ProcessBlock 的 union 右半部分完整实现（当前 mac 专属字段留零值）
- Linux 平台采集代码的**完整实现**（当前 chain 返回 NotImplementedException）

### TASK-106 (TODO)：按 Schema v1 重新采样已有两个环境

- Windsurf + Cascade
- VS Code + kilocode
- 输出覆盖到 `.context/observations/` 下同名文件，旧格式 log 保留为 `*.legacy.log` 作为历史对照

### TASK-107 (HH)：Phase 1 复盘评审

- 输入：TASK-106 的两份 Schema v1 样本
- 评审内容：
  - 哪些维度在两样本里有差异 / 哪些同质 / 哪些只单侧出现
  - Schema 是否需要新增 / 调整分组或字段（遵循「只增不减」原则）
  - Phase 2 识别模型候选方案（可能是「维度到 caller 类型」的规则表）
  - `Hasshin()` 的 API 契约草案（是否需要改 `(string, error)` 签名）
- 产出：Phase 2 启动所需的设计输入

## 六、Phase 2 — 识别模型与实现（样本驱动，待定）

> 具体 TASK 待 Phase 1 TASK-107 出结论后补齐。预期包含：

- 基于样本定义识别模型（多维标签组合 → caller 类型）
- API 设计（结构体 vs 字符串）
- 实现与 Phase 1 原型代码重构
- Windows 平台支持
- 测试策略
- 双平台人工验收

## 七、约定提醒

- 整体编译 / 运行验证仅在 Phase 收尾进行，由 Human 执行：`go mod tidy` + `./build-mac.sh` / `./build-windows.bat`
- 每 Phase 最后一个 TASK 必须是 **人工验收（HH）**
- 终端输出 / 交互类变更无法自动化，一律走 HH
- 本文档随任务推进持续更新；完成的 TASK 改为 `(DONE)`；新决策追加到台账
- Schema 变动必须同步更新第四节

## 八、特征指纹

### 定位

**特征指纹 ≠ Schema v1 观测字段**。两者职责互补：

- **Schema v1**：宽度优先的事实容器，记录"看到了什么"，不做判定，**只增不减**
- **特征指纹**：在观测之上提取的**高信号组合**，是"一眼就指向某类调用方"的 hint 集合

指纹**不是判定**，只是 Phase 2 识别模型优先喂入的候选信号。随样本积累可调整，**不受 Schema 的"只增不减"约束**。

### 前置假设

同第零节的"调用方诚实"前提——指纹的信号强度建立在"调用方不主动伪装"之上。

### FP-01：stdio 三路形态

**信号源**：`stdin` / `stdout` / `stderr` 各自的 fd 形态，通过 `os.Stdin.Stat().Mode()` 的 `ModeCharDevice` / `ModeNamedPipe` / `ModeSocket` 判断。取值集合：`tty` / `pipe` / `socket` / `file` / `null` / `other`。

**核心性质**

- **独立于进程树和 env** —— kernel 真实状态，不受 fork 继承污染、不受字符串模糊匹配打破
- **低成本** —— 一次 `Fstat` 出结果，不爬进程链
- **三路独立、不折叠成元组** —— **不对称本身是指纹**

**形态映射（持续扩充）**

| `(stdin, stdout, stderr)` | 调用模式 |
|--------------------------|---------|
| `(tty, tty, tty)` | 用户在交互终端敲命令；**或** agent 把命令塞进用户的 tty（Cascade 即此种） |
| `(tty, socket, socket)` | RPC agent — 用户 tty 进，结果走 socket 回扩展宿主（kilo 即此种） |
| `(pipe, pipe, pipe)` | 被另一进程 `exec.Command` 起来 / shell 管道 |
| `(file, file, file)` 或混合 | 后台 / `nohup` / CI redirect |
| `(null, null, null)` | daemon / launchd 托管 |

**双样本印证**

- Windsurf + Cascade：`(tty, tty, tty)` — 套壳 tty
- VS Code + kilo：`(tty, socket, socket)` — 经典 RPC agent 指纹

**与 Schema 其他块的互证**

- 与 `tty` 块：fd 是 `tty` 只说"这个 fd 连到字符设备"；`tty.controlling` 才说"会话有没有控制终端"。背离情况：daemonize 后失去 controlling，但 fd 仍可能连某个 pts。
- 与 `chain` 块：`(tty, socket, socket)` + 链里全无 shell → 高置信度 RPC agent；`(tty, tty, tty)` + 链里没 shell → Cascade 这种套壳模式。

**附加字段（采集时按需补）**

- `tty` 时：`tty_name` / `isatty` / termios
- `socket` 时：`family`（`unix` / `inet` / `inet6`） / unix domain peer path（如可取）
- `pipe` 时：对端若可识别就记

### FP-02（候选）：chain 层角色分类

**状态**：候选池，未正式立项；样本积累后可能升格为 FP-02。

**定位**：在 `chain.layers[*]` 的原始字段（`comm` / `command` / 层位置等）之上，按规则打上角色标签。**这是消费者行为，不是采集时写入 schema**。

**候选规则分层**

1. **层位置 role**（零歧义，从 slice 索引直出）
   - `direct_parent`：`layers[0]`
   - `intermediaries`：非首非末

2. **comm 精确匹配 role**（零歧义；comm normalize 后查表）

   normalize 步骤：`TrimPrefix "-"` + `ToLower` + `TrimSuffix ".exe"`

   | comm（normalized） | role |
   |---------------------|------|
   | `bash` / `zsh` / `sh` / `fish` / `dash` / `ksh` / `tcsh` / `csh` / `pwsh` / `powershell` / `cmd` | `shell` |
   | `tmux` / `screen` | `session-manager` |
   | `launchd` | `launchd` |
   | `init` | `init` |
   | `sshd` / `sudo` / `cron` | 同名 |

3. **启发式 role**（需更多样本定规则）：`ui-app` / `extension-host` / `terminal-emulator` / `hosting_runtime` / `ui_root` / `wrapper-script` / `daemon`。可能需要结合 `command` 参数（如 `--type=extensionHost`）、路径特征（如 `*.app/Contents/`）、Electron helper 命名模式等。

**为什么拆出来单列**

- Schema 只记原始事实，保持「观测 vs 判定」边界清晰
- 规则 / 映射表未来会演化，和 schema 的「只增不减」节奏不同
- Phase 2 消费者可能有不同的分类需求，schema 不应预先决定分类口径