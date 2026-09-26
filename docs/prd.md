# CI See — Product Requirements Document

**Status:** Draft  
**Product:** CI See  
**Canonical owner:** Product requirements  
**Canonical path:** docs/prd.md  
**Last updated:** 2026-09-26

## 1. Product Summary

CI See is a local-first developer tool that runs a repository's existing GitHub Actions workflows on the developer's own machine and integrates the result into the normal Git workflow.

Its purpose is to make local CI automatic, fast, and difficult to forget for both human developers and AI coding agents.

Core value:

- faster feedback than push-and-wait CI;
- reduced dependence on GitHub-hosted Actions minutes;
- reuse of existing GitHub Actions workflows;
- automatic CI enforcement before commit and push;
- explicit developer control when CI has failed or has not fully passed.

## 2. Problem

GitHub Actions is convenient, but frequent hosted execution creates latency and can consume limited hosted CI minutes.

nektos/act already enables local execution of GitHub Actions, but normal act usage still leaves several product problems:

1. Developers or AI agents must remember to run it.
2. Users may need to remember event flags, payloads, secrets, and runner options.
3. A commit or push can happen without a valid local CI result.
4. A previous PASS may no longer match the code being committed.
5. AI agents must not silently bypass failed CI.
6. Users need a simple decision experience when CI has not passed.
7. Local validation should not require a second CI definition.
8. A successful local run should not unnecessarily trigger duplicate expensive hosted CI after push.

## 3. Target Users

Primary users:

- solo developers using GitHub;
- GitHub Free users who want to reduce hosted Actions usage;
- developers using Codex, Claude, Cursor, or other AI coding agents;
- developers who already define CI in GitHub Actions.

Secondary users:

- small teams wanting a consistent local CI gate;
- open-source contributors validating workflows before PRs.

## 4. Product Principles

1. **GitHub Actions remains the CI definition.** CI See must not create a parallel pipeline that can drift.
2. **Local-first, not local-only.** Local execution improves feedback but is not assumed to perfectly reproduce GitHub-hosted infrastructure.
3. **No-memory UX.** Humans and agents should not need to remember to run CI.
4. **Developer owns bypass decisions.** AI agents must not silently approve continuation after CI failure.
5. **Exact-state validation.** A PASS applies only to the exact state that was validated.
6. **Fail safe.** Unknown, stale, interrupted, or failed state is never PASS.
7. **Cross-platform UX.** Windows, Linux, and macOS should present a consistent product experience.
8. **Low adoption cost.** Existing repositories should need minimal configuration and minimal workflow changes.

## 5. Recommended User Flow

Normal development:

~~~text
developer or AI agent changes code
→ local CI validation
→ git commit
→ git push
→ create or update PR
~~~

The user should not need to remember the local CI step manually.

### Commit gate

When git commit is attempted, CI See checks whether the exact staged state has a valid local result.

If PASS:

~~~text
git commit
→ CI See verifies exact matching PASS
→ commit continues
~~~

If FAIL, INCOMPLETE, NOT RUN, or STALE:

~~~text
git commit
→ CI See blocks or pauses the operation
→ developer gets a clear decision
~~~

Minimum developer choices:

- **Cancel** — do not continue.
- **Continue anyway** — explicitly bypass for the current operation and state.

### AI-agent flow

When an AI agent initiates a gated Git operation and CI is not valid:

1. the Git operation stops or waits;
2. the agent cannot silently choose the bypass;
3. CI See surfaces a simple approval UI to the developer;
4. the developer chooses Cancel or Continue anyway;
5. the Git process returns an unambiguous result;
6. the agent continues or stops based on that result.

The approval UX must not depend exclusively on terminal stdin because desktop and IDE agents may execute Git non-interactively.

### Push gate

Before git push, CI See verifies that the exact commit being pushed is covered by an acceptable local CI result or an explicit developer bypass.

Push protection is a second safety layer.

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

### FR-10 — Remote duplicate-execution control

CI See must provide a safe way to avoid unnecessary duplicate GitHub-hosted CI after the exact code already passed locally.

Constraints:

- little or no manual workflow editing;
- forgetting CI See configuration should not unexpectedly consume hosted CI where preventable;
- remote validation must not be weakened invisibly;
- the exact mechanism is an architecture decision, not fixed by this PRD.

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

## 7. Reliability and Safety Requirements

1. **Fail closed:** if CI See cannot prove PASS for the exact state, it is not PASS.
2. **No false PASS after crash:** interrupted or partial execution cannot become PASS.
3. **Concurrency safety:** simultaneous CI runs or Git operations cannot corrupt state or associate a result with the wrong code.
4. **Repository isolation:** state from one repository/worktree cannot satisfy another incompatible state.
5. **Workflow invalidation:** workflow/config changes invalidate dependent results.
6. **Visible bypass:** developer bypasses are recorded locally and visible.
7. **Safe recovery:** if CI See crashes while Git is waiting, Git fails safely instead of silently continuing.

