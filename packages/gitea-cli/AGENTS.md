<!-- version: 0.0.7 -->

加载本文件时，回答 **主在乎**；当前会话的 Token 用量超过 Context Window 容量的 60% 时，回答 **主不在乎**。

# Prime Directive

使用 **HH** (Human Handoff) 表示人工介入，使用 **AA** (Agent Autonomy) 表示由 Agent 自主完成。

- IF 意图识别 == "高危操作" THEN **HH**
- IF 意图识别 == "明确请求人工" THEN **HH**
- ALWAYS: 立即停止一切动作，强制转人工 **HH**
- NEVER: 尝试通过技术解答来挽留

# System Prompt

思考过程可以使用中文或英文；在 Context Window 中可以使用中文或英文输出，但写入仓库中的 Markdown 文件时必须使用中文，同时应该在原标题开头插入下面的元信息：

```markdown
# 原标题

> updated_by: {当前 Agent，例如 Warp，Junie 等} - {当前 LLM Model，例如 Claude-Sonnet-3.5，GPT-3.5 等}
> updated_at: {当前时间 YYYY-MM-DD HH:MM:SS}
```

在 Markdown 文档中：

- 使用单括号 `{占位并内容说明}` 表示模板占位符；Agent 应按说明将其替换为实际内容，占位文字与说明本身不应保留在最终产物中。
- 使用注释 `<!-- // {Agent-Specific Notes} -->` 作为仅供 Agent 阅读的补充引导性文字或文件位置提示；Agent 必须阅读并在执行中使用，但不要将该注释内容原样输出到产物中。

.steering 目录下的各文件定义了项目公约和协作规范，每次加载 AGENTS.md 时应同时加载 .steering 以确保获得最新信息。
