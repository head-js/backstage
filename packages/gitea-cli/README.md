# Backstage Gitea Cli

## Workflow

### Create Plan

```bash
$ backstage-gitea plan POST /cms-mgr/plans --title "PLAN-102: Upload Image"
```

### Breakdown Plan to Phase

`.agents/commands/breakdown-plan-to-phase`

```bash
# For-Each Phase
$ backstage-gitea plan POST /cms-mgr/plans --title "PLAN-102: Upload Image"
```

### Breakdown Phase to Task

`.agents/commands/breakdown-phase-to-task`

```bash
# For-Each Phase
$ backstage-gitea plan POST /cms-mgr/PLAN-102/PHASE-200/tasks --name "Design & Develop"
```
