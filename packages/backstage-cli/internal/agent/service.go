package agent

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

// Hasshin 采集调用方的运行环境特征（Phase 1 信号探索版）。
//
// **平台范围**：Phase 1 聚焦 macOS；Windows 代码路径已预留（保证编译
// 通过）但运行时 chain 块退化为 stub（ps 没有 Windows 实现）、
// runtime 块的 macOS-专属字段留零值、process 块的 uid/euid 由 Go
// stdlib 在 Windows 属性上自然返回 -1。Linux 不作为一等目标。
//
// 探索期策略：stderr 按块输出多维观测信号，用于在多种调用环境下采样分析；
// 不做分类，只做特征抽取（见 .context/current-task.md 决策 #7）。stdout
// 暂沿用老 shell 白名单识别（命中 bash/zsh/sh/fish 返回名字，否则 "unknown"），
// 便于现有回归命令可用；归类模型待样本充足后再设计。
//
// 信号分块（各块 schema 定义详见 internal/agent/types.go；当前
// RuntimeBlock / ProcessBlock / ChainBlock / EnvBlock 已落地）：
//   runtime  (宿主环境)
//   process  (本进程)
//   chain    (祖先链快照；语义停止条件终止，无硬层数上限)
//   env      (按 family 分组 + uncategorized + raw 全量；TASK-103 旧 D/E 二分已合并)
//
// stdio 信号不再作为独立块，改由特征指纹 FP-01 承载（见 .context/current-task.md 第八节）。
func Hasshin() (string, error) {
	dumpRuntimeBlock()
	dumpProcessBlock()
	parentComm := dumpChainBlock()
	dumpEnvBlock()

	// ========== Final: legacy whitelist match ==========
	fmt.Fprintln(os.Stderr, "[hasshin] --- final ---")
	name := filepath.Base(parentComm)
	name = strings.TrimPrefix(name, "-")
	fmt.Fprintf(os.Stderr, "  normalized parent comm: %q\n", name)

	switch name {
	case "bash", "zsh", "sh", "fish":
		return name, nil
	default:
		return "unknown", nil
	}
}

// dumpRuntimeBlock 采集并输出 runtime 块（宿主环境快照）：os / arch /
// hostname / cpu / memory / OS 版本信息 / linux 提示 / virtualization。
// 结果写入 stderr，供信号探索采样分析使用。
//
// 平台分派：mac 侧通过 sw_vers / uname / sysctl 采集
// （见 collectDarwinRuntime）；linux / windows 的平台代码留待
// TASK-105 或 Phase 2。在非 mac 平台上运行时，mac 侧专属字段保持
// 零值 / 空串。
func dumpRuntimeBlock() {
	b := RuntimeBlock{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		CPUCount: runtime.NumCPU(),
	}
	if h, err := os.Hostname(); err == nil {
		b.Hostname = h
	} else {
		b.Hostname = fmt.Sprintf("<err: %v>", err)
	}

	switch runtime.GOOS {
	case "darwin":
		collectDarwinRuntime(&b)
		// case "linux":   // TASK-105 / Phase 2
		// case "windows": // TASK-105 / Phase 2
	}

	fmt.Fprintf(os.Stderr, "[HASSHIN] ---- Runtime ----\n")
	fmt.Fprintf(os.Stderr, "  os=%q arch=%q hostname=%q\n", b.OS, b.Arch, b.Hostname)
	fmt.Fprintf(os.Stderr, "  cpu_count=%d cpu_cores_physical=%d cpu_model=%q\n", b.CPUCount, b.CPUCoresPhysical, b.CPUModel)
	fmt.Fprintf(os.Stderr, "  memory_bytes=%d\n", b.MemoryBytes)
	fmt.Fprintf(os.Stderr, "  os_version=%q os_build=%q kernel_version=%q\n", b.OSVersion, b.OSBuild, b.KernelVersion)
	fmt.Fprintf(os.Stderr, "  distro=%q session_type=%q virtualization=%q\n", b.Distro, b.SessionType, b.Virtualization)
}

