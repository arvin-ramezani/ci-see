# Repository Context Policy

## Context layers

Keep persistent instructions small. Use `AGENTS.md` as an operating map with
task-specific links to deeper material. Keep durable knowledge in versioned
repository documents; load only what the current task needs.

Adapt categories to existing content rather than creating empty directory trees:

| Knowledge | Canonical responsibility |
| --- | --- |
| Product | Outcomes, requirements, scope, constraints |
| Architecture | Boundaries, dependencies, runtime model, invariants |
| Specifications | Testable behavior, conditions, failure cases |
| Decisions | Accepted rationale, alternatives, superseded choices |
| Plans | Bounded execution sequence and progress |
| Operations | Installation, deployment, recovery, troubleshooting |
| References | External or generated technical material |

Give each concept one canonical owner. Link to it from summaries.
Keep current contracts in specifications/architecture and decision history in ADRs.
Use code, schemas, and tests as implementation evidence. Report disagreement between
accepted contracts and implementation rather than choosing either silently.
Never let a plan, index, or agent guide redefine an accepted product contract.

## Document shape

Give each document a clear responsibility and recognizable heading.
Use descriptive filenames, concise sentences, explicit conditions, and existing stable IDs.
Include minimal status/authority metadata and relevant dependencies in canonical docs.
Distinguish wording review from review against implementation; update the latter date
only when relevant implementation was actually checked.
Keep the conditions of a rule beside it or link to its precise dependency.
Preserve useful tables, diagrams, and interpretive examples. Avoid prose that merely
repeats obvious source code.

## Routing

Create an index when multiple documents are independently useful.
Use actual Markdown links with ownership and when-to-read descriptions.
Make current canonical documents reachable through the relevant index hierarchy.
Identify historical, generated, and draft material so it cannot impersonate current authority.
Introduce nested `AGENTS.md` only for persistent local operating differences.
Avoid copying root instructions into each nested file.

## Size and splitting

Use the thresholds in `SKILL.md` as this skill's default conventions.
Ordinary maintenance reviews changed documents and affected dependencies;
a full audit inventories the requested documentation tree.
Split by independently useful task context, authority, or rate of change.
Retain cohesive exceptions with a reason and useful section headings.
Consider bytes and actual token measurements when available; do not infer context fit
from physical lines. Do not pack dense paragraphs into fewer lines.
Keep routing overhead proportional; many tiny files can make retrieval worse.

## Evolution and enforcement

Update canonical docs in the same change when contracts, boundaries, assumptions,
operations, or accepted decisions change. Leave unrelated documents alone.
Move completed plans' durable outcomes into canonical docs and evidence.
Retrieve changing external facts from official sources; preserve project-specific choices.
Translate stable, measurable invariants into relevant tests, schemas, linters, or CI.
Reuse existing checks and established architecture rules.
Keep semantic review separate from structural linting; neither guarantees future behavior.
