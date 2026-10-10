# S0 — Go testing, CI environments and interfaces

**Status:** Forward-looking S1+ test plan, not an implementation or passing test report.
**Authority:** [Architecture code boundaries](../architecture.md#11-code-boundaries), [implementation slices](../implementation-plan/vertical-slices.md#4-dependency-ordered-vertical-slices), [fixture catalog](fixtures.md).

## Planned package seams (minimal contracts; no premature API freeze)

| Owner | Proposed boundary to test | Constraint |
| --- | --- | --- |
| `internal/core` | Typed execution states (`PASS`, `FAIL`, `INCOMPLETE`, `NOT_RUN`, `STALE`) and gate outcomes; pure fingerprint and allow/block policies | Must not import process, Git or OS adapters |
| `internal/adapters/git` | Repository/worktree identity, index tree, commit-tip/ref context, snapshot materialization, hook inspection | Git-provided paths, no mutable working-tree leakage |
| `internal/adapters/act` | Dependency probe, event/secret transport, process run, job/overall normalization and cancellation | `os/exec` argv, never shell command concatenation |
| `internal/adapters/state` | Read/write exact validation run, atomic durable storage, lock/lease recovery, private fingerprint inputs | No cross-worktree reusable PASS or secret plaintext |
| `internal/adapters/approval` | Human-provider feasibility interface; authenticated single-use approval response | S6 trust evidence blocks S7 bypass; no fake provider shipping |
| `internal/hooks` | Hook installation/chain, gate invocation, unambiguous exit outcomes | Preserve existing hooks or reject unsupported setup |
| `internal/cli`, `cmd/ci-see` | Commands `ci-see`, `init`, `status`, diagnostics and exit semantics | Wrap use cases; no duplicated product rules |

Interface suggestions, not final Go types: `Snapshot{Repo,Worktree,Tree,Commit?,EventContext}`, `ValidationRequest{Snapshot,Fingerprint,Workflows}`, `RunResult{State,Jobs,EngineVersion}`, `GateDecision{Outcome,Reason}`. Choose actual fields and types in S1/S2 after D3 review; secrets should not be serialized into ordinary results.

## Planned commands (after `go.mod` exists)

```sh
gofmt -l ./cmd ./internal           # empty output expected; format changed files
go test ./...
go vet ./...
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -race ./...                 # when Go race toolchain/CGO support is available
```

- S1: pin the supported Go toolchain in `go.mod` and CI, choose test package naming, establish a **measured baseline**. Proposed quality goal: comprehensive tests for every fail-closed core branch; review an **~80% core-package line coverage target** only after S1 metrics. Do not impose a fabricated repository-wide threshold on this docs-only S0.
- Unit tests: pure core and stubbed adapters. Integration: real isolated Git repos + fake `act` executable, no Docker/network. Acceptance: explicitly provision real `act`+Docker, Git and representative workflows; do not infer hosted parity.
- CI after S1: Linux runner compiles/tests/vets and runs fake-`act` fixtures by default; Windows-native tests use a **Windows binary and native Git**, while WSL2 tests run a **Linux binary inside WSL2** as a separately recorded environment. macOS build/smoke and actual UI approval feasibility have dedicated host gates.
- Race checks on supported hosts; platform-specific file locking, process-tree cancellation, hook interpreter/quoting, restrictive secret permissions, Git path semantics, and filesystem rename/sync require real-host acceptance in S3–S10.
- Fixture code must not touch developer repositories or real secrets. Use per-test temp root and no shared global Git hooks/config. Test cleanup under success, panic/crash and timeout conditions.
- S0 validation consists **only** of documentation integrity, matrix cardinalities, fixture-reference consistency and independent contract review. `go test`/coverage/race cannot pass until there is executable Go code.

## S1 entry/exit boundary

Begin only after reviewed S0 mapping and documented D3 open points; D3 must be resolved before S2/S4 where its choices affect safety. S1 exit: reproducible native Go CLI skeleton and typed state/gate outcomes, adapter seams with tests, no false PASS for missing inputs, with exact commit/SHA and actual command output. Do not implement Git gating, `act` execution, approval or remote CI in S1.
