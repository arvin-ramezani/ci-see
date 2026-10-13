# System architecture — Explicit CLI, engine boundaries, reports

**Status:** Draft aligned with merged PRD #20; not runtime acceptance
**Canonical authority:** [System architecture](../architecture.md)

## 1. Purpose

Deliver a lightweight Go CLI for explicitly running existing GitHub Actions
locally through user-installed `act` and Docker. Developers own invocation.
Ordinary Git operates independently; there are no gates or approval components.

## 2. Architecture Goals

- One workflow definition; delegate execution to `act`.
- Minimal setup and no mandatory workflow rewrites.
- Download missing dependencies and reuse verified compatible engine/runtime caches.
- Accurate outcomes, tested-input identity, logs, exit semantics and optional JSON.
- No daemon, GUI, custom cache engine, Git interception or hosted suppression.
- Windows, macOS and Linux hosts; initially Linux-container jobs.

## 3. MVP Decisions

| Area | Decision |
| --- | --- |
| Language | Go executable; standard library first |
| Engine/runtime | External `act` and Docker with Linux containers |
| Invocation | Explicit CLI only; no automatic Git integration |
| Source | Current local files, including uncommitted changes; not staged-index gating |
| Persistence | Per-run local reports, no reusable-PASS authorization store |
| Hosts | Windows-native, Linux, macOS; WSL2 is a distinct Linux environment |
| Packaging | Prebuilt native binaries; host Go/Node unnecessary for CI See itself |

CI See owns selection, setup checks, verified-artifact handoff, process lifecycle
and result reporting. `act` owns workflow semantics; Docker owns containers/images.
Package managers own workflow package caches and lockfile integrity.

## 4. System Context

A developer/agent invokes the CLI. CI See reads local workflows and Git context,
checks required tools, prepares tested inputs, and invokes `act` with argv.
`act` uses Docker to execute the selected Linux-compatible jobs.
Logs and reports return to the caller. Git commit/push are outside this path.

No GitHub registration, hosted status publication or agent-specific service is required.

## 5. Components

### CLI

Own command parsing, help/version, explicit invocation, setup guidance, selection,
terminal output and JSON mode. `ci-see init` is optional local setup only.
Exact new flags/subcommands remain a focused command-contract decision; existing
S1 placeholders are not implemented feature contracts.

### Core

A small run model owns PASS/FAIL/INCOMPLETE, lifecycle, selected scope and
report accuracy. It has no GateDecision, ALLOW_BYPASS, hook or human-provider model.
Do not adapt the existing block-only gate core into the new runner.

### Adapters

- Git/context: resolve repository/worktree metadata without altering Git settings.
- Workflow/source: discover YAML and prepare/identify actual tested inputs.
- Engine/process: preflight, explicit configuration, argv execution and cancellation.
- Reports: atomic per-run metadata/log storage with repository isolation.

### Git hooks

**Removed.** This heading remains for old links. Neither init nor execution
installs, chains, modifies, invokes or depends on CI See Git hooks.

## 6. Local State

Use Git-resolved worktree metadata for private run storage, not a literal `.git/`
assumption. Conceptual location: `<git-dir>/ci-see/runs/<run-id>/`.
Keep run IDs unique, writes atomic, secrets separate and permissions restrictive.

Lifecycle: RUNNING becomes PASS, FAIL or INCOMPLETE only after evidence is collected.
An abandoned RUNNING record is incomplete; a crash never fabricates PASS.
A failed final report write is non-PASS. Parallel runs cannot overwrite one another.
No repository-global lock is required merely for independent run directories;
shared cache coordination belongs to existing engine/runtime facilities.

Reports are historical evidence, never authorization to commit/push or skip a new run.
