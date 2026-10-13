# CI See — Product Requirements Document

**Status:** Draft revision — owner-requested explicit CLI direction (2026-10-13)
**Authority:** Product requirements
**Canonical public path:** `docs/prd.md`

This is the compatibility index. **Requirements live in the linked topic documents.**
Original headings/IDs remain for legacy links; removed requirements are labeled.
Use only the sections relevant to the current task.

The revised MVP uses explicit CLI invocation, `act`, Docker/Linux containers,
and download-on-miss with verified cache reuse. Automatic Git gates and approval
UX are removed; hosted-run suppression is deferred. See the
[revision handoff](prd/documentation-and-next.md#revision-scope-and-authority)
for the requirement map and older documents awaiting reconciliation.
This PRD change does not change runtime code or complete D3 acceptance.

## 1. Product Summary

[Read the complete section](prd/product-and-flow.md#1-product-summary).

## 2. Problem

[Read the complete section](prd/product-and-flow.md#2-problem).

## 3. Target Users

[Read the complete section](prd/product-and-flow.md#3-target-users).

## 4. Product Principles

[Read the complete section](prd/product-and-flow.md#4-product-principles).

## 5. Recommended User Flow

[Read the complete section](prd/product-and-flow.md#5-recommended-user-flow).

### Commit gate

**Removed:** [Legacy section](prd/product-and-flow.md#commit-gate).

### AI-agent flow

[Read the complete section](prd/product-and-flow.md#ai-agent-flow).

### Push gate

**Removed:** [Legacy section](prd/product-and-flow.md#push-gate).

## 6. Functional Requirements

[Read the complete section](prd/repository-and-result-requirements.md#6-functional-requirements).

### FR-1 — Repository initialization

[Read the complete section](prd/repository-and-result-requirements.md#fr-1--repository-initialization).

### FR-2 — Workflow discovery

[Read the complete section](prd/repository-and-result-requirements.md#fr-2--workflow-discovery).

### FR-3 — One-command local CI

[Read the complete section](prd/repository-and-result-requirements.md#fr-3--one-command-local-ci).

### FR-4 — Event handling

[Read the complete section](prd/repository-and-result-requirements.md#fr-4--event-handling).

### FR-5 — Durable result identity

[Read the complete section](prd/repository-and-result-requirements.md#fr-5--durable-result-identity).

### FR-6 — Exact-state gating

**Revised to report accuracy:** [Read the section](prd/gating-and-remote-requirements.md#fr-6--exact-state-gating).

### FR-7 — Native Git workflow

[Read the complete section](prd/gating-and-remote-requirements.md#fr-7--native-git-workflow).

### FR-8 — Developer decision UI

**Removed:** [Read the disposition](prd/gating-and-remote-requirements.md#fr-8--developer-decision-ui).

### FR-9 — Developer-only bypass

**Removed:** [Read the disposition](prd/gating-and-remote-requirements.md#fr-9--developer-only-bypass).

### FR-10 — Remote duplicate-execution control

**Deferred:** [Read the boundary](prd/gating-and-remote-requirements.md#fr-10--remote-duplicate-execution-control).

### FR-11 — Logs and diagnostics

[Read the complete section](prd/diagnostics-and-platform-requirements.md#fr-11--logs-and-diagnostics).

### FR-12 — Secrets and variables

[Read the complete section](prd/diagnostics-and-platform-requirements.md#fr-12--secrets-and-variables).

### FR-13 — Cross-platform host support

[Read the complete section](prd/diagnostics-and-platform-requirements.md#fr-13--cross-platform-host-support).

## 7. Reliability and Safety Requirements

[Read the complete section](prd/quality-and-scope.md#7-reliability-and-safety-requirements).

## 8. Compatibility Goals

[Read the complete section](prd/quality-and-scope.md#8-compatibility-goals).

## 9. MVP Non-Goals

[Read the complete section](prd/quality-and-scope.md#9-mvp-non-goals).

## 10. MVP Scope

[Read the complete section](prd/quality-and-scope.md#10-mvp-scope).

## 11. Success Criteria

[Read the complete section](prd/success-and-decisions.md#11-success-criteria).

## 12. Confirmed Product Decisions

[Read the complete section](prd/success-and-decisions.md#12-confirmed-product-decisions).

## 13. Current MVP Direction Requiring Technical Validation

[Read the complete section](prd/success-and-decisions.md#13-current-mvp-direction-requiring-technical-validation).

## 14. Open Questions for Later Specs

[Read the complete section](prd/success-and-decisions.md#14-open-questions-for-later-specs).

## 15. Documentation Strategy

[Read the complete section](prd/documentation-and-next.md#15-documentation-strategy).

## 16. Next Spec-Driven Artifacts

[Read the complete section](prd/documentation-and-next.md#16-next-spec-driven-artifacts).

