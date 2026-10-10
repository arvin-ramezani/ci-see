# Restructure and Maintain Documentation

## Capture the original

Record the current commit and affected uncommitted contents. A commit alone does
not capture working-tree edits. Preserve recoverable source snapshots outside the
canonical documentation tree when needed.
Identify incoming references with repository search before moving files.
Read the entire affected source and the dependencies needed to interpret it;
headings alone are insufficient for a migration.

## Decide the boundaries

Review responsibility, authority, duplication, change frequency, and likely tasks.
Explain the task that would read each destination.
Allow retrieval of a complete behavior with its conditions and failure rules without
a long chain of tiny files. Keep coupled tables, state transitions, and rules together.
Separate historical rationale from current behavior while retaining cross-links.
Create a topic index when independently useful documents need navigation.

## Build the preservation map

Identify requirements and subclauses, definitions, assumptions, invariants, decisions
and rationale, edge cases, ordering guarantees, exceptions, interpretive examples,
tests, and acceptance criteria before rewriting.
Map each unit to a destination path and heading or stable ID.
Group unchanged contiguous units only with explicit, verifiable scope.
Mark duplicates `DEDUPLICATED` with their canonical owner and evidence that
wording, conditions, and scope are equivalent.
Record intentional semantic changes separately with their authorization.
Do not mark missing material obsolete based on age or brevity alone.

## Execute

Move or rewrite according to the map. Preserve:

- `MUST` versus `SHOULD`, positive and negative obligations;
- actors, permissions, boundaries, prerequisites, supported platforms;
- timeouts, defaults, units, errors, cancellation, concurrency;
- state transitions, durability, ordering;
- requirement IDs, decision status, rationale, rejected alternatives;
- useful diagrams, tables, examples, source links, and unresolved questions.

Repair relative links from each new location. Search the repository for old paths,
heading fragments, and identifiers, including scripts, workflows, prompts, and comments.
Use a compatibility stub or update known external references when a published
canonical URL must remain usable. Give stubs a replacement link without duplicating rules.
Update indexes with ownership and task cues. Remove retired files only after
accounting for their content and references in a recoverable baseline.

## Compare in a fresh pass

Read original and resulting content side by side using the baseline, not memory.
Check every mapped unit and destination for weakened, widened, missing, or invented rules.
Inspect cross-file dependencies and conditional behavior together.
Check unresolved conflicts and accepted history remain identifiable.
Fix omissions before handing off.

Use a separate reviewer for a consequential migration when available and authorized.
Provide raw baseline and resulting artifacts, not a conclusion to confirm.
Otherwise describe the comparison as self-review.
Flag unmatched units as blocking. Report semantic changes only when separately authorized.
For a code-driven contract update, distinguish authorized behavior changes from
the structural migration and inspect relevant implementation evidence.

## Complete

Run navigation checks. Record manual anchor and non-Markdown reference checks
when the available checker does not cover them.
Produce the preservation map and size decisions using the report template.
Retain a cohesive large document when its review supports that choice; record why.
