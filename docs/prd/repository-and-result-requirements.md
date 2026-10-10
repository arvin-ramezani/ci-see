# Product requirements — FR-1–FR-5: initialization, workflow execution, and results

**Status:** Draft (inherited)  
**Canonical authority:** [Product requirements](../prd.md)

<!-- BEGIN ORIGINAL CONTRACT -->
## 6. Functional Requirements

### FR-1 — Repository initialization

Provide a simple initialization command:

~~~text
ci-see init
~~~

It should:

- detect the Git repository;
- detect GitHub Actions workflows;
- configure required Git integration;
- validate local execution dependencies;
- create only minimal local configuration;
- avoid unnecessary workflow-file edits.

### FR-2 — Workflow discovery

Automatically discover workflows under .github/workflows/.

Users should not register workflows manually.

### FR-3 — One-command local CI

Provide a simple full-validation command:

~~~text
ci-see
~~~

Current MVP direction:

- implement CI See orchestration and Git gating in **Go**;
- distribute prebuilt CLI binaries for supported host OS/architectures without a Node.js runtime dependency;
- use nektos/act as the GitHub Actions execution engine;
- use Docker or a compatible container runtime for consistent Linux execution.

These are implementation directions subject to architecture validation, not permanent product constraints.

### FR-4 — Event handling

Support at minimum:

- push;
- pull_request;
- workflow_dispatch where practical.

Normal usage should not require users to construct event JSON manually.

### FR-5 — Durable result identity

Every run must record enough identity to prove what was validated.

Minimum data:

- repository identity;
- staged-tree or Git-tree identity;
- commit identity when available;
- selected workflows and event;
- start/end timestamps;
- overall result;
- per-workflow/job result;
- execution engine and version;
- CI See configuration fingerprint.

Required states:

- PASS;
- FAIL;
- INCOMPLETE or CANCELLED;
- NOT_RUN;
- STALE.

<!-- END ORIGINAL CONTRACT -->
