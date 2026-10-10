# CI See — Documentation Index

Select the smallest relevant source set. This index does not change document status
or make a draft/plan authoritative over a product or behavior contract.

| Source | Owns | Read when |
| --- | --- | --- |
| [PRD](prd.md) | Product scope and requirements | Changing behavior or outcomes |
| [Architecture](architecture.md) | Go CLI, system boundaries, security invariants | Changing design or dependencies |
| [Local CI execution](specs/local-ci-execution.md) | `act` orchestration and failure cases | Working on local CI |
| [Git gating](specs/git-gating.md) | Commit/push exact-state gates | Working on Git integration |
| [Developer approval UX](specs/developer-approval-ux.md) | Bypass decisions and audit rules | Working on approvals |
| [Go implementation plan](implementation-plan.md) | Slices, tests and review gates | Planning or implementing slices |
| [S0 traceability and test plan](s0/index.md) | Requirement mapping, D3 proposals, fixture planning | Working on S0 or planning S1 tests |
| [Context policy](context-policy.md) | Routing, splitting and preservation | Changing documentation |

## Navigation and migration

Every contract link above is a permanent canonical entrypoint. Its section headings redirect to
verbatim task-scoped files. For the source-to-destination/legacy-anchor inventory and size
review, see the [migration evidence](migration/index.md). Issue [#9](https://github.com/arvin-ramezani/ci-see/issues/9)
(S0 traceability) should use these canonical indexes and not duplicate their clauses.

Agent entry: [`AGENTS.md`](../AGENTS.md).
Installed skill: [Context Engineering](../.agents/skills/context-engineering/SKILL.md).
