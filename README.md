# CI See

Local-first GitHub Actions execution with developer-controlled CI gating.

## Documentation and AI agents

Use [AGENTS.md](AGENTS.md) and the [documentation index](docs/index.md) to select
relevant project contracts. The repository-local
[Context Engineering skill](.agents/skills/context-engineering/SKILL.md) is installed
from [ai-skills PR #6](https://github.com/arvin-ramezani/ai-skills/pull/6)
(merge commit `6ec528e498b51a4e1a9796bb55469bc188e5fdc1`). Its [local policy](docs/context-policy.md)
requires meaning-preserving documentation changes; it does not change product scope.

Read the [PRD](docs/prd.md), [architecture](docs/architecture.md),
[three behavior specs](docs/index.md), or [implementation plan](docs/implementation-plan.md)
through their task-specific indexes. [Migration evidence](docs/migration/index.md)
records the unchanged source sections and legacy anchors.

To validate docs locally (Python standard library only):

```bash
python3 scripts/validate_docs.py --check
```
