# CI See — Agent Instructions

Use the [documentation index](docs/index.md) to locate the relevant contract;
load only the topic sections needed for the task. The public contract paths are
compatibility indexes, with authoritative wording in task-scoped child files.
Follow [context policy](docs/context-policy.md) when modifying documentation.

The PRD, architecture, local execution spec and implementation plan now describe
explicit CLI verification. Git-gating/approval contracts and previous S0/D3
procedures are historical. Read the [alignment record](docs/migration/cli-contract-alignment.md)
for preserved evidence, open decisions and the separate cleanup boundary.
This documentation change does not authorize live act or claim runtime acceptance.

## Authority and workflow

- [PRD](docs/prd.md) owns product scope and requirements.
- [Architecture](docs/architecture.md) owns the Go MVP's system boundaries.
- [Local CI execution](docs/specs/local-ci-execution.md) owns current runner behavior.
  [Git gating](docs/specs/git-gating.md) and [approval UX](docs/specs/developer-approval-ux.md)
  are historical only; do not implement their removed requirements.
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

Check links, fragments, unique authoritative IDs, and index reachability:

```bash
python .agents/skills/context-engineering/scripts/inspect_docs.py . --docs docs --entry docs/index.md --check
```

The scanner does not check fragment anchors, nested Markdown, semantic preservation,
or implementation correctness. Python is not a CI See Go runtime requirement.

Additional contract/legacy-anchor checks:

```bash
python3 scripts/validate_docs.py --check
```

Neither command proves semantic equivalence or that runtime tests pass.
