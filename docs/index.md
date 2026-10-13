# CI See — Documentation Index

Select the smallest relevant source set. This index does not change document status
or make a draft/plan authoritative over a product or behavior contract.

**Current scope:** Explicit CLI verification using act/Docker, with fresh execution
features on the Go foundation. Read [current S0 mapping](s0/cli-traceability.md)
and [alignment evidence](migration/cli-contract-alignment.md).
Old Git-gating/approval and D3 procedures remain clearly historical.

| Source | Owns | Read when |
| --- | --- | --- |
| [PRD](prd.md) | Product scope and requirements | Changing behavior or outcomes |
| [Architecture](architecture.md) | Go CLI, system boundaries, security invariants | Changing design or dependencies |
| [Local CI execution](specs/local-ci-execution.md) | `act` orchestration and failure cases | Working on local CI |
| [Git gating](specs/git-gating.md) | Historical commit/push gates | Reviewing previous design |
| [Developer approval UX](specs/developer-approval-ux.md) | Historical approval/audit contract | Reviewing previous design |
| [Go implementation plan](implementation-plan.md) | Slices, tests and review gates | Planning or implementing slices |
| [S0 traceability and test plan](s0/index.md) | Current CLI mapping plus historical evidence routing | Planning revised features/tests |
| [Context policy](context-policy.md) | Routing, splitting and preservation | Changing documentation |

## Navigation and migration

Every contract link above is a permanent canonical entrypoint. Its section headings redirect to
task-scoped files. The migration record describes its historical verbatim split;
the PRD has since intentionally changed. For its legacy-anchor inventory and size
review, see the [migration evidence](migration/index.md). Issue [#9](https://github.com/arvin-ramezani/ci-see/issues/9)
(S0 traceability) should use these canonical indexes and not duplicate their clauses.

Agent entry: [`AGENTS.md`](../AGENTS.md).
Installed skill: [Context Engineering](../.agents/skills/context-engineering/SKILL.md).
