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

Older architecture, specs, S0 mappings, and implementation slices are **not yet
reconciled**. Their gating/approval requirements and sequencing do not define the
revised product. Retain compatible reliability/security constraints until scoped
review; do not silently discard them.

Existing code remains evidence of the previous design. This PR does not change
runtime behavior or grant implementation acceptance, merge, or live-run authority.
D3 remains PARTIAL; removed gate tests need not become revised product goals
merely because historical evidence includes them.

## 16. Next Spec-Driven Artifacts

After independent review and owner acceptance of the revised PRD:

1. Revise [architecture](../architecture.md) for CLI/Go/act/Docker; remove daemon,
   interception, approval, and bypass components.
2. Revise [execution spec](../specs/local-ci-execution.md) for commands, events,
   current-file capture, outcomes/exits/JSON, dependencies, secrets, network
   boundaries, compatibility, cancellation, and cleanup.
3. Retire [Git gating](../specs/git-gating.md) and
   [approval UX](../specs/developer-approval-ux.md) as active product contracts,
   preserving history and repairing authority/navigation.
4. Replace [implementation sequencing](../implementation-plan.md) and reconcile
   [S0 traceability](../s0/index.md) with the smaller scope and platform acceptance.
5. Implement reviewed CLI slices. Add an optional agent skill once commands and
   reports are stable; a skill is not required to use the product.

Do not add a new cache service, Git protection, or hosted-suppression mechanism
during this PRD-only change.

### External capability references

Checked 2026-10-13; these inform design, not completed acceptance:

- [act usage](https://nektosact.com/usage/index.html): workflow/job/event selection and cache options.
- [act runners](https://nektosact.com/usage/runners.html): container images and hosted parity limits.
- [act unsupported features](https://nektosact.com/not_supported.html): known gaps.
- [GitHub Actions billing](https://docs.github.com/en/billing/concepts/product-billing/github-actions): allowances and public runner policy.
