> [Back to Git gating specification index](../git-gating.md)

<!-- BEGIN ORIGINAL CONTRACT -->
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

<!-- END ORIGINAL CONTRACT -->
