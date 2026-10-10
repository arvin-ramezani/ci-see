> [Back to Developer approval specification index](../developer-approval-ux.md)

<!-- BEGIN ORIGINAL CONTRACT -->
# CI See — Developer Approval UX Specification

**Status:** Draft  
**Product:** CI See  
**Canonical path:** `docs/specs/developer-approval-ux.md`  
**Depends on:** `docs/prd.md`, `docs/architecture.md`, `docs/specs/git-gating.md`  
**Last updated:** 2026-10-09

## 1. Purpose

Define how a developer explicitly decides whether a configured Git operation may continue when local CI is not PASS.

Core rule:

> An AI agent or other non-human process must not be able to silently approve a CI See bypass.

## 2. Trigger

Approval is requested only when:

- a configured commit/push gate is active; and
- the required state is FAIL, INCOMPLETE, NOT_RUN, or STALE; and
- CI See cannot obtain PASS automatically.

PASS never needs approval.

## 3. Developer Choices

The approval surface presents exactly two primary choices:

- **Cancel** — block the Git operation.
- **Continue anyway** — bypass this gate for this exact operation and state.

Closing, timing out, crashing, or losing the approval surface is equivalent to Cancel.

## 4. Approval Surface

The CLI may display status in the terminal, IDE, or agent output.

However, **Continue anyway must not be approved through:**

- a CLI flag such as `--force`;
- environment variables;
- piped/stdin-only input;
- a writable approval file;
- a reusable token intended for automation.

For the MVP, CI See invokes a **short-lived human-facing approval provider** and waits synchronously for its result.

No long-running daemon is required.

The provider is behind an adapter so platform-specific implementations can differ without changing core gate rules.

**A popup is not proof of human approval.** Before enabling Continue anyway, the implementation must document and verify the provider's trust boundary:

- identify which agent-controlled inputs/processes can reach the provider and its response channel;
- accept only a request-specific, authenticated, single-use provider response bound to the operation and validation fingerprint; reject forged/replayed responses and never trust a child exit code, stdout text, or caller-written file alone;
- establish how the provider authenticates the confirmation source and what prevents the initiating AI/CLI process from supplying approval without interaction;
- demonstrate the claimed protection with adversarial tests on the supported platform.

If no provider can substantiate its claimed human-approval boundary, **Continue anyway remains unavailable** and non-PASS operations block. Do not weaken this rule to satisfy the CLI-only MVP constraint.

## 5. Primary Platform

Primary development and acceptance target:

```text
Windows 11 + WSL2
```

When CI See runs inside WSL2, the approval provider may invoke a host-visible Windows confirmation surface.

Linux and macOS use equivalent platform-specific providers when implemented.

If no supported human-facing provider is available, the operation blocks safely.

## 6. Approval Content

The surface must clearly show:

- operation: `git commit` or `git push`;
- repository;
- CI state: FAIL / INCOMPLETE / NOT_RUN / STALE;
- concise reason;
- exact-state short identifier;
- failed/affected workflow when known;
- **Cancel**;
- **Continue anyway**.

Example:

```text
CI See

git push is blocked
Status: FAIL
Repository: ci-see
Workflow: test
State: 8f31c2…

[Cancel] [Continue anyway]
```

Do not show secrets or unnecessary raw logs in the approval surface.

<!-- END ORIGINAL CONTRACT -->
