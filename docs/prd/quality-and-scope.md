# Product requirements — Reliability, compatibility, lightweight MVP

**Status:** Draft revision — owner-requested CLI direction
**Canonical authority:** [Product requirements](../prd.md)

## 7. Reliability and Safety Requirements

1. **No false PASS:** unknown required behavior, unsupported selected jobs, missing
   tools, acquisition/integrity failures, interruptions, and crashes cannot produce PASS.
2. **Accurate scope:** show expected exclusions; failing to execute required selected
   validation is non-PASS even when other jobs succeed.
3. **Run isolation:** concurrent runs and separate repositories/worktrees cannot corrupt
   reports, overwrite each other's artifacts, or misassociate results.
4. **Stable inputs:** capture/identify tested inputs; detected changes during a run
   cannot leave ambiguous PASS attributed to the current repository.
5. **Scoped cleanup:** failure/cancellation preserve diagnostics and clean only
   run-owned resources; unrelated Docker resources and host settings are untouched.
6. **Credential control:** explicit secrets, scoped acquisition access, and redaction;
   arbitrary workflows must not be assumed safe to execute locally.
7. **Lightweight operation:** no daemon, approval GUI, cache service, or always-running
   process; avoid redundant downloads of compatible verified cached content.

## 8. Compatibility Goals

Use existing GitHub Actions YAML as written for compatible workloads. No second
CI definition or mandatory YAML rewrite. Local configuration may select events,
jobs, inputs, secrets, and supported runner-image mappings.

Initial targets: Linux runner jobs, common actions, service containers, dependencies,
matrices, expressions, and events where the tested `act` version supports them.
Publish a tested compatibility matrix rather than promise every feature works.
Unsupported required behavior must be visible and non-PASS.

See [act runners](https://nektosact.com/usage/runners.html) and
[unsupported features](https://nektosact.com/not_supported.html).
Full hosted-runner parity is not promised.

## 9. MVP Non-Goals

- Automatic CI, Git hooks, commit/push enforcement, or bypass approval/audit.
- Daemon, desktop application, hosted CI service, or new workflow interpreter.
- Second pipeline language or general Docker orchestration.
- Exact hosted parity or native Windows/macOS runner emulation.
- Custom caching engine or guaranteed offline execution of every workflow.
- Automatic GitHub suppression, remote edits, or required-check publication.
- Production sandboxing of arbitrary untrusted repositories/actions.
- Agent-specific runtime; an optional small skill can invoke the ordinary CLI later.

## 10. MVP Scope

The first usable release must prove:

1. Install prebuilt CI See and required `act`/Docker dependencies.
2. Discover workflows and explain missing setup without Git integration.
3. Explicitly invoke default or selected workflow/job/event validation.
4. Download missing dependencies, verify identities, and reuse compatible caches.
5. Execute supported existing workflows with progress and actionable diagnostics.
6. Return accurate outcomes/exits and optional JSON/local reports.
7. Demonstrate cold/warm-cache runs on supported Windows, macOS, and Linux hosts.

Hosted suppression is outside this loop. Existing code and older D3 evidence
do not establish implementation or production acceptance of the revised MVP.
