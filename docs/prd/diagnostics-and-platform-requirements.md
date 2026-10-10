# Product requirements — FR-11–FR-13: diagnostics, secrets, host support

**Status:** Draft (inherited)  
**Canonical authority:** [Product requirements](../prd.md)

<!-- BEGIN ORIGINAL CONTRACT -->
### FR-11 — Logs and diagnostics

Expose:

- workflows/jobs running;
- failures;
- skipped jobs;
- completion status;
- why a Git operation is blocked;
- why a previous result became stale.

Detailed engine logs must remain available.

### FR-12 — Secrets and variables

Support local secrets and variables without committing them to Git.

Separate:

- committed project configuration;
- local non-secret configuration;
- local secrets.

Secrets must not be printed in normal logs.

### FR-13 — Cross-platform host support

Target hosts:

- Windows 11;
- Linux;
- macOS.

The CLI and gating UX should be consistent.

MVP execution may be Linux-container-centric and is not required to fully emulate GitHub-hosted Windows or macOS runners.

<!-- END ORIGINAL CONTRACT -->
