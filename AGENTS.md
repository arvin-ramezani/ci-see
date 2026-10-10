# CI See — Agent Instructions

Start with [documentation index](docs/index.md). Load only task-specific canonical
context, and follow [context policy](docs/context-policy.md) when modifying docs.

## Authority and workflow

- [PRD](docs/prd.md) owns product scope and requirements.
- [Architecture](docs/architecture.md) owns the Go MVP's system boundaries.
- [Local CI execution](docs/specs/local-ci-execution.md), [Git gating](docs/specs/git-gating.md),
  and [developer approval UX](docs/specs/developer-approval-ux.md) own their behaviors.
- [Implementation plan](docs/implementation-plan.md) owns implementation sequencing,
  not new authority to change behavioral contracts.
- Honor the statuses of these sources. Preserve disputed evidence, escalate genuine
  owner decisions, and continue independently authorized work.
- Validate changes and obtain independent reviews per implementation plan.
  Do not merge any PR without explicit owner approval.

## Documentation maintenance

Use the installed [`$context-engineering`](.agents/skills/context-engineering/SKILL.md)
to review document growth, routing, splitting/moving, and preservation.
For authorized structural updates, do justified work without a file-count gate,
but never silently change accepted meaning or application behavior.
Capture a baseline, map every affected substantive unit, repair navigation, and
compare original to new meaning. Review over-150-line files; do not split mechanically.

Check common Markdown links and index reachability:

```bash
python .agents/skills/context-engineering/scripts/inspect_docs.py . --docs docs --entry docs/index.md --check
```

The scanner does not check fragment anchors, nested Markdown, semantic preservation,
or implementation correctness. Python is not a CI See Go runtime requirement.
