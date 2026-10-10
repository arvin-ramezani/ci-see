# Git gating specification — Post-commit/push identity, fingerprints, staleness, gate outcomes

**Status:** Draft (inherited)  
**Canonical authority:** [Git gating specification](../git-gating.md)

<!-- BEGIN ORIGINAL CONTRACT -->
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

<!-- END ORIGINAL CONTRACT -->
