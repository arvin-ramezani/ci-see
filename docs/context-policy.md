# CI See — Context Engineering Policy

**Scope:** Repository knowledge and agent routing; does not replace the product or
security contracts. **Upstream skill:** `arvin-ramezani/ai-skills` at `6ec528e498b51a4e1a9796bb55469bc188e5fdc1`.

## Retrieval and authority

Use [documentation index](index.md) to choose relevant canonical sources.
The PRD owns product requirements, architecture owns Go/system boundaries,
the three focused specs own their behavioral contracts, and the implementation
plan owns execution order. Preserve each source's actual status and unresolved
owner decisions. Use links rather than duplicate contracts; synchronize affected
documentation when authorized code changes contract behavior.

## Review and restructure

- Prefer **80–100 lines** for focused new prose when complete.
- Allow **100–150 lines** for coherent topics.
- For scoped documents **over 150 lines**, record KEEP/SPLIT/REORGANIZE with
  a concrete reason; size alone never mandates a split.
- For **300–400+ lines**, actively seek independent retrieval units, while
  retaining justified cohesive exceptions. Lines are not token limits.

Apply [Context Engineering](../.agents/skills/context-engineering/SKILL.md) for
baseline snapshots, mapping every affected requirement/decision/condition/example
to the destination, justified in-scope splitting, incoming-link/index repairs,
and a fresh original-vs-result semantic comparison. Never lose negative cases,
ordering guarantees or accepted meaning. Routine authorized structural work has
no arbitrary file-count approval gate. Escalate only out-of-scope/destructive
actions and unresolved product, security, architecture or contract decisions.
Keep conflicting evidence visible and finish unaffected work. Read-only audits
must never write files.

## Verification

Run relevant repository checks. The bundled scanner reports common links,
index reachability, counts and hashes—not semantic equivalence, fragments,
complex nested Markdown or implementation currency. Declare actual checks.
Existing long documents are **not** split by this installation; review them
during an authorized audit or edit. The Python helper is not a Go runtime dependency.