// collectDarwinRuntime 填充 RuntimeBlock 的 mac 侧字段，通过 sw_vers / uname /
// sysctl 采集。字符串字段采集失败时写入 "<err: ...>"；整数字段采集失败留零。
func collectDarwinRuntime(b *RuntimeBlock) {
	setExecStr(&b.OSVersion, "sw_vers", "-productVersion")
	setExecStr(&b.OSBuild, "sw_vers", "-buildVersion")
	setExecStr(&b.KernelVersion, "uname", "-r")
	setExecStr(&b.CPUModel, "sysctl", "-n", "machdep.cpu.brand_string")
	setExecInt(&b.CPUCoresPhysical, "sysctl", "-n", "hw.physicalcpu")
	setExecInt64(&b.MemoryBytes, "sysctl", "-n", "hw.memsize")
}

// setExecStr 运行命令并把 trim 后的 stdout 写入 dst；失败时写入 "<err: %v>"。
func setExecStr(dst *string, name string, args ...string) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		*dst = fmt.Sprintf("<err: %v>", err)
		return
	}
	*dst = strings.TrimSpace(string(out))
}

// setExecInt 运行命令并把 trim 后的 stdout 解成十进制 int 写入 dst；失败保持零值。
func setExecInt(dst *int, name string, args ...string) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return
	}
	if n, perr := strconv.Atoi(strings.TrimSpace(string(out))); perr == nil {
		*dst = n
	}
}

// setExecInt64 同 setExecInt，dst 为 *int64。
func setExecInt64(dst *int64, name string, args ...string) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return
	}
	if n, perr := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64); perr == nil {
		*dst = n
	}
}

// dumpProcessBlock 采集并输出 A 块（process 分组，本进程运行身份）：
// pid / ppid / uid / euid / username / cwd / argv。结果写入 stderr，
// 供信号探索采样分析使用。
//
// 平台分派：当前仅填充跨平台 stdlib 能给出的字段；Unix-only 数值 ID
// （pgid / sid / gid / egid）及 Windows 专属字段（SID / SessionID /
// domain / integrity / elevated）留待 TASK-105 / Phase 2 填充。
func dumpProcessBlock() {
	b := ProcessBlock{
		PID:  os.Getpid(),
		PPID: os.Getppid(),
		UID:  os.Getuid(),
		EUID: os.Geteuid(),
		Argv: os.Args,
	}
	if u, err := user.Current(); err == nil {
		b.Username = u.Username
	} else {
		b.Username = fmt.Sprintf("<err: %v>", err)
	}
	if cwd, err := os.Getwd(); err == nil {
		b.CWD = cwd
	} else {
		b.CWD = fmt.Sprintf("<err: %v>", err)
	}
	fmt.Fprintf(os.Stderr, "[HASSHIN] ---- Process ----\n")
	fmt.Fprintf(os.Stderr, "  pid=%d ppid=%d uid=%d euid=%d username=%q\n", b.PID, b.PPID, b.UID, b.EUID, b.Username)
	fmt.Fprintf(os.Stderr, "  cwd=%q\n", b.CWD)
	fmt.Fprintf(os.Stderr, "  argv=%q\n", b.Argv)
}

// maxChainDepth 是 chain 遍历的死循环保护上限，**不是业务约束**。
// 正常链通常 3~10 层；触达上限以 Termination="depth_cap" 停止。
const maxChainDepth = 64

