# S0 — Deterministic fixture catalog (design only)

**Status:** Test specifications, **not created executable fixtures**.
**Authority:** [Local CI](../specs/local-ci-execution.md), [Git gating](../specs/git-gating.md), [approval](../specs/developer-approval-ux.md).
**Harness:** Isolated `t.TempDir()` repositories; fixed author/committer/time/config; no network; stub `act` executable found on per-test `PATH`. Real `act`/Docker runs are separate opt-in integration tests.

## Test harness invariants

- Initialize real Git objects in each repository. Disable signing and unrelated global hooks; use deterministic author/email/date, `-c` overrides and platform-safe paths. Reset only the test-owned temp directory.
- Stub `act` as a controlled executable that records argv, environment **names** (not secret values), current directory and file paths; emits scripted job output/exit/status; can sleep, spawn children, be killed and fail. No shell-string construction by Go.
- For event/payload tests, compare normalized fields and selected trees; scrub temporary paths/timestamps or fix clock/IDs. Assert absence of untracked/unstaged file bytes, not merely matching tree hashes.
- For every non-PASS fixture, assert both normalized status and Git gate's **non-zero** result; assert no reusable PASS, no secret leakage, and cleanup. Unsupported required behavior is not treated as a successful skip.
- Test parallel runs with barriers/channels and unique directories, then inject interrupted/partial writes; restore permissions and stop descendants in cleanup. Never require a real developer to approve a test.

## Fixture IDs, setup/stimulus and expected oracle

| ID | Deterministic scenario; key assertion |
| --- | --- |
| **R01** | Temp repo/not-repo + repeated init, configured hooks and missing dependencies; clear repository/dependency outcomes, no unexpected workflow edits. |
| **W01** | Nested `.yml`/`.yaml`, multi-event workflow, no manual registration; discover eligible workflows from the **snapshot**, each once. |
| **W02** | Zero executable, unsupported-only, malformed YAML, required dispatch inputs; explicit unsupported/INCOMPLETE, **not PASS**. |
| **EV01** | PR+push+dispatch workflow, push-only, dispatch-only; deterministic PR → push → dispatch precedence. |
| **EV02** | Known branch/default/base/head vs missing/ambiguous inputs, detached HEAD; generated truthful temporary payload or INCOMPLETE; no user JSON needed. |
| **SN01** | Stage file A, leave conflicting unstaged A + untracked B; snapshot matches `git write-tree`, neither unstaged nor untracked bytes reach `act`. |
| **SN02** | Post-commit same/different tree, `-a`, `--amend`, explicit path/index and hook mutation; bind only equivalent PASS and surface unsupported safely. |
| **SN03** | Pre-push new tip distinct from HEAD; two refs with one failure, deletion ref; independent non-deletion tip checks, full push blocks on failure. |
| **SN04** | Linked worktree, detached HEAD, submodule/gitlink, unavailable object, symlink/path escape and mutating hook; exact snapshot or explicit block, never false PASS. |
| **FP01** | Change tree, workflow YAML, event inputs, config, schema, `act` version, runner image; previously matching PASS becomes STALE. |
| **FP02** | Change secret/var inputs without plaintext in metadata; fingerprint invalidates reuse, confidential handling and file-permission expectations verified. |
| **ST01** | PASS, FAIL, INCOMPLETE, CANCELLED, NOT_RUN, STALE, unknown/corrupt state; only exact matching PASS can auto-allow. |
| **ST02** | Concurrent runs/worktrees, lock contention, crashed writer, truncated JSON, disk-full/denied write; atomic non-PASS recovery and correct run ownership. |
| **AC01** | Missing/old `act`, missing/unusable Docker/runtime; preflight INCOMPLETE, never PASS. |
| **AC02** | Fake `act`: zero/nonzero exit, failed job, skipped job, malformed output, crash, multiple workflows; correct overall/per-job result and diagnostics. |
| **AC03** | Cancel/timeout during process with a spawned child; bounded termination/cleanup, recorded INCOMPLETE, no child surviving. |
| **AC04** | Local secret/var files, missing required secret, restricted permissions, failed cleanup and logs; no plaintext in regular output, metadata or argv. |
| **AC05** | Captured fake `act --version`, workflow/job/skip output and staleness reason; version and safe diagnostic records available. |
| **HK01** | `init commit/push/both`, default both, repeated init, mode changes, disabled gate; idempotent owned hooks only. |
| **HK02** | Pre-commit index PASS/cache hit, changed staged content, no PASS invokes CI; fail-closed if validation non-PASS. |
| **HK03** | Push-only runs if no PASS, both-mode equivalence reuses PASS, multi-ref and deletion; one bad required ref blocks. |
| **HK04** | Failed, incomplete, stale, missing CLI, corrupt, crashed gate; stable ALLOW/BLOCK outcomes and non-zero Git status on block. |
| **HK05** | Existing user hooks, `core.hooksPath`, hook chain/mutation, uninstall/restore; preserve behavior or explicit block, no silent overwrite. |
| **AP01** | PASS no approval UI; configured gate with non-PASS requests provider with exact reason and state. |
| **AP02** | CLI force flag, stdin, env, caller-controlled file/stdout/exit and forged IPC; no machine-origin bypass. |
| **AP03** | Cancel and valid Continue paths; BLOCK_CANCELLED vs ALLOW_BYPASS, clear machine-readable/IDE result. |
| **AP04** | Approval exact repo/worktree/op/fingerprint/request, expiry, replay and concurrent consume, state change; reject replays/mismatch, never set PASS. |
| **AP05** | Close, crash, absent provider, headless/no session, timeout and unverified provider; block promptly with reason. |
| **AP06** | Verify audit written+synced **before** allow; disk-full, denied, partial, verification failure; BLOCK_ERROR, no operation completion claim or secrets. |
| **AP07** | Windows 11 native and WSL2 UI routing, automation/IPC adversarial attempts, provider trust capabilities; record resistance and residual risks (S6 gate). |
| **HOST01** | Fresh Windows/Linux/macOS binary install (native Windows separate from WSL2), no Go/Node/npm; Git/act/Docker prereqs and hook paths. |
| **HOST02** | Go CLI orchestration/command paths without Node/npm, dependency edge; command output and failure behavior. |
| **REM01** | Representative GitHub required checks with eligible/ineligible local proofs, absent API/permissions; D2-approved suppression only, safe fallback. |

## Fixture promotion checkpoints

S1: establish harness interfaces and Git temp-repo helpers **without running real act**. S2–S5: materialize and execute SN/EV/ST/AC/HK fixtures. S6 only assesses approval feasibility; S7 implements AP only after security approval. S8 validates REM after D2. S9–S10 execute real host, Docker, Git and `act` acceptance separately. Each future issue records the fixture ID, test command, host and observed result (not an assumed PASS).