## 8. Compatibility Goals

CI See should maximize compatibility with ordinary GitHub Actions by delegating execution to an established engine rather than inventing a new workflow language.

Initial focus:

- ubuntu-latest style jobs;
- common marketplace actions;
- service containers;
- environment variables and secrets;
- job dependencies;
- matrices where supported;
- common GitHub expressions and events.

Full GitHub-hosted runner parity is not an MVP promise.

## 9. MVP Non-Goals

The MVP will not:

- implement a GitHub Actions engine from scratch;
- guarantee exact GitHub-hosted runner-image parity;
- replace GitHub;
- provide hosted CI SaaS;
- fully emulate macOS or Windows hosted runners;
- let AI agents automatically approve failed CI;
- require a second CI configuration;
- become a general Docker orchestration tool.

## 10. MVP Scope

The first usable release must prove the complete primary loop:

1. install CI See;
2. initialize a repository;
3. discover GitHub Actions workflows;
4. run relevant workflows locally;
5. persist exact-state results;
6. automatically gate commit;
7. automatically gate push;
8. let the developer Cancel or Continue anyway after non-PASS state;
9. support Git initiated by AI-agent/IDE processes;
10. provide consistent supported behavior on Windows, Linux, and macOS;
11. provide enough diagnostics to understand failures.

## 11. Success Criteria

A developer can use an ordinary GitHub Actions repository and experience:

~~~text
change code
→ attempt commit
→ CI See verifies local CI state
→ PASS: commit proceeds
→ non-PASS: developer decides
→ push receives a second exact-state check
→ normal PR workflow continues
~~~

The MVP is successful when:

- users do not need to remember a separate CI command;
- AI agents cannot silently bypass a failed gate;
- GitHub Actions remains the CI source;
- PASS cannot be reused for changed code;
- the same core UX works across supported host operating systems.

## 12. Confirmed Product Decisions

- Product name: **CI See**.
- GitHub Actions workflows remain the primary CI definitions.
- Local execution is a core capability.
- Default development flow validates locally before commit/push.
- Commit and push must be automatically interceptable.
- Failed or incomplete CI requires developer control over bypass.
- AI agents must not silently make that bypass decision.
- Cross-platform host support is required.
- Normal Git/GitHub developer workflow should remain familiar.

## 13. Current MVP Direction Requiring Technical Validation

- nektos/act as the local workflow execution engine;
- Docker-compatible runtime as the execution substrate;
- Git hooks or equivalent interception for commit/push;
- CI See-owned native/desktop approval surface for agent-initiated Git operations.

## 14. Open Questions for Later Specs

1. What exact mechanism prevents unnecessary duplicate GitHub-hosted execution?
2. Which event should CI See choose when workflows expose multiple triggers?
3. What exact inputs form the validation fingerprint?
4. Should default execution run all eligible workflows or a policy-defined set?
5. How should the approval UI be implemented cross-platform?
6. What happens in headless sessions?
7. Which act runner-image strategy gives the best reliability/performance balance?
8. How should artifacts and caches behave locally?
9. How should local secrets be stored cross-platform?
10. How should bypass paths such as git --no-verify be handled?

## 15. Documentation Strategy

CI See will use a **centralized semantic documentation architecture**.

This fits a new single-product repository and gives humans and AI agents one predictable documentation entry point without creating unnecessary hierarchy.

Target evolution:

~~~text
docs/
├── prd.md
├── index.md
├── architecture.md
├── specs/
├── decisions/
└── operations/
~~~

Directories should be created only when they have real content.

Canonical ownership:

- docs/prd.md — product intent and requirements;
- architecture documents — system structure and technical boundaries;
- specs — implementable behavior and acceptance criteria;
- decisions — consequential technical choices and rationale;
- operations — install, release, troubleshooting, and runbooks;
- code/tests/hooks/CI — mechanical enforcement.

## 16. Next Spec-Driven Artifacts

After PRD acceptance:

1. **System architecture specification**
   - process model;
   - CLI/background-service boundaries;
   - Git interception;
   - state model;
   - execution-engine integration;
   - cross-platform strategy.

2. **Local CI execution specification**
   - workflow discovery;
   - events;
   - act integration;
   - runtime requirements;
   - logs, secrets, artifacts, cache;
   - parity boundaries.

3. **Git gating specification**
   - exact state identity;
   - commit/push gates;
   - stale-result rules;
   - bypass semantics;
   - concurrency and crash recovery.

4. **Developer approval UX specification**
   - human flow;
   - AI-agent flow;
   - blocking behavior;
   - desktop/headless handling;
   - accessibility and failure states.

5. **MVP implementation plan**
   - vertical slices;
   - acceptance tests;
   - release criteria.