// dumpChainBlock 采集 C 块（chain，祖先链快照）并写入 stderr。
//
// 策略：从 os.Getppid() 起沿 ppid 上爬；以语义停止条件终止
// （见 ChainBlock doc）。**本函数只采集 ps 直接返回的原始信息**
// （原始字段 + RawPsLine 全量保留），不做角色分类 / 指纹匹配等
// 二次判定 —— 那些逻辑归特征指纹章节处理（见 .context/current-task.md
// 第八节）。
//
// 返回值 parentComm 是 layers[0].Comm（即 os.Getppid() 的 ps -o comm=
// 原值，可能带 "-" 前缀），供 Hasshin final whitelist 段沿用老逻辑。
// 链为空时返回 ""。
func dumpChainBlock() (parentComm string) {
	b := ChainBlock{}
	visited := make(map[int]bool)
	cur := os.Getppid()

	for len(b.Layers) < maxChainDepth {
		if cur <= 1 {
			if cur == 1 {
				// 把 init/launchd 这一终止层也记进来（成功与否都记，不丢信息）
				layer, _ := queryChainLayer(cur)
				b.Layers = append(b.Layers, layer)
				b.Termination = "init"
			} else {
				// cur == 0：极罕见的进程失孤
				b.Termination = "orphan"
			}
			break
		}
		if visited[cur] {
			b.Termination = "cycle"
			break
		}
		visited[cur] = true

		layer, err := queryChainLayer(cur)
		if err != nil {
			// 失败层也记录进来：保留 pid + 错误详情，不丢信息
			b.Layers = append(b.Layers, layer)
			b.Termination = "query_failed"
			break
		}
		b.Layers = append(b.Layers, layer)

		cur = layer.PPID
		if cur == 0 {
			b.Termination = "orphan"
			break
		}
	}
	if b.Termination == "" {
		b.Termination = "depth_cap"
	}
	b.Depth = len(b.Layers)

	fmt.Fprintf(os.Stderr, "[HASSHIN] ---- Chain ----\n")
	fmt.Fprintf(os.Stderr, "  depth=%d termination=%q\n", b.Depth, b.Termination)
	for i, l := range b.Layers {
		fmt.Fprintf(os.Stderr, "  layer[%d] pid=%d ppid=%d pgid=%d sid=%d uid=%d username=%q\n",
			i, l.PID, l.PPID, l.PGID, l.SID, l.UID, l.Username)
		fmt.Fprintf(os.Stderr, "           comm=%q tty=%q lstart=%q\n", l.Comm, l.TTY, l.Lstart)
		fmt.Fprintf(os.Stderr, "           command=%q\n", l.Command)
		fmt.Fprintf(os.Stderr, "           raw_ps_line=%q\n", l.RawPsLine)
	}

	if len(b.Layers) > 0 {
		parentComm = b.Layers[0].Comm
	}
	return parentComm
}

// queryChainLayer 查询 pid 的身份快照，按 GOOS 分派到具体实现。
//
// 返回 err != nil 表示查询失败，调用方应终止遍历；但返回的
// ChainLayer 仍填有 PID + RawPsLine（后者含错误文本），不丢信息。
//
// 当前分派：
//   - darwin：queryChainLayerDarwin（ps 拟现有实现）
//   - 其他（windows / linux）：stub，返回明确 "not implemented" 错误
func queryChainLayer(pid int) (ChainLayer, error) {
	switch runtime.GOOS {
	case "darwin":
		return queryChainLayerDarwin(pid)
	default:
		msg := fmt.Sprintf("chain query not implemented on %s yet", runtime.GOOS)
		return ChainLayer{
			PID: pid, PPID: -1, PGID: -1, SID: -1, UID: -1,
			RawPsLine: "<" + msg + ">",
		}, fmt.Errorf("%s", msg)
	}
}

