# CI See — Git Gating Specification

**Status:** Draft  
**Product:** CI See  
**Canonical path:** `docs/specs/git-gating.md`  
**Depends on:** `docs/prd.md`, `docs/architecture.md`, `docs/specs/local-ci-execution.md`  
**Last updated:** 2026-10-03

## 1. Purpose

Define how CI See protects `git commit` and/or `git push` using exact-state local CI results.

Core rule:

> A configured Git gate proceeds automatically only when CI See can prove that the exact relevant Git state has an acceptable local result.

Developer bypass behavior is defined separately in the approval UX specification.

## 2. Gate Mode

During `ci-see init`, the user chooses which Git operations are gated:

```text
commit
push
both
```

Default: **both**.

The setting is persisted in CI See configuration and can be changed later without reinitializing the repository.

Behavior:

- **commit** — run/check CI before commit; push is not gated by CI See.
- **push** — run/check CI before push; commit is not gated by CI See.
- **both** — commit is gated first; push reuses the matching commit-time PASS when still valid, so CI does not run twice unnecessarily.

Changing the gate mode does not make an unrelated previous result valid.

## 3. Git Integration

`ci-see init` installs only the hook integration required by the selected mode:

- **pre-commit** — validate the exact staged tree before commit;
- **post-commit** — when commit gating is enabled, bind a successful staged-tree result to the created commit when still equivalent;
- **pre-push** — validate each pushed ref's new tip before network transfer.

CI See must preserve existing repository hooks. Initialization must never silently overwrite unrelated hook logic.

Hook scripts contain no product rules; they invoke CI See gate commands.

Changing gate mode must update only CI See's own hook integration.

## 4. Commit Gate Target

When commit gating is enabled, CI See validates the Git **index**, not arbitrary working-tree contents.

CI See identifies the staged state using the tree object produced from the index, conceptually:

```text
git write-tree
```

Unstaged and untracked changes that are not part of the index must not affect whether that staged state has PASS.

A PASS used by the commit gate must have been produced against an execution snapshot representing that exact staged tree.

## 5. Automatic Validation

Users should not need to remember to run `ci-see` before a configured gate.

For commit mode:

```text
git commit
→ calculate staged identity
→ matching PASS exists: continue
→ otherwise run local CI for that staged snapshot
→ PASS: continue
→ non-PASS: enter developer decision boundary
```

For push mode:

```text
git push
→ identify pushed tip
→ matching PASS exists: continue
→ otherwise run local CI for that commit state
→ PASS: continue
→ non-PASS: enter developer decision boundary
```

A previous result for different content or execution inputs must not satisfy the gate.

## 6. Post-Commit Binding

When commit gating is enabled, a pre-commit run cannot know the final commit SHA.

After a successful commit, CI See may bind the validated result to the new commit when:

- the new commit tree exactly matches the validated staged tree;
- relevant workflow/config/execution inputs still match;
- the result is PASS.

In **both** mode, this allows pre-push to reuse the commit-time validation without running the same CI again.

If equivalence cannot be proven, no binding occurs.

## 7. Push Gate Target

When push gating is enabled, for each non-deletion ref update supplied to `pre-push`, CI See validates the **new local tip commit**.

A push may proceed automatically when each required new tip has:

- a directly matching PASS; or
- in **both** mode, a PASS safely bound from an equivalent staged-tree validation.

If no reusable PASS exists, push gating runs local CI once for the required commit state.

Deleting a remote ref does not require local CI validation.

If multiple refs are pushed, failure of any required ref blocks the push unless the developer approval policy permits continuation.

## 8. Validation Fingerprint

A reusable PASS is identified by a deterministic fingerprint covering at minimum:

- repository/worktree identity;
- Git tree identity;
- selected workflows;
- workflow file contents;
- selected events and meaningful generated event inputs;
- CI See configuration;
- CI See execution-policy/schema version;
- detected `act` version;
- runner/image configuration that affects execution;
- local variables supplied to execution;
- local secrets supplied to execution, represented without storing plaintext.

Secrets must influence fingerprint validity without exposing secret values in metadata. The implementation should use a private keyed or salted digest rather than persist raw secret values.

Commit SHA is recorded when available, but code equivalence is anchored by the validated Git tree plus the full execution fingerprint.

## 9. Result States

For gating purposes:

- **PASS** — exact matching validation succeeded.
- **FAIL** — matching validation ran and failed.
- **INCOMPLETE** — matching validation could not complete reliably.
- **NOT_RUN** — no result exists for the required state.
- **STALE** — a result exists but no longer matches the required fingerprint.

Only PASS allows automatic continuation at a configured gate.

## 10. Staleness

A previous PASS becomes STALE when any fingerprint input changes.

Examples:

- staged code changes;
- workflow YAML changes;
- CI See configuration changes;
- selected event/input changes;
- relevant variable/secret changes;
- `act` version changes;
- runner/image configuration changes.

CI See should report the known reason for staleness when practical.

## 11. Gate Outcomes

CI See exposes stable semantic outcomes:

- **ALLOW_PASS** — exact PASS; Git continues.
- **ALLOW_BYPASS** — explicit valid developer bypass; Git continues.
- **BLOCK_FAIL** — validation failed.
- **BLOCK_INCOMPLETE** — validation could not complete safely.
- **BLOCK_CANCELLED** — developer declined continuation.
- **BLOCK_ERROR** — CI See could not evaluate the gate safely.

Process behavior:

- allow outcomes return success to Git;
- block outcomes return non-zero to Git;
- crashes or unknown states return non-zero.

This gives terminals, IDEs, and AI agents an unambiguous continue/stop result.

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
