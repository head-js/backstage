# Linear 批量上传 Blame 操作记录

> updated_by: Kilo - GLM-5.1
> updated_at: 2026-06-13 12:37:00

## 背景

将 Gitea 上累积的 blame 记录批量同步到 Linear 的 `LSTD/Blame` project，便于在 Linear 侧统一查看与跟进。

## 前置条件

1. **工具就位**
   - `backstage-gitea` 已在全局 PATH（用于读取 blame 列表）
   - `backstage-linear` 已在全局 PATH（用于上传到 Linear）
   - `python` 可用（用于解析 JSON + 串行循环）

2. **Linear 环境变量**（`backstage-linear` 依赖）
   - `BACKSTAGE_LINEAR_URL`
   - `BACKSTAGE_LINEAR_TOKEN`

3. **Linear 目标**：`LSTD` team 下存在名为 `Blame` 的 project。

## 操作步骤

### Step 1：拉取 Gitea 全部 blame 列表，落盘到临时文件

```bash
backstage-gitea plan GET /blames > "C:\Users\shuai\AppData\Local\Temp\kilo\blames.json"
```

返回的是 JSON 数组，每条结构：

```json
{
  "id": "BLAME-007",
  "name": "[FALLBACK] xxx",
  "title": "BLAME-007: [FALLBACK] xxx",
  "gitea": {
    "type": "ISSUE",
    "id": 248,
    "no": 13,
    "owner": "blame",
    "repo": "PLAN-000",
    "title": "BLAME-007: [FALLBACK] xxx",
    "description": "",
    "body": "### Category\n产物存疑\n\n### Where\n...",
    "state": "open"
  }
}
```

验证条数：

```bash
python -c "import json; arr=json.load(open(r'C:\Users\shuai\AppData\Local\Temp\kilo\blames.json',encoding='utf-8')); print('total:', len(arr))"
```

### Step 2：写一个一次性 Python 脚本，串行上传

脚本位置：`C:\Users\shuai\AppData\Local\Temp\kilo\upload_blames.py`

```python
import json, subprocess, time, sys

BLAMES_FILE = r"C:\Users\shuai\AppData\Local\Temp\kilo\blames.json"
arr = json.load(open(BLAMES_FILE, encoding='utf-8'))
print(f"[start] total={len(arr)}", flush=True)

ok, fail = 0, 0
for i, item in enumerate(arr, 1):
    title = item["title"]                  # 顶层 title，带 BLAME ID 前缀
    body = item["gitea"]["body"]           # Gitea issue 完整 Markdown 正文
    context = "." + body                   # 必须以 "." 开头，绕过 LoadContext 前置检查

    print(f"[{i}/{len(arr)}] uploading {item['id']!r}  title_len={len(title)}  ctx_len={len(context)}", flush=True)
    proc = subprocess.run(
        ["backstage-linear", "api", "POST",
         "/teams/LSTD/projects/Blame/issues",
         "--title", title,
         "--context", context],
        capture_output=True, text=True, encoding='utf-8'
    )
    out = (proc.stdout or "").strip()
    err = (proc.stderr or "").strip()
    if proc.returncode == 0 and out and not out.startswith("Error"):
        print(f"  -> OK  {out[:200]}", flush=True)
        ok += 1
    else:
        print(f"  -> FAIL rc={proc.returncode}  stdout={out[:200]}  stderr={err[:200]}", flush=True)
        fail += 1
    time.sleep(1)                          # 每条之间 sleep 1s，节流

print(f"[done] ok={ok} fail={fail} total={len(arr)}", flush=True)
sys.exit(0 if fail == 0 else 1)
```

执行：

```bash
python "C:\Users\shuai\AppData\Local\Temp\kilo\upload_blames.py"
```

## 关键约定

1. **title**：用顶层 `title`（带 BLAME ID 前缀），不要用 `name`，便于跨系统追溯。
2. **context 前导点**：`backstage-linear` 的 `--context` 参数走 `framework.LoadContext`，该函数前置检查 `strings.HasPrefix(context, ".")`，不以 `.` 开头会报 `InvalidFormatException`。当传入的不是 `.md` 路径而是字面量 Markdown 时，需要在内容前手动拼一个 `.` 字符绕过该检查。这是已知契约缺陷，已通过 blame 上报，暂不修复。
3. **副作用**：上送到 Linear 的 description 会带一个前导 `.` 字符，已知且接受。
4. **串行 + sleep 1s**：避免对 Linear 接口产生突发压力，也便于按行排查失败。
5. **不去重**：源数据中 BLAME ID 可能重复（不同 PLAN 下分别独立编号），按原样上传，Linear 侧用自己的 `LSTD-xx` identifier 区分。
6. **投递目标**：统一 `/teams/LSTD/projects/Blame/issues`，路径用 team key + project name，UUID 解析由 `backstage-linear` 内部完成。

## 复用建议

- 后续若再次批量同步，先比对一下 Linear 现有 issue 列表（`backstage-linear api GET /teams/LSTD/projects/Blame/issues`），手动决定是否需要去重；当前脚本默认不做去重判定。
- 临时脚本与数据放在 `C:\Users\shuai\AppData\Local\Temp\kilo\` 下，跑完可清理。