// queryChainLayerDarwin 调一次 ps 查询 pid 的身份快照（macOS）。
//
// 解析失败（fields 不足）时返回填了 RawPsLine 的 layer，结构化
// 字段留哨位 -1；err=nil 表示 "ps 调用成功、只是解析不齐"，
// 调用方可继续遍历（虽然 PPID 为 -1 会导致 cur <= 1 分支）。
//
// ps 列布局：pid / ppid / pgid / sess / uid / tty / user /
// lstart(5 tokens) / comm / command(剩余)。最小字段数 = 13。
//
// command 字段从 raw 里按 token 偏移手动切出，保留原始多空格
// （strings.Fields+Join 会把多空格规整成单空格，丢失原文）。
func queryChainLayerDarwin(pid int) (ChainLayer, error) {
	// CombinedOutput 合并 stderr，失败时 ps 的错误输出（如
	// "ps: no such process"）也能被 RawPsLine 捕捉到。
	out, err := exec.Command("ps",
		"-o", "pid=,ppid=,pgid=,sess=,uid=,tty=,user=,lstart=,comm=,command=",
		"-p", strconv.Itoa(pid)).CombinedOutput()
	if err != nil {
		return ChainLayer{
			PID: pid, PPID: -1, PGID: -1, SID: -1, UID: -1,
			RawPsLine: fmt.Sprintf("<ps err: %v; output=%q>", err, strings.TrimSpace(string(out))),
		}, err
	}
	raw := strings.TrimRight(string(out), "\r\n")
	layer := ChainLayer{
		RawPsLine: raw,
		PID:       pid, // 兜底；解析成功会覆盖
		PPID:      -1,
		PGID:      -1,
		SID:       -1,
		UID:       -1,
	}
	fields := strings.Fields(raw)
	if len(fields) < 13 {
		return layer, nil
	}
	if n, e := strconv.Atoi(fields[0]); e == nil {
		layer.PID = n
	}
	if n, e := strconv.Atoi(fields[1]); e == nil {
		layer.PPID = n
	}
	if n, e := strconv.Atoi(fields[2]); e == nil {
		layer.PGID = n
	}
	if n, e := strconv.Atoi(fields[3]); e == nil {
		layer.SID = n
	}
	if n, e := strconv.Atoi(fields[4]); e == nil {
		layer.UID = n
	}
	layer.TTY = fields[5]
	layer.Username = fields[6]
	layer.Lstart = strings.Join(fields[7:12], " ")
	layer.Comm = fields[12]
	// command 保留原始空白：跳过前 13 个 token 后取原文尾部
	layer.Command = rawAfterNTokens(raw, 13)
	return layer, nil
}

// rawAfterNTokens 跳过 raw 开头的 n 个空白分隔 token 后，返回剩余原文
// （不含引导空白）。与 strings.Fields+Join 的区别：token 内部的多
// 空格不会被规整成单空格。用于保留 command 字段的原始形态。
func rawAfterNTokens(raw string, n int) string {
	s := raw
	for i := 0; i < n; i++ {
		// skip leading whitespace
		j := 0
		for j < len(s) && (s[j] == ' ' || s[j] == '\t') {
			j++
		}
		s = s[j:]
		if s == "" {
			return ""
		}
		// find end of token
		j = 0
		for j < len(s) && s[j] != ' ' && s[j] != '\t' {
			j++
		}
		s = s[j:]
	}
	// trim one leading whitespace run before command’s original content
	j := 0
	for j < len(s) && (s[j] == ' ' || s[j] == '\t') {
		j++
	}
	return s[j:]
}

