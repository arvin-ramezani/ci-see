# System architecture — Small Go packages and unresolved contracts

**Status:** Draft aligned with merged PRD #20
**Canonical authority:** [System architecture](../architecture.md)

## 11. Code Boundaries

Keep `cmd/ci-see`, `internal/cli`, Go module/build metadata and useful CLI tests.
Build fresh run, workflow, engine and report behavior in focused internal packages.
Suggested responsibilities, not mandatory empty scaffolding:

| Responsibility | Owns |
| --- | --- |
| CLI | Parse inputs and render human/JSON output |
| Run model | Outcomes, scope, lifecycle and input/report identities |
| Workflow/source | Discovery, event context and tested-file preparation |
| Engine | Version checks, artifact handoff, argv process execution/cancellation |
| Reports | Per-run storage, redaction and atomic finalization |

Interfaces belong with consumers. Pure result rules must not import process,
Git or Docker adapters. Use explicit boundary validation and typed errors.
Create adapters only when needed; do not add a daemon, framework or cache service.

The existing `internal/core` is a legacy gate skeleton, not the new run model.
`cmd/d3-preflight` and `internal/d3preflight` are separate offline test tooling,
not prerequisites or implementation of local workflow execution.
Removing/relocating that code and obsolete tests belongs in a separate cleanup PR.
This documentation PR removes no executable code.

## 12. Deferred Decisions

Resolve in a focused spec/ADR before dependent implementation:

- CLI syntax, numeric exits and versioned JSON schema.
- Portable current-file capture, ignored/untracked policy, Git metadata needs,
  symlink/submodule handling and detection of input changes.
- Explicit engine configuration isolation and verified cache handoff capabilities.
- Tool versions, runner mappings, timeout/cleanup mechanisms and host support matrix.
- Secret input methods and acquisition-versus-runtime networking policy.

Do not silently choose these security/compatibility decisions during unrelated work.
Hosted suppression, native Windows/macOS jobs and other OSs remain outside MVP.

## 13. Next Specifications

Use the revised [local execution contract](../specs/local-ci-execution.md) and
[implementation plan](../implementation-plan.md). The gating and approval specs
are historical; do not create new implementations from them.

Fresh discovery/execution/report features should reuse the small CLI foundation
rather than evolve block-only gates. No new application behavior is claimed here.
