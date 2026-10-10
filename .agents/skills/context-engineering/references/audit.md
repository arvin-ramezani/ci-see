# Audit Repository Knowledge

Respect read-only scope. Inventory navigation before reading all content.
Use the existing index to choose canonical sources.
For a full audit, inspect all scoped files over the threshold and needed dependencies.
For a focused audit, report its limits.

## Inspect evidence

| Area | Question |
| --- | --- |
| Persistent context | Are instructions necessary and links task-specific? |
| Retrieval | Can representative tasks locate the required canonical context? |
| Responsibility | Does each file serve one coherent task context? |
| Authority | Does each concept have one canonical owner and identifiable status? |
| Duplication | Are repeated rules equivalent and intentionally linked? |
| Contracts | Are conditions, negative cases, exceptions, and acceptance criteria explicit? |
| Consistency | Do accepted sources contradict plans, ADRs, or implementation evidence? |
| Freshness | Is there concrete drift against code, tests, schemas, or official external docs? |
| Enforcement | Are claimed checks implemented and executed? |
| Navigation | Do local files, anchors, and incoming references resolve? |

Treat age and line count as inspection signals; neither alone proves a defect.
Distinguish an implementation defect from stale documentation and unresolved intent.
Inspect relevant code/tests for implementation comparisons and record what ran.
Do not infer freshness from a recent file modification date.

## Report findings

Give findings an ID, severity, source locator, evidence, impact, action, and owner.
Use `BLOCKING` for missing contracts or contradictions affecting the requested change;
use `REQUIRED` for concrete navigation/authority defects; use `ADVISORY` for optional
improvements to documentation that remains usable.
Avoid findings based only on formatting preference.
Record a size decision and reason for every reviewed large file, including cohesive exceptions.

## Act according to scope

For an improvement request, fix routine navigation, ownership, and justified
structure using the preservation workflow, then rerun affected checks.
For a read-only audit, return evidence and proposed destinations without edits.
Preserve contradictory product/security intent and identify the specific owner
decision. Continue work independent of it.
Do not promote proposals or historical text to accepted status during cleanup.

## Finish

Report scope/exclusions, findings/corrections, size decisions, verification,
semantic review method, and unresolved items.
Separate proposed checks from executed checks and code review from runtime tests.
