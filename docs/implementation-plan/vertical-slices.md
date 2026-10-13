# Implementation planning — Explicit CLI slices and retired IDs

**Status:** Draft revision aligned with merged PRD #20
**Canonical authority:** [Implementation planning](../implementation-plan.md)

## 4. Dependency-Ordered Vertical Slices

IDs remain stable; changed meanings are explicit below. Old slice numbers in
historical PRs/evidence do not prove the newly defined slice is implemented.

| Slice | Revised deliverable | Current status and required evidence |
| --- | --- | --- |
| **S0 — Contract alignment** | PRD, architecture, execution spec, active test map and retirement routing | This documentation PR; review pending, no runtime acceptance |
| **S1 — CLI foundation and scoped cleanup** | Keep Go module/entrypoint/help/version; retire unused gate models, remove gate help text and isolate/retire D3 tooling in a separate PR | Foundation merged via PR #15; cleanup and revised command contract pending |
| **S2 — Repository, discovery and current source** | Fresh repo/workflow/job discovery, explicit/default events and tested local input capture; resolve CLI-01/SRC-01 | Not implemented; fake execution and temporary repositories first |
| **S3 — Dependencies and configuration** | Tool preflight, trusted commit/digest resolution, acquisition and verified reuse through existing caches; explicit config/secrets, resolve DEP-01/SEC-01 | Not implemented; cold/warm/tamper/failure fixtures; no new cache engine |
| **S4 — Workflow execution** | Fresh act adapter, argv, bounded lifecycle, complete selection/job outcomes and run-owned cleanup | Not implemented; fake process verification then separately authorized real-tool acceptance |
| **S5 — Reports and agent consumption** | Atomic per-run reports, human logs, redaction, stable exits/versioned JSON and no historical PASS reuse | Not implemented; useful report primitives may be built with earlier slices |
| **S6 — Approval-provider feasibility** | RETIRED: human-only provider is removed from product | No implementation or acceptance dependency |
| **S7 — Developer bypass and audit** | RETIRED: approval/bypass/audit-before-Git is removed | No implementation or acceptance dependency |
| **S8 — Hosted duplicate control** | DEFERRED beyond MVP under FR-10 | No MVP release dependency; existing GitHub triggers remain |
| **S9 — Release and host acceptance** | Native packages, checksums/version metadata, installation/uninstall and tested support matrix | Not implemented; existing Linux/Windows foundation CI is not full host acceptance |
| **S10 — Explicit CLI MVP acceptance** | Install → explicitly run selected existing workflows → read reliable result → independently use Git | Not implemented; supported hosts, cold/warm caches, failure/recovery and honest limitations |

**Active sequence:** S0 → S1 cleanup/contract → S2 → S3 → S4 → S5 → S9 → S10.
S6/S7 are retired and S8 deferred; none blocks this sequence.
Report/outcome primitives can accompany S2–S4 when needed to verify their behavior.

Fresh discovery/execution/report packages reuse the small Go foundation.
Do not refactor the legacy gate authorization framework into the new execution path.

## Slice clarification — 2026-10-12

**Historical heading retained.** The earlier staged-snapshot/S2 and S4 sequencing
clarification is superseded by this table. Its original wording is recoverable at
baseline `63ef75bf5a7c475675e0602a79df0294206d529a`.
D3 evidence remains PARTIAL; no new runtime execution is authorized here.
