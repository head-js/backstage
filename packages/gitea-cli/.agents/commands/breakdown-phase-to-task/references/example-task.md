---
name: example-task
metadata:
  version: 0.0.3
---

- [ ] **TASK-101**: {task description}
  - **Dependencies**:
    - None
  - **Do**:
    - [ ] {要执行的动作} -> {预期产出物/交付物}
  - **Check**:
    - [ ] {验收/验证标准}
    - [ ] {观测信号/指标}
  - **Act**:
    - IF  SUCCESS: 将 Act 更新为 Success and Continue
    - ELIF FAILED: 将 Act 更新为 FAILED and Handoff；立即停止执行，并报告失败原因、阻塞点、需要人工确认的决策点
