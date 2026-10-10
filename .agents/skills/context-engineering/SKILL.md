---
name: context-engineering
description: Govern repository documentation structure, agent routing, and preservation. Use for documentation audits, canonical-source conflicts, large-file review and splitting, or maintenance of indexes and context policy. This skill owns structural migration and splitting approval rules; use doc-strategy-engineer for feature-document authoring and agent-memory strategy.
---

# Context Engineering

Build task-relevant, indexed repository knowledge that survives agent sessions.
Perform the review and justified restructuring; produce evidence another agent can inspect.

## Skill ownership

This skill is the source of truth for documentation review thresholds, routing,
structural reorganization and splitting permissions, and preservation verification.
Use `doc-strategy-engineer` for feature/FA authoring, documentation strategy proposals,
and cross-session memory when installed. In mixed tasks, apply this skill's
structural rules while the other skill owns the content-specific workflow.
Do not use file count as an approval trigger. Within an authorized documentation
improvement, perform justified splits/moves regardless of count, with preservation
evidence. Escalate only out-of-scope changes, unrecoverable destructive changes,
or unresolved product, security, architecture, or accepted-contract decisions.

## Select the workflow

| Requested outcome | Load next |
| --- | --- |
| Establish documentation and operating rules | [Bootstrap](references/bootstrap.md) |
| Create, update, move, or split documentation | [Restructure](references/restructure.md) |
| Assess consistency, discoverability, or drift | [Audit](references/audit.md) |
| Explain the policy or verify its provenance | [Sources](references/sources.md) |

Read [policy](references/policy.md) for applicable conventions.
Load only the selected workflow and relevant project sources; avoid preloading all references.
Use [report template](assets/report-template.md) for an auditable handoff.

## Establish scope and authority

1. Identify the requested repository, applicable agent instructions, existing indexes,
   documentation tooling, and working-tree changes.
2. Locate canonical sources through task-specific routing. Inventory files before
   deciding which contents to read. Follow dependencies needed to interpret the task.
3. Distinguish accepted contracts, draft proposals, historical rationale, execution
   plans, implementation evidence, and external references.
4. Preserve existing owner decisions and project conventions. Follow the user's
   instruction when it conflicts with this skill and report the relevant convention.

## Review document size

- Prefer roughly 80–100 lines for new focused prose when clear and complete.
- Accept 100–150 lines when the material remains cohesive.
- Review every document over 150 physical lines within the requested scope.
- Review shorter documents too when authority, duplication, or retrieval problems exist.
- Record `KEEP`, `SPLIT`, or `REORGANIZE` with a concrete reason for each reviewed file.
- For 300–400+ lines, actively seek independent retrieval units; document a cohesive
  exception when retaining the file. Never shorten by deleting necessary context.

Treat these thresholds as local review conventions. A threshold triggers judgment;
it does not prove a defect or require a particular resulting line count.
When authorized to improve documentation, execute justified splits in the same task.
For an explicit read-only audit, provide findings and a destination map without edits.

## Preserve meaning while changing structure

Capture a retrievable original revision or content snapshot before editing.
Map every affected requirement, invariant, decision, definition, exception,
failure case, acceptance criterion, and necessary explanatory example to its owner.
Use existing stable IDs; otherwise use original path, heading, and baseline line range.
Compare original and destination text in a separate pass after restructuring.
Preserve obligation strength, conditions, scope, ordering, and negative behavior.
Repair incoming links and indexes, including references in code and automation.

Do not silently resolve product, security, or architectural contradictions.
Retain conflicting evidence, flag a genuine owner decision, and complete unaffected work.
Proceed with routine structural choices inside the user's authorized task.
Do not treat a split request as permission to change application behavior or merge a PR.

## Verify and hand off

Use existing repository checks when available. For a portable structural inventory, run:

```bash
python <skill-dir>/scripts/inspect_docs.py <repo-root> --docs docs --entry docs/index.md
```

Adapt documentation roots and exclusions to the repository. Read the JSON limitations.
Use `--check` only after choosing the intended index coverage.
The scanner reports line counts, hashes, headings, file links, and index reachability;
it cannot prove semantic equivalence, valid anchors, freshness, or architectural compliance.

Report scope and baseline, review decisions, changed paths, content-preservation map,
structural checks, semantic review method, remaining conflicts, and next action.
Use `PASS`, `PARTIAL`, or `BLOCKED` with evidence; never infer an independent review
from a self-review or an unrun check. A cohesive file may remain over 150 lines.
For a small edit, keep the report proportional; retain affected preservation evidence.
Finish when requested restructuring, navigation repairs, and relevant checks are complete.
