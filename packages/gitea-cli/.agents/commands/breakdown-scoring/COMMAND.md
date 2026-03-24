---
name: breakdown-scoring
metadata:
  version: 0.0.3
---

# Breakdown Scoring

## Scoring Criteria

按照以下标准对指定文件进行评审和评分，将评分附加到原文底部：

1. Comprehension - 需求理解 - 10分制
    - 对需求的理解正确性、全面性
    - 是否准确识别了所有关键需求点
    - 是否理解了需求的优先级和约束条件

2. Scope - 范围控制 - 10分制
    - 保持在需求指定的范围内工作，没有偏离方向
    - 是否避免了超出需求的额外功能开发
    - 是否明确标识了超出范围的工作（如有）

3. Skills - 技能应用 - 10分制
    - 对 Skills 使用的正确性、全面性
    - 是否正确使用了所需的技能和工具
    - 是否充分利用了可用的技能来解决问题

4. Format - 格式规范 - 10分制
    - 输出产物严格符合格式要求
    - 是否遵循了指定的输出格式标准
    - 代码风格、文档结构是否规范

5. Quality - 综合质量 - 10分制
    - 综合整体工作水平
    - 输出产物的完整性和可用性
    - 是否存在明显的缺陷或遗漏

## Output Example

```yaml
### Breakdown Scoring

- Comprehension: a
- Scope: b
- Skills: c
- Format: d
- Quality: e

- Final: (a + b + c + d + e) / 5
```