// envExactMap 固定键 → family 名的精确匹配表。优先级高于 envPrefixRules。
//
// 特异 key 放 exact（如 VSCODE_CRASH_REPORTER_PROCESS_TYPE 归 ExtensionHost，
// 而非让前缀匹配吃进 IDE）；同家族的通用 key 走 prefix。
var envExactMap = map[string]string{
	// ShellPref
	"SHELL": "ShellPref",

	// ShellVersion
	"BASH_VERSION": "ShellVersion",
	"ZSH_VERSION":  "ShellVersion",
	"FISH_VERSION": "ShellVersion",
	"PSModulePath": "ShellVersion",

	// Terminal
	"TERM":                 "Terminal",
	"TERM_PROGRAM":         "Terminal",
	"TERM_PROGRAM_VERSION": "Terminal",
	"COLORTERM":            "Terminal",
	"LC_TERMINAL":          "Terminal",
	"KITTY_WINDOW_ID":      "Terminal",
	"WEZTERM_EXECUTABLE":   "Terminal",
	"ALACRITTY_SOCKET":     "Terminal",
	"TERMINAL_EMULATOR":    "Terminal",

	// Session（比 Terminal 特异：有这些 key 基本能确定是 session manager/remote）
	"SSH_CONNECTION": "Session",
	"SSH_CLIENT":     "Session",
	"SSH_TTY":        "Session",
	"TMUX":           "Session",
	"STY":            "Session",

	// CI
	"CI":             "CI",
	"GITHUB_ACTIONS": "CI",
	"GITLAB_CI":      "CI",
	"CIRCLECI":       "CI",
	"BUILDKITE":      "CI",
	"JENKINS_URL":    "CI",
	"TRAVIS":         "CI",

	// Remote
	"KUBERNETES_SERVICE_HOST": "Remote",
	"IN_NIX_SHELL":            "Remote",

	// ExtensionHost（VSCODE_ 子集，特异性高于 IDE 前缀）
	"VSCODE_CRASH_REPORTER_PROCESS_TYPE": "ExtensionHost",
	"VSCODE_PID":                         "ExtensionHost",
	"VSCODE_IPC_HOOK":                    "ExtensionHost",
	"VSCODE_ESM_ENTRYPOINT":              "ExtensionHost",
	"VSCODE_NLS_CONFIG":                  "ExtensionHost",

	// Locale（LC_TERMINAL 已在上面走 Terminal）
	"LANG": "Locale",

	// Editor
	"EDITOR": "Editor",
	"VISUAL": "Editor",
	"PAGER":  "Editor",

	// Path
	"PATH": "Path",

	// User (Unix)
	"USER":    "User",
	"LOGNAME": "User",
	"HOME":    "User",
	"TMPDIR":  "User",

	// User (Windows; key 存 UPPER 形式，Windows 的 case-insensitive fallback
	// 会命中。以下只放 "用户身份 / 路径" 相关键；机器级 / 系统级
	// 键（COMPUTERNAME / SYSTEMROOT / WINDIR）不归 User，留 Uncategorized。）
	"USERPROFILE":  "User",
	"APPDATA":      "User",
	"LOCALAPPDATA": "User",
	"USERDOMAIN":   "User",
	"TEMP":         "User",
	"TMP":          "User",
}

// envPrefixRules 前缀匹配规则，按顺序逐条试；exact 不命中才查这里。
// 更窄的前缀应在前（目前 IDE / Locale 无冲突）。
var envPrefixRules = []struct {
	prefix string
	family string
}{
	// IDE 前缀（VSCODE_ 的特异 key 已在 envExactMap 归 ExtensionHost）
	{"VSCODE_", "IDE"},
	{"WINDSURF_", "IDE"},
	{"CURSOR_", "IDE"},
	{"CASCADE_", "IDE"},
	{"JETBRAINS_", "IDE"},
	{"IDEA_", "IDE"},
	{"PYCHARM_", "IDE"},

	// Locale（LC_TERMINAL 已在 envExactMap 归 Terminal）
	{"LC_", "Locale"},
}

// classifyEnvKey 将 env key 归类到某个 family；归不到返回空串（→ Uncategorized）。
// exact 优先 prefix；per-key 决策，每个 key 最多一个 family（单归属）。
//
// Windows 适配：Windows 的 env key 大小写不敏感（`Path` == `PATH` == `path`），
// 但 os.Environ() 返回原始大小写。策略：
//  1. 先按原始 key 试匹配（mac 行为此步即终结，与此前完全一致）
//  2. Windows 下再按 UPPER(key) 再试一次（捕捉 Path / Home 等 native 形式）
// mac 不执行第 2 步，原有行为严格不变。
func classifyEnvKey(key string) string {
	// Step 1: 先按原始 key 匹配（mac / Windows 同走此分支）
	if fam, ok := envExactMap[key]; ok {
		return fam
	}
	for _, r := range envPrefixRules {
		if strings.HasPrefix(key, r.prefix) {
			return r.family
		}
	}

	// Step 2: Windows 专用 fallback，按 UPPER(key) 再试一次
	if runtime.GOOS == "windows" {
		upper := strings.ToUpper(key)
		if upper != key { // 已是 UPPER 的 key 第一步已试过，跳过重复查
			if fam, ok := envExactMap[upper]; ok {
				return fam
			}
			for _, r := range envPrefixRules {
				if strings.HasPrefix(upper, r.prefix) {
					return r.family
				}
			}
		}
	}

	return ""
}

