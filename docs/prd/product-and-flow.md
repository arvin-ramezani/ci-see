> [Back to Product requirements index](../prd.md)

<!-- BEGIN ORIGINAL CONTRACT -->
# CI See — Product Requirements Document

**Status:** Draft  
**Product:** CI See  
**Canonical owner:** Product requirements  
**Canonical path:** docs/prd.md  
**Last updated:** 2026-10-09

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
9. **Lightweight installation.** CI See is delivered as a compiled Go CLI; users of prebuilt binaries should not need Go or Node.js installed.

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

<!-- END ORIGINAL CONTRACT -->
