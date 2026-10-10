# Git gating specification — Bypass boundary, concurrency, hooks, acceptance 1–20

**Status:** Draft (inherited)  
**Canonical authority:** [Git gating specification](../git-gating.md)

<!-- BEGIN ORIGINAL CONTRACT -->
## 12. Bypass Boundary

CI See itself must not expose a trivial machine-only `--force` bypass.

A CI See bypass must be:

- explicitly approved by the developer;
- bound narrowly to the operation and state;
- locally recorded;
- unusable as a reusable blanket approval.

The interaction mechanism belongs in the Developer Approval UX specification.

Native Git bypasses such as `--no-verify` are outside CI See's enforcement boundary and must be documented honestly; CI See must not claim hooks are impossible to bypass.

## 13. Concurrency and Crash Safety

Gating and state updates must be concurrency-safe.

Requirements:

- one run cannot overwrite another run's result;
- gate reads must not accept partially written results;
- PASS is written only after successful execution completes;
- interrupted validation remains non-PASS;
- locks must recover safely after crashes;
- repository/worktree state cannot satisfy another incompatible worktree.

If CI See crashes while Git is waiting, Git must fail rather than continue.

## 14. Hook Lifecycle

`ci-see init` and later gate-mode changes must be safe to run repeatedly.

CI See must:

- detect existing CI See hook integration;
- install only hooks required by the selected gate mode;
- remove CI See-owned hooks no longer required after a mode change;
- avoid duplicate installation;
- preserve compatible existing user/tool hooks;
- fail clearly rather than destructively overwrite an unsupported hook setup.

CI See should also provide a safe way to remove only its own hook integration.

Exact installation mechanics may vary by platform while preserving the same behavior.

## 15. Required Acceptance Tests

Implementation must prove:

1. `ci-see init` supports `commit`, `push`, and `both`;
2. `both` is the default;
3. gate mode can be changed later without reinstalling the repository;
4. disabled gates do not trigger CI See validation;
5. changing staged content invalidates an earlier PASS;
6. unstaged-only changes do not invalidate an exact staged-tree PASS;
7. commit gating automatically runs validation when no matching PASS exists;
8. push-only gating automatically runs validation when no matching PASS exists;
9. failed/incomplete validation blocks the configured gate;
10. post-commit binding occurs only when the created commit tree matches;
11. `both` mode reuses a safely bound PASS at push without unnecessary rerun;
12. a different commit/tree cannot reuse another PASS;
13. workflow/config/event changes make previous PASS stale;
14. relevant secret/variable changes make previous PASS stale without storing plaintext secrets;
15. multiple pushed refs are evaluated independently;
16. partial/corrupt state cannot become PASS;
17. concurrent gates/runs cannot cross-associate results;
18. CI See crashes return a blocking Git result;
19. existing hooks are not silently overwritten;
20. allow vs block is machine-readable and unambiguous.

## 16. Deferred

This spec does not define:

- the exact human approval surface;
- approval timeout/headless UX;
- remote GitHub-hosted duplicate-execution control;
- prevention of native Git `--no-verify`;
- native Windows/macOS hosted-runner emulation.

Next: Developer Approval UX specification.
<!-- END ORIGINAL CONTRACT -->
