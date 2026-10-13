# Product requirements — Revised Git boundaries and deferred hosted control

**Status:** Draft revision — owner-requested CLI direction
**Canonical authority:** [Product requirements](../prd.md)

IDs and headings remain for traceability. Removed requirements must not be
implemented merely because an older specification still references them.

### FR-6 — Exact-state gating

**Automatic gating removed; report accuracy retained.** Reports describe only
the inputs and selected scope actually tested under FR-5. Historical PASS cannot
claim validation of later code/configuration/dependency changes.

There is no staged-tree gate, commit-to-push PASS reuse, or authorization of Git
operations from reports. Each CLI invocation requests a new run;
dependency-cache reuse does not mean reusing a previous PASS.

### FR-7 — Native Git workflow

**Revised:** Git commands work independently of CI See. Developers and agents
explicitly invoke the CLI when they want verification.

No hooks, automatic runs, commit/push wrappers, interception, or Git-blocking
service are part of the product. Running CI See must not change Git behavior.

### FR-8 — Developer decision UI

**Removed.** No desktop approval UI, waiting Git process, or Cancel/Continue-anyway
surface. Terminal output and optional JSON provide results.
Automation must run non-interactively using explicit inputs; missing required
inputs terminate clearly rather than waiting for a GUI.

### FR-9 — Developer-only bypass

**Removed.** There is no CI See gate to bypass and no bypass token, human approval
channel, or durable bypass audit. Developers decide whether/when to commit or push.

The CLI and agent integration must still report actual results: FAIL, INCOMPLETE,
and NOT_RUN cannot be represented as PASS.

### FR-10 — Remote duplicate-execution control

**Deferred beyond MVP.** A later feature may help avoid duplicate hosted execution
after local validation, subject to feasibility and a separately accepted design.

For MVP, do not disable workflows, edit triggers, alter repository settings, add
skip markers, or publish GitHub statuses automatically. A local run neither
suppresses hosted Actions nor satisfies GitHub required checks by itself.

Future suppression needs explicit user control, visible effects, restoration,
and treatment of branch protection and collaborator runs. Its mechanism and
safety acceptance are not decided here.
