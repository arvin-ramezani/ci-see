# Product requirements — Documentation authority and revision handoff

**Status:** Draft revision — owner-requested CLI direction
**Canonical authority:** [Product requirements](../prd.md)

## 15. Documentation Strategy

Use [docs/index.md](../index.md) as the documentation entrypoint.
PRD owns intent; architecture owns system boundaries; specs own behavior; plans
own sequencing; evidence owns recorded results. Link rather than duplicate.
Create directories only when they have content.

### Revision scope and authority

This PR updates the PRD first, as requested. It deliberately changes meaning
from automatic Git enforcement to an explicit CLI. IDs, public paths, and legacy
headings remain for traceability. Earlier wording is recoverable at baseline
`d979d9b2a631c121aaee0ef60185201159ce4d11`.

| Baseline requirement | Disposition |
| --- | --- |
| FR-1: init/Git integration | Retain setup; remove Git integration |
| FR-2: discovery | Retain; add listing/selection and empty-selection handling |
| FR-3: local CI | Retain Go/act; explicit invocation, Docker, download/reuse policy |
| FR-4: events | Retain; explicit inputs and local-context limits |
| FR-5: identity | Retain reporting; actual local inputs and JSON replace staged-only gating state |
| FR-6: exact-state gates | Remove gates/PASS reuse; retain tested-input accuracy |
| FR-7: interception | Replace with independent Git and explicit verification |
| FR-8: approval UI | Remove; terminal/JSON and non-interactive CLI |
| FR-9: bypass/audit | Remove; truthful reporting retained |
| FR-10: hosted duplicates | Defer beyond MVP; no automatic GitHub changes |
| FR-11: diagnostics | Retain; remove Git-blocking reasons, add setup/cache/scope |
| FR-12: secrets | Retain separation/redaction; explicit credential inputs |
| FR-13: hosts | Retain three hosts; distinguish them from Linux job execution |
| Sections 1–5, 7–14 | Rewrite flow, principles, safety, scope, success, decisions and open questions for explicit CLI |
| Sections 15–16 | Retain ownership; replace gating/approval delivery with revision handoff |

The architecture, local execution contract, implementation slices and current S0
mapping are now aligned in the [contract alignment record](../migration/cli-contract-alignment.md).
Git-gating/approval specifications and old S0/D3 procedures are historical. Open
CLI/source/dependency/security decisions are explicit; alignment does not mean
all implementation contracts are finalized or runtime features exist.

Existing source remains unchanged. D3 stays PARTIAL; retired gate tests are not
new product goals, and unverified runtime evidence does not become PASS.

## 16. Next Spec-Driven Artifacts

1. Review and accept the aligned [architecture](../architecture.md),
   [execution spec](../specs/local-ci-execution.md) and [plan](../implementation-plan.md).
2. Resolve the relevant [open implementation contracts](../implementation-plan/authority-and-gates.md#current-cli-decisions)
   before their owning feature slice; do not invent security/source policy.
3. Open a separate scoped cleanup/command-contract PR: keep the Go foundation,
   retire gate-only models and decide disposition of separate D3 tooling.
4. Build fresh discovery, dependency, execution and report features under the
   revised plan; perform separate authorized real-tool/host acceptance before release.
5. Add an optional agent skill once commands/reports are stable.

This document and alignment PR do not authorize code deletion, closing PR #19,
VM/firewall changes, live act or production acceptance.

### External capability references

Checked 2026-10-13; these inform design, not completed acceptance:

- [act usage](https://nektosact.com/usage/index.html): workflow/job/event selection and cache options.
- [act runners](https://nektosact.com/usage/runners.html): container images and hosted parity limits.
- [act unsupported features](https://nektosact.com/not_supported.html): known gaps.
- [GitHub Actions billing](https://docs.github.com/en/billing/concepts/product-billing/github-actions): allowances and public runner policy.