// mapForFamily 返回 EnvBlock 里指定 family 对应的 map 指针。
// family 名必须与 envExactMap / envPrefixRules 的 family 字段一致；
// 不识别的返回 Uncategorized 兜底。
func (b *EnvBlock) mapForFamily(family string) *map[string]string {
	switch family {
	case "ExtensionHost":
		return &b.ExtensionHost
	case "IDE":
		return &b.IDE
	case "Session":
		return &b.Session
	case "CI":
		return &b.CI
	case "Remote":
		return &b.Remote
	case "Terminal":
		return &b.Terminal
	case "ShellVersion":
		return &b.ShellVersion
	case "Editor":
		return &b.Editor
	case "Locale":
		return &b.Locale
	case "Path":
		return &b.Path
	case "User":
		return &b.User
	case "ShellPref":
		return &b.ShellPref
	default:
		return &b.Uncategorized
	}
}

// dumpEnvBlock 采集 env 块并写入 stderr。
//
// 策略：
//   - Raw 保留 os.Environ() 原样（顺序 + 重复），不丢信息的第一道防线
//   - 每个 env 走 classifyEnvKey 单归属到 family，归不到进 Uncategorized
//   - stderr 按 family 顺序打（空 family 跳过 header），uncategorized 之后
//     打全量 key 清单（按字母排序），便于扫视和跨样本 diff
func dumpEnvBlock() {
	b := EnvBlock{
		Raw: os.Environ(), // 全量原样
	}

	allKeys := make([]string, 0, len(b.Raw))
	for _, kv := range b.Raw {
		var key, val string
		if eq := strings.IndexByte(kv, '='); eq < 0 {
			// 罕见：env 里没有 = 的条目（不合规但理论可能）
			key = kv
			val = ""
		} else {
			key = kv[:eq]
			val = kv[eq+1:]
		}
		allKeys = append(allKeys, key)

		target := b.mapForFamily(classifyEnvKey(key))
		if *target == nil {
			*target = make(map[string]string)
		}
		(*target)[key] = val
	}

	// stderr 输出（按 family 顺序；空 family 跳过 header）
	fmt.Fprintf(os.Stderr, "[HASSHIN] ---- Env ----\n")
	printEnvFamily("extension_host", b.ExtensionHost)
	printEnvFamily("ide", b.IDE)
	printEnvFamily("session", b.Session)
	printEnvFamily("ci", b.CI)
	printEnvFamily("remote", b.Remote)
	printEnvFamily("terminal", b.Terminal)
	printEnvFamily("shell_version", b.ShellVersion)
	printEnvFamily("editor", b.Editor)
	printEnvFamily("locale", b.Locale)
	printEnvFamily("path", b.Path)
	printEnvFamily("user", b.User)
	printEnvFamily("shell_pref", b.ShellPref)
	printEnvFamily("uncategorized", b.Uncategorized)

	// 全量 key 清单，按字母排序（便于扫视和跨样本 diff）
	sortedKeys := append([]string(nil), allKeys...)
	sort.Strings(sortedKeys)
	fmt.Fprintf(os.Stderr, "  all_keys (N=%d, sorted):\n    %s\n",
		len(sortedKeys), strings.Join(sortedKeys, " "))
}

// printEnvFamily 打一个 family 块到 stderr；空 map 跳过不打 header。
// family 内 key 按字母排序，便于阅读和跨样本 diff。
func printEnvFamily(label string, m map[string]string) {
	if len(m) == 0 {
		return
	}
	fmt.Fprintf(os.Stderr, "  [%s]\n", label)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(os.Stderr, "    %s=%q\n", k, m[k])
	}
}

