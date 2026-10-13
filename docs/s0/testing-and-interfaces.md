# S0 — Current CLI test seams and verification boundaries

**Status:** Forward-looking revised plan; not passing fixture evidence
**Authority:** [architecture](../architecture.md#11-code-boundaries),
[current traceability](cli-traceability.md), [plan](../implementation-plan.md)

## Planned package seams (minimal contracts; no premature API freeze)

| Responsibility | Test boundary |
| --- | --- |
| CLI | Explicit arguments, selection, human/JSON rendering and exit semantics |
| Run model | Lifecycle/outcome aggregation, tested-input scope, no vacuous PASS |
| Workflow/source | Repository discovery, payloads/current files, identity/mutation/path safety |
| Engine | Tool/config probes, verified artifact handoff, argv and process cancellation |
| Reports | Atomic per-run files, redaction, concurrent repository/worktree isolation |

Create seams only as their consumer needs them. No GitGate, approval provider,
hook installer or reusable-PASS store belongs in the new feature.

## Planned commands (after `go.mod` exists)

The Go module already exists. Implementation checks:
`gofmt -l cmd internal`, focused `go test`, `go vet ./...`, and appropriate builds.
Run the existing full deterministic suite when practical/required by CI.
Report race/coverage outcomes only when actually executed; no invented threshold.

Temporary repositories and fake engines provide fast local regression evidence.
They never use real secrets, global Git configuration, developer repositories,
Docker-wide cleanup or public network probes.
Process, Docker, cache, cancellation and permission acceptance on real supported
hosts is separate, authorized work; fake fixtures cannot certify it.

## S1 entry/exit boundary

The merged foundation provides help/version and safe unavailable placeholders.
Keep that base, update legacy wording and retire gate-only code in a separate PR.
Resolve CLI-01 before command implementation; SRC-01 before source execution;
DEP-01/SEC-01 before trusted dependency/runtime integration; D4 before release.

Implementation sessions perform related verification and report what remains.
Independent reviews are separate gates, not self-review claims. Existing source
is not rewritten or removed by this planning change.
