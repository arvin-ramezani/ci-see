# Developer approval specification — Binding, bypass, cancel/headless behavior, durable audit

**Status:** Draft (inherited)  
**Canonical authority:** [Developer approval specification](../developer-approval-ux.md)

<!-- BEGIN ORIGINAL CONTRACT -->
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
- must not make future commits/pushes automatically allowed;
- requires its audit record to be durably persisted before returning `ALLOW_BYPASS`.

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

**Mandatory persistence gate:** After verifying an explicit, valid approval, CI See must atomically and durably save its audit record **before** returning `ALLOW_BYPASS`. If writing, syncing, or verifying the record fails (including full disk or denied permission), CI See returns `BLOCK_ERROR` and the Git operation remains blocked. Do not treat the UI approval or a partial file write as sufficient.

A saved approval decision records authorization for the attempted operation, not proof that Git eventually completed. Repeated or concurrent requests must not reuse it. For Cancel/deny decisions, audit recording is best-effort and never permits continuation.

<!-- END ORIGINAL CONTRACT -->
