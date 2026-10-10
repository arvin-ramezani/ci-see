# Issue #10 — Documentation Structure and Preservation Review

**Scope:** Docs-only, original baseline commit `770238b55a1c532c541eea1e5dfe4efe66c594ab` (2026-10-10).
**Skill:** [context-engineering](../../.agents/skills/context-engineering/SKILL.md) installed from
`arvin-ramezani/ai-skills` at `6ec528e498b51a4e1a9796bb55469bc188e5fdc1`.
**Decision:** No behavioral, security, or owner-decision change authorized.

## Reviewed baseline and decisions

| Canonical path | Before (physical lines) | Decision | Rationale / disposition |
| --- | ---: | --- | --- |
| [`docs/prd.md`](../prd.md) | 467 | SPLIT | 7 task-related children; preserve source verbatim; each legacy heading routes to its original section. |
| [`docs/architecture.md`](../architecture.md) | 225 | SPLIT | 3 task-related children; preserve source verbatim; each legacy heading routes to its original section. |
| [`docs/implementation-plan.md`](../implementation-plan.md) | 87 | SPLIT | 3 task-related children; preserve source verbatim; each legacy heading routes to its original section. |
| [`docs/specs/local-ci-execution.md`](../specs/local-ci-execution.md) | 244 | SPLIT | 3 task-related children; preserve source verbatim; each legacy heading routes to its original section. |
| [`docs/specs/git-gating.md`](../specs/git-gating.md) | 275 | SPLIT | 3 task-related children; preserve source verbatim; each legacy heading routes to its original section. |
| [`docs/specs/developer-approval-ux.md`](../specs/developer-approval-ux.md) | 253 | SPLIT | 3 task-related children; preserve source verbatim; each legacy heading routes to its original section. |

`docs/implementation-plan.md` is short (87 physical lines), but contains dense
slices, D1–D4 gates and evidence tables. Splitting is justified by task retrieval,
not line count. `README.md`, `AGENTS.md`, `docs/index.md`, and
`docs/context-policy.md` remain compact routers/policy; no further split.
A child may be >100 lines where the cohesive topic warrants it; there is no hard cap.

## Preservation: moved versus edited

- **Moved verbatim:** Every character of each of the six baseline source documents,
  in source order, to disjoint, contiguous child segments. The transformation
  reconstructed every original byte-for-byte string from its ordered child segments
  before uploading. Original source Git blobs and exact baseline line ranges are
  recorded in [manifest](manifest.json). No original clause was summarized away.
- **Edited:** The six canonical files became navigational compatibility indexes;
  README, root AGENTS and docs index were updated with contextual routes. Added
  this evidence, link/contract checker and GitHub Actions docs-only check.
- **Legacy anchors:** Every original H2/H3 (and deeper) heading is retained under
  the same text and Markdown level at the same canonical path. Links point to
  the original heading preserved in its new file; [manifest](manifest.json)
  records heading, legacy slug, original source line, and destination.
- **Single owners:** Child files own normative wording. Compatibility headings,
  navigation and this evidence are not independent contract definitions.

## Contract preservation coverage

| Contract group | Unique authoritative destination | Automated check |
| --- | --- | --- |
| FR-1–FR-13 | `docs/prd/*-requirements.md` | Exactly 13 FR heading IDs |
| Local execution acceptance 1–12 | `docs/specs/local-ci-execution/results-and-acceptance.md` | Consecutive 1–12 |
| Git gating acceptance 1–20 | `docs/specs/git-gating/bypass-lifecycle-and-acceptance.md` | Consecutive 1–20 |
| Approval acceptance 1–23 | `docs/specs/developer-approval-ux/security-and-acceptance.md` | Consecutive 1–23 |
| S0–S10 and D1–D4 | `docs/implementation-plan/vertical-slices.md`, `authority-and-gates.md` | IDs once in defining table |

Semantic self-check: All normative body text remained verbatim, including
staged-index versus pushed-tip identity, no unstaged/untracked leakage, audit
persisted **before** bypass, anti-spoofing, fail-closed outcomes, FR-10,
Windows-native versus WSL2, and unresolved D1–D4. This was not an independent
contract/security review. No Go application tests were performed or claimed.

## Verification and maintenance

Run `python3 scripts/validate_docs.py --check`. The checker validates Markdown
path/fragment targets, entrypoint reachability, legacy anchors and their routed
links, 13 FR heading owners, numbered 12+20+23 acceptance clauses, and S/D table
IDs. It reports line counts and over-150 review signals, not line-limit failures.
It cannot prove future semantic equivalence or conformance of application code.

**Issue #9 coordination:** Keep S0's requirements trace map separate, using
these canonical paths or direct topic links. D1–D4 are still unresolved gates;
no issue #9 traceability edits, Go code, hooks or releases are included here.

## Primary OpenAI guidance (verified 2026-10-10)

- [Rethinking skills and prompts for GPT-6 Astra (2026-09-11)](https://developers.openai.com/blog/rethinking-skills-and-prompts-for-gpt-6-astra): scope AGENTS.md guidance and use progressive disclosure.
- [Custom Code Review rules for Codex (2026-07-20)](https://developers.openai.com/blog/custom-code-review-rules-for-codex): narrowly scope non-obvious invariants.

80–100/150 physical line review conventions belong to this repository's skill,
not an OpenAI document-length requirement. Independent review and explicit owner
approval are required before merging this PR.
