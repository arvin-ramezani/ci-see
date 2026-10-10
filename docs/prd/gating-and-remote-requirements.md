# Product requirements — FR-6–FR-10: Git gates, approvals, hosted CI

**Status:** Draft (inherited)  
**Canonical authority:** [Product requirements](../prd.md)

<!-- BEGIN ORIGINAL CONTRACT -->
### FR-6 — Exact-state gating

A PASS must never be reused for different code.

Changes to relevant code, workflow files, CI See configuration, event inputs, or other validation inputs must invalidate the result.

### FR-7 — Native Git workflow

Normal git commit and git push commands must be gateable automatically.

Users should not be forced to replace Git commands with special CI See wrappers.

### FR-8 — Developer decision UI

For non-PASS state, CI See must present a clear product-owned decision surface.

Example:

~~~text
Local CI did not pass.

Operation: git commit
Status: FAILED

[Cancel] [Continue anyway]
~~~

It must work when Git is initiated by:

- a human terminal;
- an IDE;
- a desktop AI agent;
- an IDE AI agent.

### FR-9 — Developer-only bypass

A bypass must be an explicit developer decision.

CI See must not expose a trivial machine-only approval path that lets an AI agent silently self-approve.

A bypass should be bound as narrowly as practical to the exact operation and repository state.

A successful bypass must have its audit record durably saved **before** Git is allowed to continue. If the record cannot be saved or verified, CI See must block the operation; the approval alone is not sufficient.

### FR-10 — Remote duplicate-execution control

CI See must provide a safe way to avoid unnecessary duplicate GitHub-hosted CI after the exact code already passed locally.

Constraints:

- little or no manual workflow editing;
- forgetting CI See configuration should not unexpectedly consume hosted CI where preventable;
- remote validation must not be weakened invisibly;
- the exact mechanism is an architecture decision, not fixed by this PRD.

<!-- END ORIGINAL CONTRACT -->
