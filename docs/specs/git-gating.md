# CI See — Git Gating Specification

**Status:** Draft  
**Product:** CI See  
**Canonical path:** `docs/specs/git-gating.md`  
**Depends on:** `docs/prd.md`, `docs/architecture.md`, `docs/specs/local-ci-execution.md`  
**Last updated:** 2026-10-03

## 1. Purpose

Define how CI See automatically protects `git commit` and `git push` using exact-state local CI results.

Core rule:

> Git proceeds automatically only when CI See can prove that the exact relevant Git state has an acceptable local result.

Developer bypass behavior is defined separately in the approval UX specification.

## 2. Git Integration

`ci-see init` installs thin Git integration for:

- **pre-commit** — validate the exact staged tree before commit;
- **post-commit** — bind a successful staged-tree result to the created commit when still equivalent;
- **pre-push** — validate each pushed ref's new tip before network transfer.

CI See must preserve existing repository hooks. Initialization must never silently overwrite unrelated hook logic.

Hook scripts contain no product rules; they invoke CI See gate commands.

## 3. Commit Gate Target

The commit gate validates the Git **index**, not arbitrary working-tree contents.

CI See identifies the staged state using the tree object produced from the index, conceptually:

```text
git write-tree
```

Unstaged and untracked changes that are not part of the index must not affect whether that staged state has PASS.

A PASS used by the commit gate must have been produced against an execution snapshot representing that exact staged tree.

## 4. Automatic Validation

Users should not need to remember to run `ci-see` before Git.

On `git commit`:

```text
calculate staged identity
→ matching PASS exists: continue
→ otherwise run local CI for that staged snapshot
→ PASS: continue
→ non-PASS: enter developer decision boundary
```

A previous result for different content or execution inputs must not satisfy the gate.

## 5. Post-Commit Binding

A pre-commit run cannot know the final commit SHA.

After a successful commit, CI See may bind the validated result to the new commit when:

- the new commit tree exactly matches the validated staged tree;
- relevant workflow/config/execution inputs still match;
- the result is PASS.

This allows pre-push to reuse the commit-time validation without running the same CI again.

If equivalence cannot be proven, no binding occurs.

## 6. Push Gate Target

For each non-deletion ref update supplied to `pre-push`, CI See validates the **new local tip commit**.

A push may proceed automatically when each required new tip has:

- a directly matching PASS; or
- a PASS safely bound from an equivalent staged-tree validation.

Deleting a remote ref does not require local CI validation.

If multiple refs are pushed, failure of any required ref blocks the push unless the developer approval policy permits continuation.

## 7. Validation Fingerprint

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

## 8. Result States

For gating purposes:

- **PASS** — exact matching validation succeeded.
- **FAIL** — matching validation ran and failed.
- **INCOMPLETE** — matching validation could not complete reliably.
- **NOT_RUN** — no result exists for the required state.
- **STALE** — a result exists but no longer matches the required fingerprint.

Only PASS allows automatic continuation.

## 9. Staleness

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

## 10. Gate Outcomes

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

## 11. Bypass Boundary

CI See itself must not expose a trivial machine-only `--force` bypass.

A CI See bypass must be:

- explicitly approved by the developer;
- bound narrowly to the operation and state;
- locally recorded;
- unusable as a reusable blanket approval.

The interaction mechanism belongs in the Developer Approval UX specification.

Native Git bypasses such as `--no-verify` are outside CI See's enforcement boundary and must be documented honestly; CI See must not claim hooks are impossible to bypass.

## 12. Concurrency and Crash Safety

Gating and state updates must be concurrency-safe.

Requirements:

- one run cannot overwrite another run's result;
- gate reads must not accept partially written results;
- PASS is written only after successful execution completes;
- interrupted validation remains non-PASS;
- locks must recover safely after crashes;
- repository/worktree state cannot satisfy another incompatible worktree.

If CI See crashes while Git is waiting, Git must fail rather than continue.

## 13. Hook Lifecycle

`ci-see init` must be safe to run repeatedly.

It must:

- detect existing CI See hook integration;
- avoid duplicate installation;
- preserve compatible existing user/tool hooks;
- fail clearly rather than destructively overwrite an unsupported hook setup.

CI See should also provide a safe way to remove only its own hook integration.

Exact installation mechanics may vary by platform while preserving the same behavior.

## 14. Required Acceptance Tests

Implementation must prove:

1. changing staged content invalidates an earlier PASS;
2. unstaged-only changes do not invalidate an exact staged-tree PASS;
3. commit automatically runs validation when no matching PASS exists;
4. failed/incomplete validation blocks automatic commit;
5. post-commit binding occurs only when the created commit tree matches;
6. push reuses a safely bound PASS without unnecessary rerun;
7. a different commit/tree cannot reuse another PASS;
8. workflow/config/event changes make previous PASS stale;
9. relevant secret/variable changes make previous PASS stale without storing plaintext secrets;
10. multiple pushed refs are evaluated independently;
11. partial/corrupt state cannot become PASS;
12. concurrent gates/runs cannot cross-associate results;
13. CI See crashes return a blocking Git result;
14. existing hooks are not silently overwritten;
15. allow vs block is machine-readable and unambiguous.

## 15. Deferred

This spec does not define:

- the exact human approval surface;
- approval timeout/headless UX;
- remote GitHub-hosted duplicate-execution control;
- prevention of native Git `--no-verify`;
- native Windows/macOS hosted-runner emulation.

Next: Developer Approval UX specification.
