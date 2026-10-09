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

## 7. Request Binding

Each approval request is bound to:

- repository/worktree identity;
- gate operation;
- exact validation fingerprint;
- current non-PASS result;
- one approval request ID.

Before accepting Continue anyway, CI See verifies the provider response is genuine and unconsumed for this request, then rechecks that the request still matches the current operation/state.

If the state changed while approval was open, the request becomes invalid and the operation blocks or reevaluates.

## 8. Bypass Semantics

A successful developer approval produces:

```text
ALLOW_BYPASS
```

The bypass:

- applies only to the current Git operation;
- applies only to the exact bound state;
- is consumed once;
- must not become a reusable PASS;
- must not make future commits/pushes automatically allowed.

## 9. Cancellation and Failure

These outcomes block Git:

- developer chooses Cancel;
- approval window is closed;
- approval provider cannot start;
- provider crashes;
- request expires;
- request/state mismatch;
- CI See crashes.

They map to the appropriate blocking outcome such as:

```text
BLOCK_CANCELLED
BLOCK_ERROR
BLOCK_INCOMPLETE
```

Git receives a non-zero result.

## 10. Headless Behavior

CI See must not wait indefinitely when no human approval surface can be presented.

In a headless or unsupported environment:

```text
non-PASS
→ no approval provider available
→ block
→ explain how to rerun from a supported interactive environment
```

There is no automatic fallback bypass.

## 11. Local Audit Record

Each completed bypass decision is recorded locally with:

- timestamp;
- repository/worktree identity;
- operation;
- validation fingerprint;
- non-PASS status/reason;
- decision;
- approval-provider type;
- request ID.

Do not record secret values.

Audit records are diagnostic history, not reusable authorization.

## 12. User Experience Rules

- Keep the decision surface small and direct.
- Default focus must not make accidental bypass easy.
- Cancel is the safe/default outcome.
- Continue anyway must be visually explicit.
- Keyboard operation must be supported.
- Closing the UI must never mean approval.
- The terminal/IDE/agent receives a concise final allow/block message.

## 13. Security Boundary

**In-scope threat:** an AI agent can initiate Git/CLI commands, set process arguments/environment, write stdin, read output, and modify its checkout. None of these capabilities alone may authorize a CI See bypass.

**Residual risk:** a same-user agent with desktop automation, accessibility APIs, or control over the provider process/IPC may be able to click or forge approval. A normal Windows popup, including one launched from WSL2, does not by itself establish that a human approved it. Full host compromise and native Git `--no-verify` are outside CI See's enforcement guarantee.

Before shipping bypass, a security review must document the provider implementation, process/IPC permissions, UI automation exposure, caller-versus-approver trust separation, and limitations. Validate that the provider enforces the documented boundary. Where this cannot be demonstrated, disable Continue anyway and fail closed, rather than claiming agent-proof authorization.

The product must state its actual assurance level; never claim tamper-proof or guaranteed human-only enforcement on an uncontrolled local host.

## 14. Required Acceptance Tests

Implementation must prove:

1. PASS never opens approval UI;
2. non-PASS at a configured gate requests developer approval;
3. no CI See CLI flag can approve bypass;
4. stdin/piped input alone cannot approve bypass;
5. Cancel blocks Git;
6. Continue anyway returns ALLOW_BYPASS;
7. bypass applies only to the exact bound state;
8. bypass is one-time and does not become PASS;
9. changing state while approval is open invalidates the request;
10. closing the approval surface blocks Git;
11. unavailable/crashed approval provider blocks Git;
12. headless mode does not wait indefinitely or auto-approve;
13. bypass audit metadata contains no secret values;
14. terminal/IDE/agent receives an unambiguous final result;
15. Windows 11 + WSL2 can present the human-facing approval surface;
16. caller-controlled flags, environment, stdin, output, files, or process exit codes cannot forge approval;
17. provider response authentication, request binding, expiry, single-use consumption, and replay rejection work;
18. adversarial attempts to automate the provider UI or spoof its IPC are tested on Windows 11 + WSL2, with residual risks recorded;
19. a provider that fails trust-boundary verification has Continue anyway disabled and blocks non-PASS Git operations;
20. security review documents precisely which automated-agent capabilities are and are not resisted.

## 15. Deferred

This spec does not define:

- final visual styling of native approval windows;
- remote/team approval;
- policy servers;
- tamper-resistant enterprise enforcement;
- prevention of native Git `--no-verify`.

Next: MVP implementation plan.
